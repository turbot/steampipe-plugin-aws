#!/usr/bin/env node

const _ = require("lodash");
const chalk = require("chalk");
const diff = require("diff");
const fs = require("fs-extra");
const mm = require("micromatch");
const os = require("os");
const path = require("path");
const { spawn } = require("child_process");
const envVariable = require("custom-env");
const json_diff = require("json-diff");

const TMP_DIR = os.tmpdir();

const _targetGlobs = function () {
  var args = process.argv.slice(2);
  if (args.length == 0) {
    args.push("*");
  }
  var targets = args.map(i => {
    var tmp;
    if (i.includes("/")) {
      tmp = i;
    } else {
      tmp = "tests/" + i;
    }
    // Strip a trailing slash if there is one
    if (tmp[tmp.length - 1] == "/") {
      return tmp.slice(0, -1);
    }
    return tmp;
  });
  return targets;
};

const _targetDirs = function (targetGlobs) {
  var targetDirs = targetGlobs.map(i => {
    return i.split("/")[0];
  });
  return [...new Set(targetDirs)];
};

const _availableTests = function (targetDirs) {
  var result = [];
  for (const d of targetDirs) {
    result = result.concat(
      fs.readdirSync(d).map(i => {
        return d + "/" + i;
      })
    );
  }
  return result;
};

const _testPrerequisites = function (testDir) {
  var dependenciesSrc;
  try {
    dependenciesSrc = fs.readFileSync(testDir + "/dependencies.txt", "utf8");
  } catch (e) {
    return [];
  }
  if (!dependenciesSrc) {
    return [];
  }
  var dependencies = dependenciesSrc.split(/\s+/g);
  var testsDir = path.resolve(testDir, "..");
  var integrationTestsDir = path.resolve(testsDir, "..");
  var prereqs = _.compact(dependencies).map(i => {
    return path.relative(integrationTestsDir, path.resolve(testsDir, i));
  });
  return prereqs;
};

const _testsWithPrerequisites = function (tests) {
  var result = [];
  for (const t of tests) {
    var prereqs = _testPrerequisites(t);
    var testObj = {
      dir: t,
      failed: false,
      tmpDir: path.resolve(TMP_DIR, t),
      prereqs,
      output: {},
      setup: {},
      pretest: {},
      test: {},
      posttest: {},
      teardown: {}
    };
    result.push(_testsWithPrerequisites(prereqs));
    result.push(testObj);
  }
  var flatResult = _.flattenDeep(result);
  return [...new Set(flatResult)];
};

const _renderToTmp = function (
  test,
  sourcePath,
  defaultRendered = null,
  tmpDir = null
) {
  _.templateSettings.interpolate = /{{([\s\S]+?)}}/g;
  var template, tmpPath, rendered;
  tmpPath = path.resolve(tmpDir || test.tmpDir, path.basename(sourcePath));
  try {
    const templateSrc = fs.readFileSync(sourcePath, "utf8");
    template = _.template(templateSrc);
    rendered = template(test);
  } catch (e) {
    if (defaultRendered) {
      rendered = defaultRendered;
    } else {
      e.sourcePath = sourcePath;
      throw e;
    }
  }
  fs.writeFileSync(tmpPath, rendered);
  return tmpPath;
};

const _runTerraformInit = async function (test, phase) {
  return new Promise((resolve, reject) => {
    const tfDir = path.resolve(test.tmpDir, "terraform", phase);
    const args = ["init"];
    const cmd = spawn("terraform", args, { encoding: "utf8", cwd: tfDir });

    var result = {
      stdout: "",
      stderr: ""
    };
    var resultSent = false;

    console.log(chalk.yellow("Running terraform"));
    cmd.stdout.on("data", data => {
      result.stdout += data;
    });

    cmd.stderr.on("data", data => {
      process.stdout.write(chalk.gray(data));
      result.stderr += data;
    });

    cmd.on("exit", code => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      result.status = code;
      try {
        // TODO: What is result.output used for?
        result.output = result.stdout;
        _.merge(test.test, { terraform: { init: result } });
      } catch (e) {
        result.output = {};
      }
      resolve(test);
    });

    cmd.on("error", err => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      reject(err);
    });
  });
};

const _runTerraformApply = function (test, phase) {
  return new Promise((resolve, reject) => {
    const tfDir = path.resolve(test.tmpDir, "terraform", phase);

    let environment = Object.assign(process.env, {
      TF_VAR_resource_name: process.env.TURBOT_TEST_RESOURCE_NAME,
      TF_VAR_resource_name_1: process.env.TURBOT_TEST_RESOURCE_NAME_1,
      TF_VAR_resource_name_2: process.env.TURBOT_TEST_RESOURCE_NAME_2
    });
    const args = ["apply", "-auto-approve", "-no-color"];
    const cmd = spawn("terraform", args, {
      encoding: "utf8",
      cwd: tfDir,
      env: environment
    });

    var result = {
      stdout: "",
      stderr: ""
    };
    var resultSent = false;

    cmd.stdout.on("data", data => {
      process.stdout.write(chalk.dim(data));
      result.stdout += data;
    });

    cmd.stderr.on("data", data => {
      process.stdout.write(chalk.gray(data));
      result.stderr += data;
    });

    cmd.on("exit", code => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      result.status = code;
      try {
        // TODO: What is result.output used for?
        result.output = result.stdout;
        _.merge(test[phase], { terraform: { apply: result } });
      } catch (e) {
        result.output = {};
      }
      resolve(test);
    });

    cmd.on("error", err => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      reject(err);
    });
  });
};

const _runTerraformDestroy = function (test, phase = "test") {
  return new Promise((resolve, reject) => {
    const tfDir = path.resolve(test.tmpDir, "terraform", phase);
    console.log("");
    const args = ["destroy", "-auto-approve", "-no-color"];

    const cmd = spawn("terraform", args, {
      encoding: "utf8",
      cwd: tfDir
    });

    var result = {
      stdout: "",
      stderr: ""
    };
    var resultSent = false;

    cmd.stdout.on("data", data => {
      result.stdout += data;
    });

    cmd.stderr.on("data", data => {
      process.stdout.write(chalk.gray(data));
      result.stderr += data;
    });

    cmd.on("exit", code => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      result.status = code;

      try {
        // TODO: What is result.output used for?
        result.output = result.stdout;
        _.set(test.test, ["terraform", phase, "destroy"], result);
      } catch (e) {
        result.output = {};
      }
      resolve(test);
    });

    cmd.on("error", err => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      reject(err);
    });
  });
};

const _runTerraformOutput = async function (test, phase) {
  return new Promise((resolve, reject) => {
    const tfDir = path.resolve(test.tmpDir, "terraform", phase);
    const args = ["output", "--json"];
    const cmd = spawn("terraform", args, { encoding: "utf8", cwd: tfDir });

    var result = {
      stdout: "",
      stderr: ""
    };
    var resultSent = false;

    cmd.stdout.on("data", data => {
      result.stdout += data;
    });

    cmd.stderr.on("data", data => {
      process.stdout.write(chalk.gray(data));
      result.stderr += data;
    });

    cmd.on("exit", code => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      result.status = code;
      try {
        let terraformOutput = JSON.parse(result.stdout);
        Object.assign(test.output, {
          resourceId: _.get(terraformOutput, "resource_id.value", "")
        });
        Object.assign(test.output, {
          resourceName: _.get(terraformOutput, "resource_name.value", "")
        });
        Object.assign(test, {
          resourceId: _.get(terraformOutput, "resource_id.value", "")
        });
        Object.assign(test, {
          resourceName: _.get(terraformOutput, "resource_name.value", "")
        });
        Object.assign(test.output, terraformOutput);
      } catch (e) {
        result.output = {};
      }
      resolve(test);
    });

    cmd.on("error", err => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      reject(err);
    });
  });
};

const _runTerraformApplyForTestPhase = async function (test, phase) {
  const files = fs
    .readdirSync(test.dir)
    .sort()
    .filter(i => {
      if (![".tf", ".tfvars"].includes(path.extname(i))) {
        return false;
      }
      var basename = path.basename(i);
      if (phase != "test") {
        return basename.startsWith(phase + "-");
      } else {
        return !/^(setup|pretest|posttest|teardown)-/.test(basename);
      }
    });
  if (_.isEmpty(files)) {
    return test;
  }
  const tfDir = path.resolve(test.tmpDir, "terraform", phase);

  try {
    fs.ensureDirSync(tfDir);
    // TODO: We don't use the tmp var, what is this for loop for?
    for (const f of files) {
      // Due to a bug with ncc (https://github.com/zeit/ncc/issues/444), we
      // need to pass in the test dir and filename together instead of passing
      // them in as separate arguments
      var testDirWithFile = test.dir + "/" + f;
      _renderToTmp(test, path.resolve(".", testDirWithFile), null, tfDir);
    }
  } catch (e) {
    console.log(chalk.red(`Template rendering error in: ${test.dir}`));
    console.log(chalk.red(e.message));
    throw e;
  }

  test = await _runTerraformInit(test, phase);
  test = await _runTerraformApply(test, phase);
  test = await _runTerraformOutput(test, phase);

  return test;
};

const _runGraphqlQuery = function (test, query) {
  return new Promise((resolve, reject) => {
    try {
      var queryTmp = _renderToTmp(test, query.query);
      // var variablesTmp = _renderToTmp(test, query.variables, {});
      var expectedTmp = _renderToTmp(test, query.expected);
    } catch (e) {
      console.log(chalk.red(`Template Error: ${e.sourcePath}`));
      console.log(chalk.red(e.message));
      throw e;
    }

    console.log(
      chalk.yellow(`\nRunning SQL query: ${path.basename(query.query)}`)
    );

    /*
    const args = [
      "graphql",
      "--query",
      queryTmp,
      "--variables",
      variablesTmp,
      "--expected",
      expectedTmp,
      "--format",
      "json",
      "--timeout",
      process.env.TURBOT_TEST_EXPECTED_TIMEOUT || 180
    ];
    const cmd = spawn("turbot", args, { encoding: "utf8" });
    */

    q = fs.readFileSync(queryTmp, { encoding: "utf8" });
    q = q.replace(/\n/, " ");
    q = q.replace(/\r/, "");
    // console.log({ q, queryTmp });
    const args = [
      "query",
      "--output",
      "json",
      q
      //"select name from aws_s3_bucket whe order by name"
    ];
    const cmd = spawn("steampipe", args, { encoding: "utf8" });

    var result = {
      stdout: "",
      stderr: ""
    };
    var resultSent = false;

    cmd.stdout.on("data", data => {
      result.stdout += data;
    });

    cmd.stderr.on("data", data => {
      process.stdout.write(chalk.dim(data));
      result.stderr += data;
    });

    cmd.on("exit", code => {
      // console.log(result);
      if (resultSent) {
        return;
      }
      resultSent = true;
      result.status = code;
      try {
        result.output = JSON.parse(result.stdout);
      } catch (e) {
        result.output = {};
      }
      if (true) {
        // if (code) { <<<< WHY?
        var outputStr = JSON.stringify(result.output, null, 2) + "\n";
        var expectedStr =
          JSON.stringify(JSON.parse(fs.readFileSync(expectedTmp)), null, 2) +
          "\n";
        var differences = diff.diffJson(outputStr, expectedStr);

        differences.forEach(part => {
          if (part.added) {
            result.status = 1
            process.stdout.write(chalk.green(part.value));
          } else if (part.removed) {
            result.status = 1
            process.stdout.write(chalk.red(part.value));
          } else {
            process.stdout.write(chalk.dim(part.value));
          }
        });
      }
      resolve(result);
    });

    cmd.on("error", err => {
      if (resultSent) {
        return;
      }
      resultSent = true;
      reject(err);
    });
  });
};

const _runGraphqlQueriesForTestPhase = async function (
  test,
  phase,
  terraformSuccessful
) {
  const queries = fs
    .readdirSync(test.dir)
    .sort()
    .map(i => {
      if (phase == "test" && i == "query.sql") {
        return {
          name: "",
          phase,
          query: `${test.dir}/query.sql`,
          variables: `${test.dir}/variables.json`,
          expected: `${test.dir}/expected.json`
        };
      }
      if (i.startsWith(phase + "-") && i.endsWith("-query.sql")) {
        const name = i.slice((phase + "-").length, -"-query.sql".length);
        return {
          name,
          phase,
          query: `${test.dir}/${phase}-${name}-query.sql`,
          variables: `${test.dir}/${phase}-${name}-variables.json`,
          expected: `${test.dir}/${phase}-${name}-expected.json`
        };
      }
      return null;
    })
    .filter(i => !!i);

  // console.log({queries: queries})
  for (const q of queries) {
    var queryResult;
    // Skip running GraphQL to avoid wasting time if the Terraform run
    // wasn't successful
    if (terraformSuccessful) {
      queryResult = await _runGraphqlQuery(test, q);
    } else {
      console.log(
        chalk.redBright.bold("Terraform run failed, skipping SQL queries")
      );
      // TODO: Do we need additional info in the queryResult object?
      queryResult = {
        status: 1
      };
    }
    // console.log({status: queryResult.status})
    if (queryResult.status) {
      // If any of the queries in the test fail, then the test has failed
      // and should not run further steps.
      test.failed = true;
      console.log("");
      console.log(chalk.redBright.bold("✘ FAILED"));
    } else {
      if (q.phase == "test") {
        console.log(chalk.greenBright.bold("✔ PASSED"));
      }
    }
    test[q.phase][q.name] = queryResult;
    _.merge(test.output, queryResult.output);
  }
  return test;
};

const _run = async function (tests) {
  var testNames = tests.map(i => i.dir);
  var results = _.keyBy(tests, "dir");

  // Set the test resource names before running any steps to allow same name
  // use during consecutive tests
  process.env.TURBOT_TEST_RESOURCE_NAME_PREFIX =
    process.env.TURBOT_TEST_RESOURCE_NAME_PREFIX || "turbottest";
  process.env.TURBOT_TEST_RESOURCE_NAME =
    process.env.TURBOT_TEST_RESOURCE_NAME ||
    process.env.TURBOT_TEST_RESOURCE_NAME_PREFIX + _.random(100000);
  process.env.TURBOT_TEST_RESOURCE_NAME_1 =
    process.env.TURBOT_TEST_RESOURCE_NAME_1 ||
    process.env.TURBOT_TEST_RESOURCE_NAME_PREFIX + _.random(100000);
  process.env.TURBOT_TEST_RESOURCE_NAME_2 =
    process.env.TURBOT_TEST_RESOURCE_NAME_2 ||
    process.env.TURBOT_TEST_RESOURCE_NAME_PREFIX + _.random(100000);

  try {
    for (const t of testNames) {
      results[t] = await _runSetup(results[t]);
      const phases = ["pretest", "test", "posttest"];
      for (const phase of phases) {
        results[t] = await _runTestPhase(results[t], phase);
      }
    }
  } catch (e) {
    console.log(
      chalk.bold.red(
        "\nERROR DETECTED: Stopping test run and entering teardown phase."
      )
    );
    console.log(chalk.bold.red(e.message));
    if (process.env.TURBOT_TEST_LOG_LEVEL == "debug") {
      console.log(e);
    }
  }
  // Teardown tests in reverse order
  for (var i = testNames.length - 1; i >= 0; i--) {
    const t = testNames[i];
    results[t] = await _runTeardown(results[t]);
  }
  return results;
};

const _runSetup = async function (test) {
  /**
   * Set the environment variable from the env file
   *
   * To set the environment variable create a file with
   * extension ".env.staging"
   */
  envVariable.env("staging", `${test.dir}`);
  console.log(
    "customEnv TURBOT_TEST_EXPECTED_TIMEOUT",
    process.env.TURBOT_TEST_EXPECTED_TIMEOUT
  );

  if (test.failed) {
    return test;
  }
  console.log(chalk.bold(`\nSETUP: ${test.dir} [${test.prereqs}]`));
  fs.ensureDirSync(test.tmpDir);
  return test;
};

const _runTestPhase = async function (test, phase) {
  if (test.failed) {
    return test;
  }
  console.log(chalk.bold(`\n${phase.toUpperCase()}: ${test.dir}`));
  test = await _runTerraformApplyForTestPhase(test, phase);
  let terraformSuccessful = true;
  // Default to 0 if there is no status because this phase didn't have
  // Terraform to run
  terraformSuccessful =
    _.get(test, "test.terraform.init.status", 0) === 0 &&
    _.get(test, "test.terraform.apply.status", 0) === 0;
  test = await _runGraphqlQueriesForTestPhase(test, phase, terraformSuccessful);
  return test;
};

const _runTeardown = async function (test) {
  console.log(chalk.bold(`\nTEARDOWN: ${test.dir}`));
  // Check for pretest / test / posttest keys
  let getKeys = Object.keys(test);
  if (getKeys.includes("pretest") && !_.isEmpty(test.pretest)) {
    test = await _runTerraformDestroy(test, "pretest");
  }
  if (getKeys.includes("test") && !_.isEmpty(test.test)) {
    test = await _runTerraformDestroy(test, "test");
  }
  if (getKeys.includes("posttest") && !_.isEmpty(test.posttest)) {
    test = await _runTerraformDestroy(test, "posttest");
  }
  fs.removeSync(test.tmpDir);
  return test;
};

async function main() {
  const targetGlobs = _targetGlobs();
  const targetDirs = _targetDirs(targetGlobs);
  const availableTests = _availableTests(targetDirs);
  const targets = mm(availableTests, targetGlobs);
  const resolvedTargets = _testsWithPrerequisites(targets);
  if (!resolvedTargets.length) {
    return console.log("No matching targets. Stopping.");
  }
  var result = await _run(resolvedTargets);
  if (process.env.TURBOT_TEST_LOG_LEVEL == "debug") {
    console.log(JSON.stringify(result, null, 2));
  }
  var numTests = resolvedTargets.length;
  var numTestsPassing = _.filter(result, testResult => !testResult.failed)
    .length;
  var failingTests = _.chain(result)
    .map(i => {
      return i.failed ? i : null;
    })
    .compact()
    .value();
  var numTestsFailing = failingTests.length;
  console.log(chalk.bold("SUMMARY:"));
  if (numTestsFailing) {
    console.log("");
    for (let i of failingTests) {
      console.log(chalk.red(`✘ ${i.dir} failed.`));
    }
  }
  var summaryColor = "redBright";
  if (numTestsPassing == numTests) {
    summaryColor = "greenBright";
  } else if (numTestsPassing) {
    summaryColor = "yellowBright";
  }
  console.log("");
  console.log(chalk[summaryColor](`${numTestsPassing}/${numTests} passed.`));
  console.log();
  process.exit(numTests - numTestsPassing);
}

if (process.env.TURBOT_TEST_COLOR_LEVEL == "0") {
  chalk.level = 0;
}

main();
