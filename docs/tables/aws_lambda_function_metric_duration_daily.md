---
title: "Steampipe Table: aws_lambda_function_metric_duration_daily - Query AWS Lambda Function Metrics using SQL"
description: "Allows users to query AWS Lambda Function daily duration metrics."
folder: "Lambda"
---

# Table: aws_lambda_function_metric_duration_daily - Query AWS Lambda Function Metrics using SQL

The AWS Lambda Function Metrics allows you to monitor and troubleshoot your Lambda functions. It provides near real-time metrics on the performance and health of your functions, including duration, which measures the elapsed wall clock time from when the function code starts executing as a result of an invocation to when it stops executing. The daily metric duration gives insights into the function's performance and efficiency over a 24-hour period.

## Table Usage Guide

The `aws_lambda_function_metric_duration_daily` table in Steampipe provides you with information about daily duration metrics of AWS Lambda functions. This table allows you, as a DevOps engineer, to query duration-specific details, including average, maximum, and minimum execution times of functions, along with the total count of requests. You can utilize this table to monitor the performance of Lambda functions, such as identifying functions with long execution times, tracking daily changes in function duration, and more. The schema outlines the various attributes of the Lambda function duration metrics for you, including the timestamp, function name, region, and associated statistics.

The `aws_lambda_function_metric_duration_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Analyze the daily performance metrics of AWS Lambda functions to understand their execution duration trends. This can help in assessing the efficiency of your serverless applications and optimizing resource allocation.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have a daily maximum duration over 100 milliseconds in your Lambda functions. This allows you to identify potential performance issues and optimize your functions for better efficiency.

```sql+postgres
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

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_duration,
  round(maximum,2) as max_duration,
  round(average,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where maximum > 100
order by
  name,
  timestamp;
```

### Lambda function daily average duration less than 5 milliseconds
Identify instances where the average daily duration of AWS Lambda functions is less than 5 milliseconds, in order to assess the efficiency and performance of your serverless applications. This information can be used to optimize resource allocation and improve application responsiveness.

```sql+postgres
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

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_duration,
  round(maximum,2) as max_duration,
  round(average,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where average < 5
order by
  name,
  timestamp;
```