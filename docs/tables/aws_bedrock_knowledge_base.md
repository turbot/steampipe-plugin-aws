---
title: "Steampipe Table: aws_bedrock_knowledge_base - Query AWS Bedrock Knowledge Bases using SQL"
description: "Allows users to query AWS Bedrock Knowledge Bases, providing information about knowledge bases including their configuration, status, and associated metadata."
folder: "Bedrock"
---

# Table: aws_bedrock_knowledge_base - Query AWS Bedrock Knowledge Bases using SQL

Amazon Bedrock Knowledge Base is a feature that enables you to build retrieval augmented generation (RAG) applications by creating a container for your data. It allows you to store, manage, and retrieve information that can be used to enhance the capabilities of foundation models in Amazon Bedrock.

## Table Usage Guide

The `aws_bedrock_knowledge_base` table in Steampipe provides you with information about knowledge bases within AWS Bedrock. This table allows you, as a developer or data scientist, to query knowledge base-specific details, including configuration, status, and associated metadata. You can utilize this table to gather insights on knowledge bases, such as their storage configuration, embeddings configuration, and associated roles. The schema outlines the various attributes of the knowledge base for you, including the name, status, creation date, and associated tags.

## Examples

### Basic info
Explore the basic details of your AWS Bedrock knowledge bases to understand their current status and when they were created. This can help in managing and monitoring your knowledge base resources effectively.

```sql+postgres
select
  name,
  knowledge_base_id,
  status,
  created_at,
  arn
from
  aws_bedrock_knowledge_base;
```

```sql+sqlite
select
  name,
  knowledge_base_id,
  status,
  created_at,
  arn
from
  aws_bedrock_knowledge_base;
```

### List knowledge bases that are not active
Identify knowledge bases that may require attention by finding those that are not in an active state. This can help in troubleshooting or maintenance activities.

```sql+postgres
select
  name,
  knowledge_base_id,
  status,
  created_at,
  failure_reasons
from
  aws_bedrock_knowledge_base
where
  status != 'ACTIVE';
```

```sql+sqlite
select
  name,
  knowledge_base_id,
  status,
  created_at,
  failure_reasons
from
  aws_bedrock_knowledge_base
where
  status != 'ACTIVE';
```

### Get knowledge base configuration details
Analyze the configuration settings of your knowledge bases to understand their embeddings setup and other technical specifications.

```sql+postgres
select
  name,
  knowledge_base_id,
  knowledge_base_configuration,
  storage_configuration
from
  aws_bedrock_knowledge_base;
```

```sql+sqlite
select
  name,
  knowledge_base_id,
  knowledge_base_configuration,
  storage_configuration
from
  aws_bedrock_knowledge_base;
```

### List knowledge bases with their associated IAM roles
Examine the IAM roles associated with your knowledge bases to understand their permissions and access configurations.

```sql+postgres
select
  name,
  knowledge_base_id,
  role_arn,
  status
from
  aws_bedrock_knowledge_base;
```

```sql+sqlite
select
  name,
  knowledge_base_id,
  role_arn,
  status
from
  aws_bedrock_knowledge_base;
```