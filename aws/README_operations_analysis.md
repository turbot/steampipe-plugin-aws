# AWS Operations Analyzer

This is one of the two test tools that validate that our IAM policy allows all the queries in our `extractor-steampipe` project. This test is designed to statically analyze Steampipe table definitions and from the plugin's source code and SQL queries from Catio's JSON definition. It determines the complete set of AWS API operations that could be called by a set of queries and tests them against the AWS IAM Policy Simulator.

## Quickstart

Run from the root directory of this project. Assumes you have the `extractor-steampipe` project cloned in a sibling directory. If not, change the `queryFile` path in `aws/operations_analyzer_test.go`

Run only the tests that validate the queries from json:

```sh
setDev
go test -v -run TestAnalyzeQueriesFromJSON ./aws
```

Run all the unit tests for the feature:

```sh
setDev
go test -v -run TestAnalyze ./aws
```

## How It Works

### Key Functions

The analyzer exposes three main functions:

- `GetAWSOperationsFromSQL(ctx, sqlQuery)`: This is the most common entry point. It takes a SQL query string and returns a flat, de-duplicated, and sorted slice of all unique AWS operations that could be triggered.
- `AnalyzeSQLQuery(ctx, sqlQuery)`: This function is similar but returns a `map[string][]string`, grouping the discovered AWS operations by the table that requires them.
- `AnalyzeTableOperations(ctx, table)`: This is a lower-level function that analyzes a single `*plugin.Table` object and its dependencies.

The analyzer uses a combination of dynamic plugin introspection and sophisticated SQL parsing to build a complete picture of the required API calls.

### 1. Dynamic Table Loading & Caching

On its first run, the analyzer inspects the `Plugin()` definition to get a complete list of all tables available in the `steampipe-plugin-aws` plugin. It builds and caches two key maps:

- A map of table names to their `*plugin.Table` definitions.
- A map of `List` function names to the table they belong to.

This dynamic approach means that as new tables are added to the plugin, the analyzer will automatically discover them without needing any manual code changes or registration.

### 2. Hydration Graph Traversal

When analyzing a table, the analyzer inspects the following parts of its definition to find operations tagged with a `service` and `action`:

- `Get`
- `List`
- `HydrateConfig`

Crucially, if it finds a `ParentHydrate` function in the `List` config, it uses the cached map to find the corresponding parent table and adds it to a queue for analysis. This recursive process ensures that the operations from all dependent tables in the hydration chain are discovered.

### 3. SQL Parsing

To identify which tables are involved in a query, the analyzer uses the `pganalyze/pg_query_go` library, a robust PostgreSQL parser. This allows it to reliably extract table names even from complex SQL statements involving `JOINs`, `CTEs` (Common Table Expressions), and nested subqueries.
