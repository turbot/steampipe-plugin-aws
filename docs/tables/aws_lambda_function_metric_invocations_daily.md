---
title: "Table: aws_lambda_function_metric_invocations_daily - Query AWS Lambda Function Metrics using SQL"
description: "Allows users to query AWS Lambda Function Metrics on a daily basis."
---

# Table: aws_lambda_function_metric_invocations_daily - Query AWS Lambda Function Metrics using SQL

The `aws_lambda_function_metric_invocations_daily` table in Steampipe provides information about the daily invocation metrics of AWS Lambda functions. This table enables DevOps engineers to query function-specific details, including the number of invocations, the function name, and the timestamp of the data point. Users can utilize this table to gather insights on function usage, such as the number of invocations over time, peak usage times, and more. The schema outlines the various attributes of the Lambda function metrics, including the function name, the namespace, the metric name, and the timestamp.

The `aws_lambda_function_metric_invocations_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_function_metric_invocations_daily` table, you can use the `.inspect aws_lambda_function_metric_invocations_daily` command in Steampipe.

**Key columns**:

- `function_name`: The name of the Lambda function. This is useful for joining with other tables that contain Lambda function details.
- `namespace`: The namespace of the function. This can be used to join with other tables that contain AWS namespace information.
- `timestamp`: The timestamp of the data point. This allows for temporal analysis of function invocations over time.

## Examples

### Basic info

```sql
select
  name,
  timestamp,
  sum
from
  aws_lambda_function_metric_invocations_daily
order by
  name,
  timestamp;
```


### Lambda function daily invocations over 10 in last 3 days

```sql
select
  name,
  timestamp,
  round(sum::numeric,2) as sum_invocations,
  sample_count
from
  aws_lambda_function_metric_invocations_daily
where 
  date_part('day', now() - timestamp) <=3
and sum > 10
order by
  name,
  timestamp;
```