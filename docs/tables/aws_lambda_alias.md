---
title: "Steampipe Table: aws_lambda_alias - Query AWS Lambda Alias using SQL"
description: "Allows users to query AWS Lambda Alias, providing detailed information about each alias associated with AWS Lambda functions."
folder: "Lambda"
---

# Table: aws_lambda_alias - Query AWS Lambda Alias using SQL

The AWS Lambda Alias is a feature of AWS Lambda service that provides a pointer to a specific Lambda function version. It enables you to manage your function versions and routing configurations, and also allows you to shift incoming traffic between two versions of a function based on preassigned weights. This allows for a gradual code rollout and testing of new function versions in a production environment.

## Table Usage Guide

The `aws_lambda_alias` table in Steampipe provides you with information about alias resources within AWS Lambda. This table allows you, as a DevOps engineer, to query alias-specific details, including the associated function name, function version, and alias ARN. You can utilize this table to gather insights on aliases, such as the alias description, routing configuration, and revision ID. The schema outlines the various attributes of the Lambda alias for you, including the name, ARN, function version, and associated routing configuration.

## Examples

### Lambda alias basic info
Analyze the settings to understand the basic information of AWS Lambda aliases such as their names, associated function names, and function versions. This can be used to manage and track different versions of your Lambda functions.

```sql+postgres
select
  name,
  function_name,
  function_version
from
  aws_lambda_alias;
```

```sql+sqlite
select
  name,
  function_name,
  function_version
from
  aws_lambda_alias;
```

### Count of lambda alias per Lambda function
Determine the areas in which each AWS Lambda function is being used by counting the number of aliases associated with each function. This can help optimize resource usage and manage function versions effectively.

```sql+postgres
select
  function_name,
  count(function_name) count
from
  aws_lambda_alias
group by
  function_name;
```

```sql+sqlite
select
  function_name,
  count(function_name) as count
from
  aws_lambda_alias
group by
  function_name;
```

### List policy details
Explore the specifics of policies associated with AWS Lambda aliases. This is useful for understanding the permissions and configurations of your Lambda functions.

```sql+postgres
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_alias;
```

```sql+sqlite
select
  policy,
  policy_std
from
  aws_lambda_alias;
```

### List URL configuration details for each alias
This query is useful for examining the URL configurations linked to each alias in your AWS Lambda service. It allows you to understand and manage how your functions are accessed, enhancing your ability to control and optimize your serverless applications.

```sql+postgres
select
  name,
  function_name,
  jsonb_pretty(url_config) as url_config
from
  aws_lambda_alias;
```

```sql+sqlite
select
  name,
  function_name,
  url_config
from
  aws_lambda_alias;
```