---
title: "Table: aws_lambda_alias - Query AWS Lambda Alias using SQL"
description: "Allows users to query AWS Lambda Alias, providing detailed information about each alias associated with AWS Lambda functions."
---

# Table: aws_lambda_alias - Query AWS Lambda Alias using SQL

The `aws_lambda_alias` table in Steampipe provides information about alias resources within AWS Lambda. This table allows DevOps engineers to query alias-specific details, including the associated function name, function version, and alias ARN. Users can utilize this table to gather insights on aliases, such as the alias description, routing configuration, and revision ID. The schema outlines the various attributes of the Lambda alias, including the name, ARN, function version, and associated routing configuration.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_alias` table, you can use the `.inspect aws_lambda_alias` command in Steampipe.

**Key columns**:

- `name`: This column contains the name of the alias. It can be used to join with other tables that contain Lambda alias names.
- `function_name`: This column contains the name of the Lambda function associated with the alias. It is useful for joining with other tables that contain Lambda function names.
- `arn`: This column contains the Amazon Resource Name (ARN) of the alias. This is a unique identifier that can be used to join with other tables that contain Lambda alias ARNs.

## Examples

### Lambda alias basic info

```sql
select
  name,
  function_name,
  function_version
from
  aws_lambda_alias;
```

### Count of lambda alias per Lambda function

```sql
select
  function_name,
  count(function_name) count
from
  aws_lambda_alias
group by
  function_name;
```

### List policy details

```sql
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_alias;
```

### List URL configuration details for each alias

```sql
select
  name,
  function_name,
  jsonb_pretty(url_config) as url_config
from
  aws_lambda_alias;
```