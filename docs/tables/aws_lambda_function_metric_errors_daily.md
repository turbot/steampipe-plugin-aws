---
title: "Steampipe Table: aws_lambda_function_metric_errors_daily - Query AWS Lambda Function using SQL"
description: "Allows users to query AWS Lambda Function error metrics on a daily basis."
folder: "Lambda"
---

# Table: aws_lambda_function_metric_errors_daily - Query AWS Lambda Function using SQL

The AWS Lambda Function is a key component of the AWS Lambda service, allowing you to run your code without provisioning or managing servers. It executes your code only when needed and scales automatically, from a few requests per day to thousands per second. The daily metric errors provide insights into the execution errors of your Lambda Function, enabling you to monitor and troubleshoot performance issues.

## Table Usage Guide

The `aws_lambda_function_metric_errors_daily` table in Steampipe gives you information about the error metrics of AWS Lambda Functions on a daily basis. This table lets you, as a DevOps engineer, data analyst, or other technical professional, query error-specific details, including the number of errors, the timestamp of errors, and associated metadata. You can utilize this table to gather insights on Lambda function errors, such as error frequency, error patterns, and more. The schema outlines the various attributes of the Lambda function error metrics, including the function name, namespace, and statistical information.

The `aws_lambda_function_metric_errors_daily` table provides you with metric statistics at 24-hour intervals for the last year.

## Examples

### Basic info
Determine the areas in which AWS Lambda functions have encountered errors on a daily basis. This allows you to analyze the performance of your functions over time and identify any potential issues that may be causing frequent errors.

```sql+postgres
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

```sql+sqlite
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
This query helps to identify and monitor AWS Lambda functions that are performing well, with a daily average error rate of less than 1. It's beneficial for maintaining high-quality performance and quickly addressing any functions that may not be meeting this standard.

```sql+postgres
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

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_error,
  round(maximum,2) as max_error,
  round(average,2) as avg_error,
  sample_count
from
  aws_lambda_function_metric_errors_daily
where average < 1
order by
  name,
  timestamp;
```