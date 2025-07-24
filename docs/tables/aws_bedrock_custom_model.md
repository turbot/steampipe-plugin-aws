---
title: "Steampipe Table: aws_bedrock_custom_model - Query AWS Bedrock Custom Models using SQL"
description: "Allows users to query AWS Bedrock Custom Models, providing information about custom foundation models including their configuration, status, and associated metadata."
folder: "Bedrock"
---

# Table: aws_bedrock_custom_model - Query AWS Bedrock Custom Models using SQL

AWS Bedrock Custom Models allows you to customize foundation models for your specific use cases. These custom models are built on top of base models and can be fine-tuned or pre-trained to better suit your specific requirements. The service provides a managed environment for customizing, deploying, and managing these models.

## Table Usage Guide

The `aws_bedrock_custom_model` table in Steampipe provides you with information about custom models within AWS Bedrock. This table allows you, as a data scientist or ML engineer, to query custom model-specific details, including model configuration, status, and associated metadata. You can utilize this table to gather insights on custom models, such as their base models, customization types, creation times, and more. The schema outlines the various attributes of the custom model for you, including the model ARN, base model name, creation time, and associated tags.

## Examples

### Basic info
Explore the basic details of your custom Bedrock models to understand their configuration and status. This can help in managing and monitoring your AI/ML resources effectively.

```sql+postgres
select
  model_name,
  model_arn,
  base_model_name,
  customization_type,
  model_status
from
  aws_bedrock_custom_model;
```

```sql+sqlite
select
  model_name,
  model_arn,
  base_model_name,
  customization_type,
  model_status
from
  aws_bedrock_custom_model;
```

### List custom models created in the last 30 days
Identify recently created custom models to track new model deployments and monitor resource usage.

```sql+postgres
select
  model_name,
  model_arn,
  base_model_name,
  creation_time
from
  aws_bedrock_custom_model
where
  creation_time >= (current_timestamp - interval '30' day)
order by
  creation_time;
```

```sql+sqlite
select
  model_name,
  model_arn,
  base_model_name,
  creation_time
from
  aws_bedrock_custom_model
where
  creation_time >= datetime('now','-30 day')
order by
  creation_time;
```

### List custom models by base model type
Analyze the distribution of custom models across different base models to understand your model customization patterns.

```sql+postgres
select
  base_model_name,
  count(*) as model_count,
  array_agg(model_name) as custom_models
from
  aws_bedrock_custom_model
group by
  base_model_name;
```

```sql+sqlite
select
  base_model_name,
  count(*) as model_count,
  group_concat(model_name) as custom_models
from
  aws_bedrock_custom_model
group by
  base_model_name;
```

### Get details of a specific custom model
Retrieve detailed information about a particular custom model using its ARN for in-depth analysis or troubleshooting.

```sql+postgres
select
  model_name,
  model_arn,
  base_model_name,
  customization_type,
  model_status,
  creation_time,
  owner_account_id
from
  aws_bedrock_custom_model
where
  model_arn = 'arn:aws:bedrock:us-east-1:123456789012:custom-model/example-model';
```

```sql+sqlite
select
  model_name,
  model_arn,
  base_model_name,
  customization_type,
  model_status,
  creation_time,
  owner_account_id
from
  aws_bedrock_custom_model
where
  model_arn = 'arn:aws:bedrock:us-east-1:123456789012:custom-model/example-model';
``` 