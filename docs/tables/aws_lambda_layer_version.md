---
title: "Steampipe Table: aws_lambda_layer_version - Query AWS Lambda Layer Versions using SQL"
description: "Allows users to query AWS Lambda Layer Versions, providing detailed information about each layer version, including its ARN, description, license info, compatible runtimes, created date, and more."
folder: "Lambda"
---

# Table: aws_lambda_layer_version - Query AWS Lambda Layer Versions using SQL

The AWS Lambda Layer Version is a distribution mechanism for libraries, custom runtimes, and other function dependencies. Layers promote code sharing and separation of responsibilities so that you can iterate faster on writing business logic. With layers, you can manage your in-house and external code dependencies separately.

## Table Usage Guide

The `aws_lambda_layer_version` table in Steampipe provides you with information about each version of a Lambda layer within AWS Lambda. This table allows you, as a DevOps engineer, to query version-specific details, including the layer ARN, description, license info, compatible runtimes, and the date it was created. You can utilize this table to gather insights on layer versions, such as their compatibility with different runtimes, license information, and more. The schema outlines the various attributes of the Lambda layer version for you, including the layer version ARN, layer name, version, description, and associated tags.

## Examples

### Basic Info
Explore which AWS Lambda layer versions are available and when they were created to better manage and update your serverless applications. This can help you ensure your applications are always running on the latest and most secure versions.

```sql+postgres
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std,
  version
from
  aws_lambda_layer_version;
```

```sql+sqlite
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  policy,
  policy_std,
  version
from
  aws_lambda_layer_version;
```