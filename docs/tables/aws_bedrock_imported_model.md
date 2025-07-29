---
title: "Steampipe Table: aws_bedrock_imported_model - Query AWS Bedrock Imported Models using SQL"
description: "Allows users to query AWS Bedrock Imported Models, providing information about models that have been imported into Amazon Bedrock."
---

# Table: aws_bedrock_imported_model - Query AWS Bedrock Imported Models using SQL

Amazon Bedrock Imported Models are custom models that have been imported into the Amazon Bedrock service. These models can be used for various AI/ML tasks through Amazon Bedrock's unified API. The service allows you to import and manage your own models while leveraging AWS's infrastructure and scalability.

## Table Usage Guide

The `aws_bedrock_imported_model` table provides insights into imported models within AWS Bedrock. As a data scientist or ML engineer, explore model-specific details through this table, including model names, architectures, and associated metadata. You can use this table to:

- Monitor and track imported models
- Verify model configurations and settings
- Audit model creation times and support features
- Analyze model architectures and capabilities

## Examples

### Basic info
Explore the basic information about imported models in AWS Bedrock to understand their configuration and capabilities.

```sql
select
  model_name,
  model_arn,
  creation_time,
  instruct_supported,
  model_architecture
from
  aws_bedrock_imported_model;
```

### List models created in the last 7 days
Analyze recently imported models to track new additions and monitor recent model imports.

```sql
select
  model_name,
  model_arn,
  creation_time,
  instruct_supported
from
  aws_bedrock_imported_model
where
  creation_time >= now() - interval '7 days';
```

### List models that support instructions
Identify which imported models have instruction support capabilities to better understand their potential applications.

```sql
select
  model_name,
  model_arn,
  creation_time,
  model_architecture
from
  aws_bedrock_imported_model
where
  instruct_supported = true;
```

### Get details for a specific model by ARN
Retrieve comprehensive information about a particular model using its unique identifier (ARN).

```sql
select
  model_name,
  model_arn,
  creation_time,
  instruct_supported,
  model_architecture
from
  aws_bedrock_imported_model
where
  model_arn = 'arn:aws:bedrock:us-east-1:123456789012:imported-model/example-model';
``` 