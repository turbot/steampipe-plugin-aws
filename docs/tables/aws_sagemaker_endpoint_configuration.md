---
title: "Steampipe Table: aws_sagemaker_endpoint_configuration - Query AWS SageMaker Endpoint Configurations using SQL"
description: "Allows users to query AWS SageMaker Endpoint Configurations to retrieve detailed information about each endpoint configuration in the AWS SageMaker service."
folder: "Config"
---

# Table: aws_sagemaker_endpoint_configuration - Query AWS SageMaker Endpoint Configurations using SQL

The AWS SageMaker Endpoint Configuration is a feature of Amazon SageMaker, a fully managed service that provides developers and data scientists with the ability to build, train, and deploy machine learning (ML) models quickly. It allows you to define the settings for deploying your model in SageMaker, including the ML compute instances to deploy and the initial variant weights. This configuration is used when you create a SageMaker model endpoint, enabling real-time inferences from the model.

## Table Usage Guide

The `aws_sagemaker_endpoint_configuration` table in Steampipe provides you with information about endpoint configurations within AWS SageMaker. This table enables you, as a data scientist, machine learning engineer, or DevOps professional, to query endpoint configuration specific details, including the Amazon Resource Name (ARN), creation time, endpoint configuration name, and production variants. You can utilize this table to gather insights on endpoint configurations, such as the instances count per variant, instance type per variant, variant name, and more. The schema outlines the various attributes of the SageMaker endpoint configuration that are available to you, including the endpoint configuration ARN, creation time, endpoint configuration name, and the production variants.

## Examples

### Basic info
Explore the various configurations of your AWS SageMaker endpoints to gain insights into their setup, including details like creation time and associated tags. This can help optimize resource usage and improve management of machine learning models.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which SageMaker endpoint configurations are not encrypted, which could potentially expose sensitive data. This query is useful in identifying potential security risks within your AWS environment.

```sql+postgres
select
  name,
  arn,
  kms_key_id
from
  aws_sagemaker_endpoint_configuration
where
  kms_key_id is null;
```

```sql+sqlite
select
  name,
  arn,
  kms_key_id
from
  aws_sagemaker_endpoint_configuration
where
  kms_key_id is null;
```