---
title: "Table: aws_lambda_function_metric_errors_daily - Query AWS Lambda Function using SQL"
description: "Allows users to query AWS Lambda Function error metrics on a daily basis."
---

# Table: aws_lambda_function_metric_errors_daily - Query AWS Lambda Function using SQL

The `aws_lambda_function_metric_errors_daily` table in Steampipe provides information about the error metrics of AWS Lambda Functions on a daily basis. This table allows DevOps engineers, data analysts, and other technical professionals to query error-specific details, including the number of errors, the timestamp of errors, and associated metadata. Users can utilize this table to gather insights on Lambda function errors, such as error frequency, error patterns, and more. The schema outlines the various attributes of the Lambda function error metrics, including the function name, namespace, and statistical information.

The `aws_lambda_function_metric_errors_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_function_metric_errors_daily` table, you can use the `.inspect aws_lambda_function_metric_errors_daily` command in Steampipe.

**Key columns**:

- `function_name`: The name of the Lambda function. This can be used to join with other tables that contain Lambda function details.
- `timestamp`: The timestamp of the error metric. This can be used to analyze error patterns over time.
- `region`: The AWS region in which the Lambda function is located. This can be used for regional analysis and to join with other region-specific AWS tables.

## Examples

### Basic info

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_lambda_function_metric_errors_daily
order by
  name,
  timestamp;
```

### Lambda function daily average error less than 1

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_error,
  round(maximum::numeric,2) as max_error,
  round(average::numeric,2) as avg_error,
  sample_count
from
  aws_lambda_function_metric_errors_daily
where average < 1
order by
  name,
  timestamp;
```