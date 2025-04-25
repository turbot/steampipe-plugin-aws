---
title: "Steampipe Table: aws_lambda_function_metric_invocations_daily - Query AWS Lambda Function Metrics using SQL"
description: "Allows users to query AWS Lambda Function Metrics on a daily basis."
folder: "Lambda"
---

# Table: aws_lambda_function_metric_invocations_daily - Query AWS Lambda Function Metrics using SQL

The AWS Lambda Function Metrics service allows you to monitor and troubleshoot your Lambda functions. It provides real-time metrics with granularity down to one minute and trace data sampling. By querying these metrics using SQL, you can gain insights into function execution such as memory usage, execution time, and failure rates.

## Table Usage Guide

The `aws_lambda_function_metric_invocations_daily` table in Steampipe provides you with information about the daily invocation metrics of AWS Lambda functions. This table enables you, as a DevOps engineer, to query function-specific details, including the number of invocations, the function name, and the timestamp of the data point. You can utilize this table to gather insights on function usage, such as the number of invocations over time, peak usage times, and more. The schema outlines the various attributes of the Lambda function metrics for you, including the function name, the namespace, the metric name, and the timestamp.

The `aws_lambda_function_metric_invocations_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Gain insights into the daily usage patterns of your AWS Lambda functions. This query helps to understand the frequency and timing of function invocations, which can aid in optimizing resource allocation and cost management.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which AWS Lambda functions are being invoked more than 10 times daily over the past three days. This is useful for tracking function usage and identifying potential areas of optimization or troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  name,
  timestamp,
  round(sum,2) as sum_invocations,
  sample_count
from
  aws_lambda_function_metric_invocations_daily
where 
  julianday('now') - julianday(timestamp) <=3
and sum > 10
order by
  name,
  timestamp;
```