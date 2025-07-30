---
title: "Steampipe Table: aws_bedrock_agent - Query AWS Bedrock Agents using SQL"
description: "Allows users to query AWS Bedrock Agents, providing information about AI-powered assistants including their configuration, status, and associated metadata."
folder: "Bedrock"
---

# Table: aws_bedrock_agent - Query AWS Bedrock Agents using SQL

AWS Bedrock Agents are AI-powered assistants that can be customized to perform specific tasks using foundation models. These agents can be configured with specific instructions, foundation models, and security settings to handle various use cases. The service provides a managed environment for creating, deploying, and managing these agents.

## Table Usage Guide

The `aws_bedrock_agent` table in Steampipe provides you with information about agents within AWS Bedrock. This table allows you, as a developer or ML engineer, to query agent-specific details, including configuration, status, and associated metadata. You can utilize this table to gather insights on agents, such as their foundation models, instructions, session settings, and more. The schema outlines the various attributes of the agent for you, including the agent ID, name, status, and associated IAM roles.

## Examples

### Basic info
Explore the basic details of your Bedrock agents to understand their configuration and status. This can help in managing and monitoring your AI assistants effectively.

```sql+postgres
select
  agent_id,
  agent_name,
  agent_status,
  description,
  updated_at
from
  aws_bedrock_agent;
```

```sql+sqlite
select
  agent_id,
  agent_name,
  agent_status,
  description,
  updated_at
from
  aws_bedrock_agent;
```

### List agents with their foundation models and instructions
Analyze the configuration of your agents to understand their underlying models and behavior instructions.

```sql+postgres
select
  agent_name,
  foundation_model,
  instruction,
  agent_status
from
  aws_bedrock_agent;
```

```sql+sqlite
select
  agent_name,
  foundation_model,
  instruction,
  agent_status
from
  aws_bedrock_agent;
```

### Get details of a specific agent
Retrieve comprehensive information about a particular agent using its ID for in-depth analysis or troubleshooting.

```sql+postgres
select
  agent_name,
  agent_id,
  arn,
  agent_status,
  foundation_model,
  instruction,
  created_at,
  updated_at,
  description
from
  aws_bedrock_agent
where
  agent_id = 'example-agent-id';
```

```sql+sqlite
select
  agent_name,
  agent_id,
  arn,
  agent_status,
  foundation_model,
  instruction,
  created_at,
  updated_at,
  description
from
  aws_bedrock_agent
where
  agent_id = 'example-agent-id';
```

### List agents with their session timeout settings
Monitor agent session configurations to ensure appropriate timeout settings for user interactions.

```sql+postgres
select
  agent_name,
  idle_session_ttl_in_seconds,
  agent_status
from
  aws_bedrock_agent;
```

```sql+sqlite
select
  agent_name,
  idle_session_ttl_in_seconds,
  agent_status
from
  aws_bedrock_agent;
```

### List agents with their IAM roles
Examine the IAM role configurations of your agents to audit security and permissions.

```sql+postgres
select
  agent_name,
  agent_resource_role_arn,
  agent_status
from
  aws_bedrock_agent;
```

```sql+sqlite
select
  agent_name,
  agent_resource_role_arn,
  agent_status
from
  aws_bedrock_agent;
``` 

### List agents created in the last 30 days
Monitor recently created agents to track new deployments and changes in your AWS Bedrock environment.

```sql+postgres
select
  agent_name,
  agent_id,
  agent_status,
  created_at,
  foundation_model,
  description
from
  aws_bedrock_agent
where
  created_at >= now() - interval '30 days';
```

```sql+sqlite
select
  agent_name,
  agent_id,
  agent_status,
  created_at,
  foundation_model,
  description
from
  aws_bedrock_agent
where
  created_at >= datetime('now', '-30 days');
``` 