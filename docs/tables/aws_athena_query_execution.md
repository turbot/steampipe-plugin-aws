---
title: "Table: aws_athena_query_execution - Query AWS Athena Query Executions using SQL"
description: "Allows users to query AWS Athena Query Executions to retrieve detailed information about each individual query execution."
---

# Table: aws_athena_query_execution - Query AWS Athena Query Executions using SQL

The `aws_athena_query_execution` table in Steampipe provides information about query executions within AWS Athena. This table allows data analysts and developers to query execution-specific details, including execution status, result configuration, and associated metadata. Users can utilize this table to track the progress of queries, analyze the performance of queries, and understand the cost of running specific queries. The schema outlines the various attributes of the Athena query execution, including the query execution id, query, output location, data scanned, and execution time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_athena_query_execution` table, you can use the `.inspect aws_athena_query_execution` command in Steampipe.

Key columns:

- `query_execution_id`: This is the unique identifier for each query execution. It can be used to join this table with others that contain query execution details.
- `query`: This is the SQL query statement that Athena executed. This can be useful for understanding what data was requested in a particular execution.
- `status_state`: This column indicates the current state of the query execution. It can be useful for tracking the progress and success of queries.

## Examples

### List all queries in error

```sql
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

```sql
select 
  workgroup, 
  sum(data_scanned_in_bytes) 
from 
  aws_athena_query_execution
group by 
  workgroup;
```

### Find queries with biggest execution time

```sql
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

```sql
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
