---
title: "Steampipe Table: aws_sagemaker_model - Query AWS SageMaker Models using SQL"
description: "Allows users to query AWS SageMaker Models for detailed information about each model, including its name, ARN, creation time, execution role, and more."
folder: "SageMaker"
---

# Table: aws_sagemaker_model - Query AWS SageMaker Models using SQL

An AWS SageMaker Model in Amazon SageMaker represents the Amazon S3 location where model artifacts are stored, and the Docker registry path where the image that contains the inference code is stored. These models are immutable and can be used for multiple purposes such as predictions, transformations, and associations. SageMaker model provides the entry point for services to access the model artifacts and image.

## Table Usage Guide

The `aws_sagemaker_model` table in Steampipe provides you with information about models within AWS SageMaker. This table allows you, as a DevOps engineer, to query model-specific details, including the model name, ARN, creation time, execution role, and more. You can utilize this table to gather insights on models, such as their associated containers, data input configurations, and VPC configurations. The schema outlines the various attributes of the SageMaker model for you, including the model ARN, creation time, model name, and associated tags.

## Examples

### Basic info
Explore the settings of the AWS SageMaker model to understand its network isolation status and the time it was created. This can help in auditing and managing your machine learning models effectively.

```sql+postgres
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model;
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model;
```

### List network isolated models
Determine the areas in which network isolation is enabled within SageMaker models. This is useful for ensuring security and data privacy by preventing any unnecessary network access to these models.

```sql+postgres
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model
where
  enable_network_isolation;
```

```sql+sqlite
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model
where
  enable_network_isolation = 1;
```