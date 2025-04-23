---
title: "Steampipe Table: aws_lambda_layer - Query AWS Lambda Layer using SQL"
description: "Allows users to query AWS Lambda Layers and retrieve information including layer ARNs, layer versions, compatible runtimes, and more."
folder: "Lambda"
---

# Table: aws_lambda_layer - Query AWS Lambda Layer using SQL

The AWS Lambda Layer is a distribution mechanism for libraries, custom runtimes, and other function dependencies. Layers promote code sharing and separation of responsibilities so that you can iterate faster on writing business logic. They allow you to manage and share your function code in a more efficient way by reducing duplication and the size of deployment packages.

## Table Usage Guide

The `aws_lambda_layer` table in Steampipe provides you with information about Lambda Layers within AWS Lambda. This table allows you, as a DevOps engineer, to query layer-specific details, including layer ARNs, layer versions, and compatible runtimes. You can utilize this table to gather insights on Lambda Layers, such as the versions available for a layer, the runtimes that a layer is compatible with, the size of the layer, and more. The schema outlines the various attributes of the Lambda Layer for you, including the layer ARN, layer version, compatible runtimes, layer size, and associated tags.

## Examples

### Basic Info
Discover the segments that have been created within your AWS Lambda layers, including their compatibility with different runtimes and architectures. This can help assess the elements within your system for potential updates or troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  compatible_runtimes,
  compatible_architectures,
  version
from
  aws_lambda_layer;
```