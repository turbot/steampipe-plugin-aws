---
title: "Steampipe Table: aws_bedrock_foundation_model - Query AWS Bedrock Foundation Models using SQL"
description: "Allows users to query AWS Bedrock Foundation Models, providing information about pre-trained models including their configuration, status, and associated metadata."
folder: "Bedrock"
---

# Table: aws_bedrock_foundation_model

AWS Bedrock Foundation Models are pre-trained models that serve as the base for various AI/ML applications. These models are provided by AWS and its partners, offering different capabilities and specializations. The service allows you to access and use these models for various AI tasks through a unified API.

## Table Usage Guide

The `aws_bedrock_foundation_model` table provides insights into foundation models within AWS Bedrock. As a data scientist or ML engineer, explore model-specific details through this table, including model capabilities, supported modalities, and lifecycle status. Utilize it to monitor model availability, verify supported features, and ensure appropriate model selection for your AI/ML workloads.

## Examples

### Basic info
Get a quick overview of all AWS Bedrock foundation models, including their IDs, names, providers, and status. This is useful for inventory and monitoring purposes.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status,
  customizations_supported
from
  aws_bedrock_foundation_model;
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status,
  customizations_supported
from
  aws_bedrock_foundation_model;
```

### List models that support text input modality
Identify which foundation models can process text inputs. This helps you select appropriate models for text-based applications.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status
from
  aws_bedrock_foundation_model
where
  input_modalities ? 'TEXT';
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status
from
  aws_bedrock_foundation_model
where
  json_extract(input_modalities, '$.TEXT') is not null;
```

### Get models that support customization
Discover which foundation models can be customized to meet specific requirements. This is useful for identifying models that can be fine-tuned for your use case.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  customizations_supported,
  model_lifecycle_status
from
  aws_bedrock_foundation_model
where
  customizations_supported is not null;
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  customizations_supported,
  model_lifecycle_status
from
  aws_bedrock_foundation_model
where
  customizations_supported is not null;
```

### List models by provider
Analyze the distribution of models across different providers. This helps you understand the available options from each provider.

```sql+postgres
select
  provider_name,
  count(*) as model_count,
  array_agg(model_name) as models
from
  aws_bedrock_foundation_model
group by
  provider_name
order by
  model_count desc;
```

```sql+sqlite
select
  provider_name,
  count(*) as model_count,
  group_concat(model_name, ', ') as models
from
  aws_bedrock_foundation_model
group by
  provider_name
order by
  model_count desc;
```

### List active foundation models
Identify all foundation models that are currently active and available for use. This helps ensure you're only using models that are ready for production.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  customizations_supported,
  input_modalities,
  output_modalities
from
  aws_bedrock_foundation_model
where
  model_lifecycle_status = 'ACTIVE';
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  customizations_supported,
  input_modalities,
  output_modalities
from
  aws_bedrock_foundation_model
where
  model_lifecycle_status = 'ACTIVE';
```

### Get details for a specific model
Retrieve comprehensive information about a particular foundation model. This is useful when you need to understand all capabilities and configurations of a specific model.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status,
  customizations_supported,
  input_modalities,
  output_modalities,
  response_streaming_supported,
  inference_types_supported
from
  aws_bedrock_foundation_model
where
  model_id = 'anthropic.claude-v2';
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  model_lifecycle_status,
  customizations_supported,
  input_modalities,
  output_modalities,
  response_streaming_supported,
  inference_types_supported
from
  aws_bedrock_foundation_model
where
  model_id = 'anthropic.claude-v2';
```

### List models with specific inference types
Find models that support specific inference types to ensure compatibility with your application requirements.

```sql+postgres
select
  model_id,
  model_name,
  provider_name,
  inference_types_supported
from
  aws_bedrock_foundation_model
where
  inference_types_supported ? 'ON_DEMAND';
```

```sql+sqlite
select
  model_id,
  model_name,
  provider_name,
  inference_types_supported
from
  aws_bedrock_foundation_model
where
  json_extract(inference_types_supported, '$.ON_DEMAND') is not null;
```