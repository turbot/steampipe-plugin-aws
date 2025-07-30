---
title: "Steampipe Table: aws_bedrock_imported_model - Query AWS Bedrock Imported Models using SQL"
description: "Allows users to query AWS Bedrock Imported Models, providing information about models that have been imported into Amazon Bedrock."
folder: "Bedrock"
---

# Table: aws_bedrock_imported_model

Amazon Bedrock Imported Models are custom models that have been imported into the Amazon Bedrock service. These models can be used for various AI/ML tasks through Amazon Bedrock's unified API. The service allows you to import and manage your own models while leveraging AWS's infrastructure and scalability.

## Table Usage Guide

The `aws_bedrock_imported_model` table provides insights into imported models within AWS Bedrock. As a data scientist or ML engineer, explore model-specific details through this table, including model names, architectures, and associated metadata. Utilize it to monitor imported models, verify configurations, and ensure appropriate model deployment and management.

## Examples

### Basic info
Get a quick overview of all imported models in AWS Bedrock, including their names, ARNs, and key configurations. This is useful for inventory and monitoring purposes.

```sql+postgres
select
  model_name,
  arn,
  creation_time,
  instruct_supported,
  model_architecture,
  job_arn,
  job_name,
  model_data_source,
  model_kms_key_arn,
  custom_model_units
from
  aws_bedrock_imported_model;
```

```sql+sqlite
select
  model_name,
  arn,
  creation_time,
  instruct_supported,
  model_architecture,
  job_arn,
  job_name,
  model_data_source,
  model_kms_key_arn,
  custom_model_units
from
  aws_bedrock_imported_model;
```

### List models created in the last 7 days
Identify recently imported models to track new additions and monitor recent model imports. This helps in maintaining an up-to-date inventory of your models.

```sql+postgres
select
  model_name,
  model_arn,
  creation_time,
  instruct_supported,
  model_architecture
from
  aws_bedrock_imported_model
where
  creation_time >= now() - interval '7 days';
```

```sql+sqlite
select
  model_name,
  model_arn,
  creation_time,
  instruct_supported,
  model_architecture
from
  aws_bedrock_imported_model
where
  creation_time >= datetime('now', '-7 days');
```

### List models that support instructions
Identify which imported models have instruction support capabilities. This helps in selecting models suitable for instruction-based tasks.

```sql+postgres
select
  model_name,
  model_arn,
  creation_time,
  model_architecture,
  tags
from
  aws_bedrock_imported_model
where
  instruct_supported = true;
```

```sql+sqlite
select
  model_name,
  model_arn,
  creation_time,
  model_architecture,
  tags
from
  aws_bedrock_imported_model
where
  instruct_supported = 1;
```

### Get details for a specific model by ARN
Retrieve comprehensive information about a particular model. This is useful when you need to understand all aspects of a specific imported model.

```sql+postgres
select
  model_name,
  arn,
  creation_time,
  instruct_supported,
  model_architecture,
  job_arn,
  job_name,
  model_data_source,
  model_kms_key_arn,
  custom_model_units
from
  aws_bedrock_imported_model
where
  arn = 'arn:aws:bedrock:us-east-1:123456789012:imported-model/example-model';
```

```sql+sqlite
select
  model_name,
  arn,
  creation_time,
  instruct_supported,
  model_architecture,
  job_arn,
  job_name,
  model_data_source,
  model_kms_key_arn,
  custom_model_units
from
  aws_bedrock_imported_model
where
  arn = 'arn:aws:bedrock:us-east-1:123456789012:imported-model/example-model';
```

### List models by architecture type
Group imported models by their architecture to understand the distribution of different model types in your environment.

```sql+postgres
select
  model_architecture,
  count(*) as model_count,
  array_agg(model_name) as models
from
  aws_bedrock_imported_model
group by
  model_architecture
order by
  model_count desc;
```

```sql+sqlite
select
  model_architecture,
  count(*) as model_count,
  group_concat(model_name, ', ') as models
from
  aws_bedrock_imported_model
group by
  model_architecture
order by
  model_count desc;
```
