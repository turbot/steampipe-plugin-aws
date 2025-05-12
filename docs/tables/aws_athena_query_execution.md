---
title: "Steampipe Table: aws_athena_query_execution - Query AWS Athena Query Executions using SQL"
description: "Allows users to query AWS Athena Query Executions to retrieve detailed information about each individual query execution."
folder: "Athena"
---

# Table: aws_athena_query_execution - Query AWS Athena Query Executions using SQL

AWS Athena Query Execution is a feature of Amazon Athena that allows you to run SQL queries on data stored in Amazon S3. It executes queries using an interactive query service that leverages standard SQL. This enables you to analyze data directly in S3 without the need for complex ETL jobs.

## Table Usage Guide

The `aws_athena_query_execution` table in Steampipe provides you with information about query executions within AWS Athena. This table allows you, as a data analyst or developer, to query execution-specific details, including execution status, result configuration, and associated metadata. You can utilize this table to track the progress of queries, analyze the performance of queries, and understand the cost of running specific queries. The schema outlines the various attributes of the Athena query execution for you, including the query execution id, query, output location, data scanned, and execution time.

## Examples

### List all queries in error
Explore which queries have resulted in errors to understand the issues and rectify them accordingly. This is useful in identifying and resolving potential problems within your AWS Athena query execution.

```sql+postgres
select
  id,
  query,
  error_message,
  error_type
from
  aws_athena_query_execution
where
  error_message is not null;
```

```sql+sqlite
select
  id,
  query,
  error_message,
  error_type
from
  aws_athena_query_execution
where
  error_message is not null;
```

### Estimate data read by each workgroup
Analyze the volume of data processed by each workgroup to understand workload distribution and optimize resources accordingly. This can be useful in identifying workgroups that are processing large amounts of data and may require additional resources or optimization.

```sql+postgres
select 
  workgroup, 
  sum(data_scanned_in_bytes) 
from 
  aws_athena_query_execution
group by 
  workgroup;
```

```sql+sqlite
select 
  workgroup, 
  sum(data_scanned_in_bytes) 
from 
  aws_athena_query_execution
group by 
  workgroup;
```

### Find queries with biggest execution time
Discover the queries that have the longest execution times to identify potential areas for performance optimization and enhance the efficiency of your AWS Athena operations.

```sql+postgres
select
  id,
  query,
  workgroup,
  engine_execution_time_in_millis 
from
  aws_athena_query_execution 
order by
  engine_execution_time_in_millis limit 5;
```

```sql+sqlite
select
  id,
  query,
  workgroup,
  engine_execution_time_in_millis 
from
  aws_athena_query_execution 
order by
  engine_execution_time_in_millis limit 5;
```

### Find most used databases
Discover the databases that are frequently used in your AWS Athena environment. This can help optimize resource allocation and identify potential areas for performance improvement.

```sql+postgres
select
  database,
  count(id) as nb_query 
from
  aws_athena_query_execution 
group by
  database 
order by
  nb_query limit 5;
```

```sql+sqlite
select
  database,
  count(id) as nb_query 
from
  aws_athena_query_execution 
group by
  database 
order by
  nb_query limit 5;
```