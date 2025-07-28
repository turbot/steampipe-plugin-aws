---
title: "Steampipe Table: aws_bedrock_foundation_model - Query AWS Bedrock Foundation Models using SQL"
description: "Allows users to query AWS Bedrock Foundation Models, providing information about pre-trained models including their configuration, status, and associated metadata."
---

# Table: aws_bedrock_foundation_model - Query AWS Bedrock Foundation Models using SQL

AWS Bedrock Foundation Models are pre-trained models that serve as the base for various AI/ML applications. These models are provided by AWS and its partners, offering different capabilities and specializations. The service allows you to access and use these models for various AI tasks through a unified API.

## Table Usage Guide

The `aws_bedrock_foundation_model` table provides insights into foundation models within AWS Bedrock. As a data scientist or ML engineer, explore model-specific details through this table, including model capabilities, supported modalities, and lifecycle status. You can use this table to:

- Monitor the availability and status of foundation models
- Identify models that support specific input and output modalities
- Verify customization options for different models
- Assess model providers and their offerings

## Examples

### Basic info for all foundation models
```sql
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status
from
  aws_bedrock_foundation_model;
```

### List models that support text input modality
```sql
select
  model_id,
  model_name,
  provider_name
from
  aws_bedrock_foundation_model
where
  input_modalities ? 'TEXT';
```

### Get models that support customization
```sql
select
  model_id,
  model_name,
  provider_name,
  customizations_supported
from
  aws_bedrock_foundation_model
where
  customizations_supported is not null;
```

### List models by provider
```sql
select
  provider_name,
  count(*) as model_count,
  array_agg(model_name) as models
from
  aws_bedrock_foundation_model
group by
  provider_name;
```

## Fields

| Field | Type | Description |
|-------|------|-------------|
|model_id|string|The unique identifier of the foundation model.|
|model_name|string|The name of the foundation model.|
|provider_name|string|The name of the model provider.|
|input_modalities|json|The input modalities supported by the model.|
|output_modalities|json|The output modalities supported by the model.|
|customizations_supported|json|The customizations supported by the model.|
|inference_types_supported|json|The inference types supported by the model.|
|model_lifecycle_status|string|The lifecycle status of the model.|
|region|string|The AWS region in which the resource is located.|
|account_id|string|The AWS Account ID in which the resource is located.| 