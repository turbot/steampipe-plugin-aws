---
title: "Table: aws_lambda_version - Query AWS Lambda Versions using SQL"
description: "Allows users to query AWS Lambda Versions to fetch detailed information about each version of a specific AWS Lambda function."
---

# Table: aws_lambda_version - Query AWS Lambda Versions using SQL

The `aws_lambda_version` table in Steampipe provides information about each version of a specific AWS Lambda function. This table allows DevOps engineers to query version-specific details, including function name, function ARN, runtime environment, and associated metadata. Users can utilize this table to gather insights on function versions, such as code size, last modification time, version, and more. The schema outlines the various attributes of the Lambda function version, including the function ARN, version, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_version` table, you can use the `.inspect aws_lambda_version` command in Steampipe.

**Key columns**:

- `function_name`: The name of the function. This can be used to join with other tables that contain function-specific information.
- `function_arn`: The Amazon Resource Name (ARN) of the function. This can be used to join with other tables that require the ARN for specific AWS resources.
- `version`: The version of the function. This can be used to join with other tables that contain version-specific information.

## Examples

### Runtime info of each lambda version

```sql
select
  function_name,
  version,
  runtime,
  handler
from
  aws_lambda_version;
```

### List of lambda versions where code run timout is more than 2 mins

```sql
select
  function_name,
  version,
  timeout
from
  aws_lambda_version
where
  timeout :: int > 120;
```

### VPC info of each lambda version

```sql
select
  function_name,
  version,
  vpc_id,
  vpc_security_group_ids,
  vpc_subnet_ids
from
  aws_lambda_version;
```

### List policy details

```sql
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_version;
```
