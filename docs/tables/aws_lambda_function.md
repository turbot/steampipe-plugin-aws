---
title: "Steampipe Table: aws_lambda_function - Query AWS Lambda Function using SQL"
description: "Allows users to query AWS Lambda Functions, providing information about each function's configuration, including runtime, code size, timeout, and associated tags."
folder: "Lambda"
---

# Table: aws_lambda_function - Query AWS Lambda Function using SQL

The AWS Lambda Function is a compute service that lets you run code without provisioning or managing servers. AWS Lambda executes your code only when needed and scales automatically, from a few requests per day to thousands per second. You pay only for the compute time you consume - there is no charge when your code is not running.

## Table Usage Guide

The `aws_lambda_function` table in Steampipe provides you with information about AWS Lambda Functions. This table allows you, as a DevOps engineer, to query function-specific details, including the function's runtime, code size, timeout, and associated tags. You can utilize this table to gather insights on functions, such as the function's configuration, handler, last modified date, and more. The schema outlines the various attributes of the AWS Lambda Function for you, including the function name, ARN, description, and associated environment variables.

## Examples

### Basic Info
Explore which AWS Lambda functions are in use and determine if they have a Key Management Service (KMS) key associated with them, which is essential for managing cryptographic keys. This is beneficial for understanding your security configuration and ensuring that sensitive data is properly encrypted.

```sql+postgres
select
  name,
  arn,
  handler,
  kms_key_arn
from
  aws_lambda_function;
```

```sql+sqlite
select
  name,
  arn,
  handler,
  kms_key_arn
from
  aws_lambda_function;
```

### List of lambda functions which are not encrypted with CMK
Identify instances where AWS Lambda functions are lacking encryption with a Customer Master Key (CMK), which is crucial for enhancing data security and complying with regulatory standards.

```sql+postgres
select
  name,
  kms_key_arn
from
  aws_lambda_function
where
  kms_key_arn is null;
```

```sql+sqlite
select
  name,
  kms_key_arn
from
  aws_lambda_function
where
  kms_key_arn is null;
```

### Count of lambda functions by runtime engines
Discover the distribution of different runtime engines used in your AWS Lambda functions. This helps to understand the prevalence of different programming languages in your serverless architecture.

```sql+postgres
select
  runtime,
  count(*)
from
  aws_lambda_function
group by
  runtime;
```

```sql+sqlite
select
  runtime,
  count(*)
from
  aws_lambda_function
group by
  runtime;
```

### List of lambda function whose retention period is less than 30 days
Determine the areas in which AWS Lambda functions have a retention period of less than 30 days. This is beneficial for identifying functions that may need their retention periods adjusted to ensure data is not lost prematurely.

```sql+postgres
select
  fn.name,
  lg.name,
  lg.retention_in_days
from
  aws_lambda_function as fn
  inner join aws_cloudwatch_log_group as lg on (
    (lg.name = '/aws/lambda/')
    or (lg.name = fn.name)
  )
where
  lg.retention_in_days < 30;
```

```sql+sqlite
select
  fn.name,
  lg.name,
  lg.retention_in_days
from
  aws_lambda_function as fn
  inner join aws_cloudwatch_log_group as lg on (
    (lg.name = '/aws/lambda/')
    or (lg.name = fn.name)
  )
where
  lg.retention_in_days < 30;
```

### Availability zone count for each VPC lambda function
Determine the number of availability zones for each AWS Lambda function within a Virtual Private Cloud (VPC). This can help in understanding the distribution of your Lambda functions across different zones, which is crucial for optimizing performance and managing costs.

```sql+postgres
select
  fn.name,
  fn.region,
  count (availability_zone) as zone_count
from
  aws_lambda_function as fn
  cross join jsonb_array_elements_text(vpc_subnet_ids) as vpc_subnet
  join aws_vpc_subnet as sub on sub.subnet_id = vpc_subnet
group by
  fn.name,
  fn.region
order by
  zone_count;
```

```sql+sqlite
select
  fn.name,
  fn.region,
  count(sub.availability_zone) as zone_count
from
  aws_lambda_function as fn,
  json_each(fn.vpc_subnet_ids) as vpc_subnet
  join aws_vpc_subnet as sub on sub.subnet_id = json_extract(vpc_subnet.value, '$')
group by
  fn.name,
  fn.region
order by
  zone_count;
```

### List all the actions allowed by managed policies for a Lambda execution role
Explore which actions are permitted by managed policies for a specific Lambda execution role. This is useful for assessing the level of access a role has, and can help identify any potential security risks or areas where permissions may need to be adjusted.

```sql+postgres
select
  f.name,
  f.role,
  a.action,
  a.access_level,
  a.description
from
  aws_lambda_function as f,
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action_glob,
  glob(action_glob) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  f.role = r.arn
  and pol_arn = p.arn
  and stmt ->> 'Effect' = 'Allow'
  and f.name = 'hellopython';
```

```sql+sqlite
select
  f.name,
  f.role,
  a.action,
  a.access_level,
  a.description
from
  aws_lambda_function as f
join aws_iam_role as r on f.role = r.arn
join aws_iam_policy as p on p.arn in (select json_extract(r.attached_policy_arns, '$[*]'))
join aws_iam_action as a on a.action in (select json_extract(json_extract(p.policy_std, '$.Statement[*].Action'), '$[*]'))
where
  json_extract(p.policy_std, '$.Statement[*].Effect') = 'Allow'
  and f.name = 'hellopython';
```

### List functions not configured with a dead-letter queue
Determine the areas in which AWS Lambda functions are potentially at risk due to the absence of a configured dead-letter queue, which is crucial for handling failed asynchronous invocations and preventing data loss.

```sql+postgres
select
  arn,
  dead_letter_config_target_arn
from
  aws_lambda_function
where
  dead_letter_config_target_arn is null;
```

```sql+sqlite
select
  arn,
  dead_letter_config_target_arn
from
  aws_lambda_function
where
  dead_letter_config_target_arn is null;
```

### List runtime settings for each function
Discover the segments that have varied runtime settings for each function in your AWS Lambda service. This can help in understanding how different functions are configured and optimize them for better performance.

```sql+postgres
select
  name,
  runtime,
  handler,
  architectures
from
  aws_lambda_function;
```

```sql+sqlite
select
  name,
  runtime,
  handler,
  architectures
from
  aws_lambda_function;
```

### List URL configuration details for each function
Review the configuration for each AWS Lambda function to gain insights into their URL settings. This can help in understanding the routing and request handling setup of your serverless applications.

```sql+postgres
select
  name,
  arn,
  jsonb_pretty(url_config) as url_config
from
  aws_lambda_function;
```

```sql+sqlite
select
  name,
  arn,
  url_config
from
  aws_lambda_function;
```

### List functions that have tracing disabled
Analyze the settings to understand which AWS Lambda functions have their tracing feature disabled. This can be useful in identifying potential gaps in the monitoring and debugging process of your serverless applications.

```sql+postgres
select
  name,
  arn,
  jsonb_pretty(tracing_config) as tracing_config
from
  aws_lambda_function
where
  tracing_config ->> 'Mode' = 'PassThrough';
```

```sql+sqlite
select
  name,
  arn,
  tracing_config
from
  aws_lambda_function
where
  json_extract(tracing_config, '$.Mode') = 'PassThrough';
```