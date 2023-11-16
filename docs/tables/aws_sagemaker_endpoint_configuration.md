---
title: "Table: aws_sagemaker_endpoint_configuration - Query AWS SageMaker Endpoint Configurations using SQL"
description: "Allows users to query AWS SageMaker Endpoint Configurations to retrieve detailed information about each endpoint configuration in the AWS SageMaker service."
---

# Table: aws_sagemaker_endpoint_configuration - Query AWS SageMaker Endpoint Configurations using SQL

The `aws_sagemaker_endpoint_configuration` table in Steampipe provides information about endpoint configurations within AWS SageMaker. This table allows data scientists, machine learning engineers, and DevOps professionals to query endpoint configuration specific details, including the Amazon Resource Name (ARN), creation time, endpoint configuration name, and production variants. Users can utilize this table to gather insights on endpoint configurations, such as the instances count per variant, instance type per variant, variant name, and more. The schema outlines the various attributes of the SageMaker endpoint configuration, including the endpoint configuration ARN, creation time, endpoint configuration name, and the production variants.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_endpoint_configuration` table, you can use the `.inspect aws_sagemaker_endpoint_configuration` command in Steampipe.

### Key columns:

- `endpoint_config_name`: The name of the endpoint configuration. This column can be used to join with other tables that contain information about SageMaker endpoint configurations.
- `endpoint_config_arn`: The Amazon Resource Name (ARN) of the endpoint configuration. This is a unique identifier that can be used to join with other tables that require the ARN of the SageMaker endpoint configuration.
- `creation_time`: The time when the endpoint configuration was created. This column can be useful for tracking the lifecycle of endpoint configurations over time.

## Examples

### Basic info

```sql
select
  name,
  arn,
  kms_key_id,
  creation_time,
  production_variants,
  tags
from
  aws_sagemaker_endpoint_configuration;
```

### List unencrypted endpoint configurations

```sql
select
  name,
  arn,
  kms_key_id
from
  aws_sagemaker_endpoint_configuration
where
  kms_key_id is null;
```
