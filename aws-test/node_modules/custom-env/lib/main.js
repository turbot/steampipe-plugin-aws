/*
 * This library wraps the popular node dotenv module
 * For it's configuration
 *
 */

var fs = require('fs')
var path = require('path')
var {
    isObj
} = require('./funcs')

var NODE_ENV = process.env.NODE_ENV

function main(env = true, dir = null, encoding = null, defaultEnvFallback = true) {
    this.env = env
    this.dir = dir
    this.encoding = encoding

    this.customEnvConfig = {
        dir: this.dir,
        envname: this.env,
        defaultEnvFallback
    }

    this.configDotEnv = {
        path: null,
        encoding: this.encoding
    }

    this.setCustomEnvConfig = function (obj) {
        // var dir = typeof (dir) == 'string' ? dir : process.cwd()
        // var envname = envname
        if (obj != null) {
            if (obj && typeof obj == 'object' && obj.constructor.name.toLowerCase() == 'object' && Object.keys(obj).length > 0) {
                for (var v in obj) {
                    this.customEnvConfig[v] = obj[v]
                }
            }
        }
    }

    this.getCustomEnvConfig = function () {
        return this.customEnvConfig
    }


    this.setDotEnvConfig = function (obj = null) {
        if (obj != null) {
            if (typeof obj == 'object' && obj.constructor.name.toLowerCase() == 'object') {
                for (var v in obj) {
                    this.configDotEnv[v] = obj[v]
                }
            }
        }
    }

    this.getDotEnvConfig = function () {
        return this.configDotEnv
    }

    this.getEnvFile = function (dir, envname = null) {
        // Check if dir exists and its a directory
        var dir = typeof (dir) == 'string' ? dir : process.cwd()
        envname = typeof (envname) == 'string' ? envname : envname === true && process.env.NODE_ENV ? process.env.NODE_ENV : undefined
        var defaultEnvRegex = "\.env$"
        var envRegex = defaultEnvRegex


        var dirIsDirectory = fs.lstatSync(dir).isDirectory()

        if (!dirIsDirectory) {
            throw new Error('dir must be a directory')
            return false
        }


        var filesInDir = fs.readdirSync(dir)

        // console.log(filesInDir)

        if (Array.isArray(filesInDir) && filesInDir.length > 0) {
            if (envname) {
                //@UPDATE Update envname to support both `dev or development`
                if (['dev', 'development'].includes(envname)) {
                    envRegex = {
                        "dev": "\.env\.dev?$",
                        "development": "\.env\.development?$"
                    }
                } else {
                    envRegex = "\.env\." + envname + "$"
                }

            } else {
                if (!defaultEnvFallback) {
                    return false;
                }
                envRegex = defaultEnvRegex
                console.warn("No env file present for the current environment: ", envname, '\n Falling back to .env config');
            }

            let testCase
            let objKeysArray
            let keyIndex
            let isEnvObject = isObj(envRegex);



            for (var file of filesInDir) {
                if (isEnvObject) {
                    //@UPDATE check if envname matches `dev or development`
                    objKeysArray = Object.keys(envRegex)
                    keyIndex = objKeysArray.indexOf(envname)

                    testCase = new RegExp(envRegex[envname]).test(file)

                } else {
                    testCase = new RegExp(envRegex).test(file)
                }

                if (testCase) {
                    return path.join(dir, file);
                }
            }

            // Now check for .dev|.development substitutes
            if (isEnvObject) {
                for (var file of filesInDir) {
                    if (new RegExp(envRegex[objKeysArray[1 - keyIndex]]).test(file)) {
                        return path.join(dir, file);
                    }
                }
            }
            // throw new Error("env files not present in directory " + dir)

            if (defaultEnvFallback && envRegex != defaultEnvRegex) {
                console.warn("No env file present for the current environment: ", envname, '\n Falling back to .env config');
                for (var file of filesInDir) {
                    if (new RegExp(defaultEnvRegex).test(file)) {
                        return path.join(dir, file);
                    }
                }
            }

        } else {
            if (defaultEnvFallback && envRegex != defaultEnvRegex) {
                console.warn("No env file present for the current environment: ", envname, '\n Falling back to .env config');
                for (var file of filesInDir) {
                    if (new RegExp(defaultEnvRegex).test(file)) {
                        return path.join(dir, file);
                    }
                }
            }
        }

        console.warn("No env file present for the current environment: ", envname)
        return false;
    }


    this.loadDotEnvConfig = function () {
        var dotenvobj = this.getDotEnvConfig()

        if (dotenvobj.path) {
            require('dotenv-expand')(require('dotenv').config(dotenvobj));
        }
    }


    this.loadCustomDotEnv = function () {
        // setTheCustomEnv
        this.setCustomEnvConfig(this.dir, this.env)

        var currentCustomEnvConfig = this.getCustomEnvConfig()
        // Set and load dotenv
        this.setDotEnvConfig({
            'path': this.getEnvFile(currentCustomEnvConfig.dir, currentCustomEnvConfig.envname)
        })


        // Start dotenv with the available settings
        this.loadDotEnvConfig()

    }

    this.loadCustomDotEnv()
}

// Create a new constructor to make some properties available

var pubLicMain = new Object()

module.exports = {
    env: function () {
        this.main = new main(arguments[0], arguments[1])
        if (arguments[2] === true) {
            this.main.setCustomEnvConfig({
                'envname': process.env.NODE_ENV
            })
        }

        if (arguments[3] === false) {
            this.main.setCustomEnvConfig({
                defaultEnvFallback: false
            })
        }
        return this
    },
    dotenvConfig: function () {
        this.main.setDotEnvConfig(...arguments)
        return this
    },
    config: function () {
        this.main.setCustomEnvConfig(...arguments)
        return this
    }
}