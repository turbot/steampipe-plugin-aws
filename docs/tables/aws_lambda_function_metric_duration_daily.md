---
title: "Table: aws_lambda_function_metric_duration_daily - Query AWS Lambda Function Metrics using SQL"
description: "Allows users to query AWS Lambda Function daily duration metrics."
---

# Table: aws_lambda_function_metric_duration_daily - Query AWS Lambda Function Metrics using SQL

The `aws_lambda_function_metric_duration_daily` table in Steampipe provides information about daily duration metrics of AWS Lambda functions. This table allows DevOps engineers to query duration-specific details, including average, maximum, and minimum execution times of functions, along with the total count of requests. Users can utilize this table to monitor the performance of Lambda functions, such as identifying functions with long execution times, tracking daily changes in function duration, and more. The schema outlines the various attributes of the Lambda function duration metrics, including the timestamp, function name, region, and associated statistics.

The `aws_lambda_function_metric_duration_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_function_metric_duration_daily` table, you can use the `.inspect aws_lambda_function_metric_duration_daily` command in Steampipe.

**Key columns**:

- `timestamp`: This is the timestamp for the data point. It can be used to track the performance of a function over time.
- `function_name`: The name of the Lambda function. This can be used to join with other tables that provide more detailed information about the function.
- `region`: The AWS region in which the function is hosted. This can be used to join with other tables that contain region-specific information.

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
  aws_lambda_function_metric_duration_daily
order by
  name,
  timestamp;
```

### Lambda function daily maximum duration over 100 milliseconds

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_duration,
  round(maximum::numeric,2) as max_duration,
  round(average::numeric,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where maximum > 100
order by
  name,
  timestamp;
```

### Lambda function daily average duration less than 5 milliseconds

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_duration,
  round(maximum::numeric,2) as max_duration,
  round(average::numeric,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where average < 5
order by
  name,
  timestamp;
```