---
title: "Table: aws_lambda_layer - Query AWS Lambda Layer using SQL"
description: "Allows users to query AWS Lambda Layers and retrieve information including layer ARNs, layer versions, compatible runtimes, and more."
---

# Table: aws_lambda_layer - Query AWS Lambda Layer using SQL

The `aws_lambda_layer` table in Steampipe provides information about Lambda Layers within AWS Lambda. This table allows DevOps engineers to query layer-specific details, including layer ARNs, layer versions, and compatible runtimes. Users can utilize this table to gather insights on Lambda Layers, such as the versions available for a layer, the runtimes that a layer is compatible with, the size of the layer, and more. The schema outlines the various attributes of the Lambda Layer, including the layer ARN, layer version, compatible runtimes, layer size, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_layer` table, you can use the `.inspect aws_lambda_layer` command in Steampipe.

* Key columns:
    * `layer_name` - The name of the lambda layer. This can be used to join this table with other tables that contain lambda layer names.
    * `layer_arn` - The Amazon Resource Name (ARN) of the Lambda layer. This can be used to join this table with other tables that contain lambda layer ARNs.
    * `layer_version` - The version number of the layer. This is important as different versions of a layer may have different attributes or compatible runtimes.

## Examples

### Basic Info

```sql
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  jsonb_pretty(compatible_runtimes) as compatible_runtimes,
  jsonb_pretty(compatible_architectures) as compatible_architectures,
  version
from
  aws_lambda_layer;
```
