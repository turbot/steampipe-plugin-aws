---
title: "Steampipe Table: aws_bedrock_guardrail - Query AWS Bedrock Guardrails using SQL"
description: "Allows users to query AWS Bedrock Guardrails, providing information about content filtering policies, security configurations, and model protection settings."
folder: "Bedrock"
---

# Table: aws_bedrock_guardrail - Query AWS Bedrock Guardrails using SQL

AWS Bedrock Guardrails enable you to implement safeguards for your generative AI applications based on your use cases and responsible AI policies. Guardrails help control the interaction between users and foundation models by filtering undesirable content, redacting personally identifiable information (PII), and blocking topics that are not appropriate for your application.

## Table Usage Guide

The `aws_bedrock_guardrail` table in Steampipe provides you with information about guardrails within AWS Bedrock. This table allows you, as a developer, ML engineer, or security administrator, to query guardrail-specific details, including content policies, sensitive information policies, topic policies, and associated metadata. You can utilize this table to gather insights on guardrails, such as their configuration, status, blocked messaging, and policy settings. The schema outlines the various attributes of the guardrail for you, including the guardrail ID, ARN, version, and various policy configurations.

## Examples

### Basic info
Explore the basic details of your Bedrock guardrails to understand their configuration and status. This can help in managing and monitoring your content filtering policies effectively.

```sql+postgres
select
  guardrail_id,
  name,
  status,
  version,
  description,
  created_at,
  updated_at
from
  aws_bedrock_guardrail;
```

```sql+sqlite
select
  guardrail_id,
  name,
  status,
  version,
  description,
  created_at,
  updated_at
from
  aws_bedrock_guardrail;
```

### Get details of a specific guardrail
Retrieve comprehensive information about a particular guardrail using its ARN for in-depth analysis or troubleshooting.

```sql+postgres
select
  name,
  guardrail_id,
  arn,
  status,
  version,
  blocked_input_messaging,
  blocked_outputs_messaging,
  kms_key_arn,
  created_at,
  updated_at
from
  aws_bedrock_guardrail
where
  arn = 'arn:aws:bedrock:us-east-1:123456789012:guardrail/example-id';
```

```sql+sqlite
select
  name,
  guardrail_id,
  arn,
  status,
  version,
  blocked_input_messaging,
  blocked_outputs_messaging,
  kms_key_arn,
  created_at,
  updated_at
from
  aws_bedrock_guardrail
where
  arn = 'arn:aws:bedrock:us-east-1:123456789012:guardrail/example-id';
```

### List guardrails with their content policies
Analyze the content policy configurations of your guardrails to ensure appropriate content filtering is in place.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  content_policy
from
  aws_bedrock_guardrail
where
  content_policy is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  content_policy
from
  aws_bedrock_guardrail
where
  content_policy is not null;
```

### List guardrails with sensitive information policies
Identify guardrails that have policies configured to handle sensitive information such as PII.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  sensitive_information_policy
from
  aws_bedrock_guardrail
where
  sensitive_information_policy is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  sensitive_information_policy
from
  aws_bedrock_guardrail
where
  sensitive_information_policy is not null;
```

### List guardrails with topic blocking policies
Find guardrails that have specific topic-based filtering configured to block certain types of conversations.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  topic_policy
from
  aws_bedrock_guardrail
where
  topic_policy is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  topic_policy
from
  aws_bedrock_guardrail
where
  topic_policy is not null;
```

### List guardrails with word filtering policies
Examine guardrails that have word-based filtering configured to block specific words or phrases.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  word_policy
from
  aws_bedrock_guardrail
where
  word_policy is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  word_policy
from
  aws_bedrock_guardrail
where
  word_policy is not null;
```

### List guardrails that are in READY status
Monitor guardrails that are active and ready to be used in your applications.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  version,
  updated_at
from
  aws_bedrock_guardrail
where
  status = 'READY';
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  version,
  updated_at
from
  aws_bedrock_guardrail
where
  status = 'READY';
```

### List guardrails created in the last 30 days
Monitor recently created guardrails to track new deployments and changes in your AWS Bedrock environment.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  created_at,
  description
from
  aws_bedrock_guardrail
where
  created_at >= now() - interval '30 days';
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  created_at,
  description
from
  aws_bedrock_guardrail
where
  created_at >= datetime('now', '-30 days');
```

### List guardrails with encryption enabled
Identify guardrails that have KMS encryption configured for additional security.

```sql+postgres
select
  name,
  guardrail_id,
  kms_key_arn,
  status
from
  aws_bedrock_guardrail
where
  kms_key_arn is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  kms_key_arn,
  status
from
  aws_bedrock_guardrail
where
  kms_key_arn is not null;
```

### List guardrails with contextual grounding policies
Find guardrails that have contextual grounding policies configured to ensure responses are grounded in provided context.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  contextual_grounding_policy
from
  aws_bedrock_guardrail
where
  contextual_grounding_policy is not null;
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  contextual_grounding_policy
from
  aws_bedrock_guardrail
where
  contextual_grounding_policy is not null;
```

### List guardrails with failure information
Identify guardrails that have failed and review recommendations for resolution.

```sql+postgres
select
  name,
  guardrail_id,
  status,
  status_reasons,
  failure_recommendations
from
  aws_bedrock_guardrail
where
  status = 'FAILED';
```

```sql+sqlite
select
  name,
  guardrail_id,
  status,
  status_reasons,
  failure_recommendations
from
  aws_bedrock_guardrail
where
  status = 'FAILED';
```

### Get guardrail custom blocked messages
View the custom messages that are returned when content is blocked by guardrails.

```sql+postgres
select
  name,
  guardrail_id,
  blocked_input_messaging,
  blocked_outputs_messaging
from
  aws_bedrock_guardrail;
```

```sql+sqlite
select
  name,
  guardrail_id,
  blocked_input_messaging,
  blocked_outputs_messaging
from
  aws_bedrock_guardrail;
```

