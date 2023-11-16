---
title: "Table: aws_sagemaker_model - Query AWS SageMaker Models using SQL"
description: "Allows users to query AWS SageMaker Models for detailed information about each model, including its name, ARN, creation time, execution role, and more."
---

# Table: aws_sagemaker_model - Query AWS SageMaker Models using SQL

The `aws_sagemaker_model` table in Steampipe provides information about models within AWS SageMaker. This table allows DevOps engineers to query model-specific details, including the model name, ARN, creation time, execution role, and more. Users can utilize this table to gather insights on models, such as their associated containers, data input configurations, and VPC configurations. The schema outlines the various attributes of the SageMaker model, including the model ARN, creation time, model name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_model` table, you can use the `.inspect aws_sagemaker_model` command in Steampipe.

### Key columns:

- `model_name`: The name of the model. This can be used to join this table with other tables that contain model-specific information.
- `model_arn`: The Amazon Resource Name (ARN) of the model. This can be used to join this table with other tables that contain resource-specific information.
- `creation_time`: The time at which the model was created. This can be used to track the age and usage of models over time.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  enable_network_isolation
from
  aws_sagemaker_model;
```

### List network isolated models

```sql
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
