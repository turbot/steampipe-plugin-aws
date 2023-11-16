---
title: "Table: aws_lambda_function - Query AWS Lambda Function using SQL"
description: "Allows users to query AWS Lambda Functions, providing information about each function's configuration, including runtime, code size, timeout, and associated tags."
---

# Table: aws_lambda_function - Query AWS Lambda Function using SQL

The `aws_lambda_function` table in Steampipe provides information about AWS Lambda Functions. This table allows DevOps engineers to query function-specific details, including the function's runtime, code size, timeout, and associated tags. Users can utilize this table to gather insights on functions, such as the function's configuration, handler, last modified date, and more. The schema outlines the various attributes of the AWS Lambda Function, including the function name, ARN, description, and associated environment variables.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_function` table, you can use the `.inspect aws_lambda_function` command in Steampipe.

**Key columns**:

- `function_name`: The name of the function. This can be used to join with other tables that reference the function by name.
- `arn`: The Amazon Resource Name (ARN) of the function. This can be used to join with other tables that reference the function by ARN.
- `runtime`: The runtime environment for the function. This is useful to identify and manage functions based on their runtime environments.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  handler,
  kms_key_arn
from
  aws_lambda_function;
```

### List of lambda functions which are not encrypted with CMK

```sql
select
  name,
  kms_key_arn
from
  aws_lambda_function
where
  kms_key_arn is null;
```

### Count of lambda functions by runtime engines

```sql
select
  runtime,
  count(*)
from
  aws_lambda_function
group by
  runtime;
```

### List of lambda function whose retention period is less than 30 days

```sql
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

```sql
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

### List all the actions allowed by managed policies for a Lambda execution role

```sql
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

### List functions not configured with a dead-letter queue

```sql
select
  arn,
  dead_letter_config_target_arn
from
  aws_lambda_function
where
  dead_letter_config_target_arn is null;
```

### List runtime settings for each function

```sql
select
  name,
  runtime,
  handler,
  architectures
from
  aws_lambda_function;
```

### List URL configuration details for each function

```sql
select
  name,
  arn,
  jsonb_pretty(url_config) as url_config
from
  aws_lambda_function;
```

### List functions that have tracing disabled

```sql
select
  name,
  arn,
  jsonb_pretty(tracing_config) as tracing_config
from
  aws_lambda_function
where
  tracing_config ->> 'Mode' = 'PassThrough';
```
