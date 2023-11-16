---
title: "Table: aws_lambda_layer_version - Query AWS Lambda Layer Versions using SQL"
description: "Allows users to query AWS Lambda Layer Versions, providing detailed information about each layer version, including its ARN, description, license info, compatible runtimes, created date, and more."
---

# Table: aws_lambda_layer_version - Query AWS Lambda Layer Versions using SQL

The `aws_lambda_layer_version` table in Steampipe provides information about each version of a Lambda layer within AWS Lambda. This table allows DevOps engineers to query version-specific details, including the layer ARN, description, license info, compatible runtimes, and the date it was created. Users can utilize this table to gather insights on layer versions, such as their compatibility with different runtimes, license information, and more. The schema outlines the various attributes of the Lambda layer version, including the layer version ARN, layer name, version, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_layer_version` table, you can use the `.inspect aws_lambda_layer_version` command in Steampipe.

**Key columns**:

- `layer_name`: The name of the layer. This is useful for joining with other tables that may reference the layer by name.
- `layer_version_arn`: The Amazon Resource Name (ARN) of the layer version. This unique identifier is useful for joining with other tables that may reference the layer version by its ARN.
- `version`: The version number of the layer. This is useful for tracking the progression of layer versions and for joining with other tables that may reference the layer version by its version number.

## Examples

### Basic Info

```sql
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
