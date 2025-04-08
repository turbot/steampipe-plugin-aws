---
title: "Steampipe Table: aws_lambda_version - Query AWS Lambda Versions using SQL"
description: "Allows users to query AWS Lambda Versions to fetch detailed information about each version of a specific AWS Lambda function."
folder: "Lambda"
---

# Table: aws_lambda_version - Query AWS Lambda Versions using SQL

The AWS Lambda Version is a distinct AWS Lambda function configuration that includes the function code and settings. Each Lambda function version has a unique Amazon Resource Name (ARN) and is immutable once it is published. It allows you to manage, invoke, or traffic shift between different versions of a function.

## Table Usage Guide

The `aws_lambda_version` table in Steampipe provides you with information about each version of a specific AWS Lambda function. This table allows you, as a DevOps engineer, to query version-specific details, including function name, function ARN, runtime environment, and associated metadata. You can utilize this table to gather insights on function versions, such as code size, last modification time, version, and more. The schema outlines the various attributes of the Lambda function version for you, including the function ARN, version, description, and associated tags.

## Examples

### Runtime info of each lambda version
Identify instances where specific versions of Lambda functions are running to understand their operational status. This allows for efficient management and potential optimization of cloud resources.

```sql+postgres
select
  function_name,
  version,
  runtime,
  handler
from
  aws_lambda_version;
```

```sql+sqlite
select
  function_name,
  version,
  runtime,
  handler
from
  aws_lambda_version;
```

### List of lambda versions where code run timout is more than 2 mins
Identify instances where AWS Lambda versions have a code run timeout exceeding 2 minutes. This is useful to pinpoint potential performance issues and optimize your Lambda functions for efficient execution.

```sql+postgres
select
  function_name,
  version,
  timeout
from
  aws_lambda_version
where
  timeout :: int > 120;
```

```sql+sqlite
select
  function_name,
  version,
  timeout
from
  aws_lambda_version
where
  cast(timeout as integer) > 120;
```

### VPC info of each lambda version
Explore the network configurations of different versions of AWS Lambda functions. This can be used to understand how these functions are interacting with your Virtual Private Cloud (VPC) infrastructure, helping to identify potential security or connectivity issues.

```sql+postgres
select
  function_name,
  version,
  vpc_id,
  vpc_security_group_ids,
  vpc_subnet_ids
from
  aws_lambda_version;
```

```sql+sqlite
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
Gain insights into the specifics of your AWS Lambda version policies. This query is beneficial for understanding and reviewing your policy configurations in a user-friendly format.

```sql+postgres
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_version;
```

```sql+sqlite
select
  policy,
  policy_std
from
  aws_lambda_version;
```