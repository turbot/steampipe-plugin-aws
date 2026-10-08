---
title: "Steampipe Table: aws_bedrock_agentcore_agent_runtime - Query AWS Bedrock AgentCore Agent Runtimes using SQL"
description: "Allows users to query AWS Bedrock AgentCore Runtimes, providing information about hosted AI agents including their network configuration, container or code artifacts, IAM roles, and status."
folder: "Bedrock"
---

# Table: aws_bedrock_agentcore_agent_runtime - Query AWS Bedrock AgentCore Agent Runtimes using SQL

Amazon Bedrock AgentCore Runtime is a serverless hosting environment for AI agents and tools. Each agent runtime packages agent code as a container image or code bundle, runs it in an isolated session, and exposes it over HTTP, MCP, or A2A protocols. Runtimes can be configured with public internet access or attached to your VPC subnets and security groups.

## Table Usage Guide

The `aws_bedrock_agentcore_agent_runtime` table in Steampipe provides you with information about agent runtimes within Amazon Bedrock AgentCore. This table allows you, as a security engineer or platform administrator, to query runtime-specific details, including network mode, VPC subnets and security groups, execution role, artifact configuration, and status. You can utilize this table to verify that runtimes are attached to a VPC, review which IAM roles they assume, and audit their protocol and authorizer settings. The schema outlines the various attributes of the agent runtime for you, including the runtime ID, name, version, ARN, and associated tags.

**Important Notes**
- The `network_configuration`, `role_arn`, `agent_runtime_artifact`, and other configuration columns are populated by a per-row `GetAgentRuntime` call in addition to the `ListAgentRuntimes` call. Querying by `agent_runtime_id` fetches a single runtime directly.
- `network_configuration` holds the `NetworkMode` (`PUBLIC` or `VPC`) and, for VPC runtimes, the `Subnets` and `SecurityGroups` under `NetworkModeConfig`.
- `agent_runtime_artifact` is keyed by artifact type: either `ContainerConfiguration` (container image URI) or `CodeConfiguration` (S3 code location, runtime, and entry point). `authorizer_configuration` is keyed by authorizer type, for example `CustomJWTAuthorizer`, and is null when the runtime uses AWS IAM authorization. Each entry in `filesystem_configurations` is keyed by mount type, for example `EfsAccessPoint`.

## Examples

### Basic info
Explore the basic details of your agent runtimes to understand their status and versions. This can help in managing and monitoring your hosted agents effectively.

```sql+postgres
select
  agent_runtime_name,
  agent_runtime_id,
  agent_runtime_version,
  status,
  last_updated_at,
  region
from
  aws_bedrock_agentcore_agent_runtime;
```

```sql+sqlite
select
  agent_runtime_name,
  agent_runtime_id,
  agent_runtime_version,
  status,
  last_updated_at,
  region
from
  aws_bedrock_agentcore_agent_runtime;
```

### List runtimes that are not attached to a VPC
Identify agent runtimes running in public network mode. These runtimes have direct internet access through the AWS network rather than routing traffic through your VPC, which may violate your organization's network security policy.

```sql+postgres
select
  agent_runtime_name,
  agent_runtime_id,
  network_configuration ->> 'NetworkMode' as network_mode,
  account_id,
  region
from
  aws_bedrock_agentcore_agent_runtime
where
  network_configuration ->> 'NetworkMode' <> 'VPC';
```

```sql+sqlite
select
  agent_runtime_name,
  agent_runtime_id,
  network_configuration ->> 'NetworkMode' as network_mode,
  account_id,
  region
from
  aws_bedrock_agentcore_agent_runtime
where
  network_configuration ->> 'NetworkMode' <> 'VPC';
```

### List VPC-attached runtimes with their subnets and security groups
Review the VPC placement of your agent runtimes to confirm they are using the expected subnets and security groups.

```sql+postgres
select
  agent_runtime_name,
  jsonb_array_elements_text(network_configuration -> 'NetworkModeConfig' -> 'Subnets') as subnet_id,
  network_configuration -> 'NetworkModeConfig' -> 'SecurityGroups' as security_groups
from
  aws_bedrock_agentcore_agent_runtime
where
  network_configuration ->> 'NetworkMode' = 'VPC';
```

```sql+sqlite
select
  agent_runtime_name,
  json_each.value as subnet_id,
  json_extract(network_configuration, '$.NetworkModeConfig.SecurityGroups') as security_groups
from
  aws_bedrock_agentcore_agent_runtime,
  json_each(json_extract(network_configuration, '$.NetworkModeConfig.Subnets'))
where
  json_extract(network_configuration, '$.NetworkMode') = 'VPC';
```

### List runtimes with the VPC each subnet belongs to
Join with the subnet table to see which VPC each runtime is placed in. This is useful for confirming runtimes are in approved VPCs rather than only checking that a VPC is configured.

```sql+postgres
select
  r.agent_runtime_name,
  r.region,
  s.subnet_id,
  s.vpc_id,
  s.availability_zone
from
  aws_bedrock_agentcore_agent_runtime as r,
  jsonb_array_elements_text(r.network_configuration -> 'NetworkModeConfig' -> 'Subnets') as subnet
  join aws_vpc_subnet as s on s.subnet_id = subnet
where
  r.network_configuration ->> 'NetworkMode' = 'VPC';
```

```sql+sqlite
select
  r.agent_runtime_name,
  r.region,
  s.subnet_id,
  s.vpc_id,
  s.availability_zone
from
  aws_bedrock_agentcore_agent_runtime as r,
  json_each(json_extract(r.network_configuration, '$.NetworkModeConfig.Subnets')) as subnet
  join aws_vpc_subnet as s on s.subnet_id = subnet.value
where
  json_extract(r.network_configuration, '$.NetworkMode') = 'VPC';
```

### List runtimes with their execution roles
Review which IAM role each runtime assumes. This helps you audit the permissions granted to your hosted agents.

```sql+postgres
select
  agent_runtime_name,
  role_arn,
  status,
  region
from
  aws_bedrock_agentcore_agent_runtime;
```

```sql+sqlite
select
  agent_runtime_name,
  role_arn,
  status,
  region
from
  aws_bedrock_agentcore_agent_runtime;
```

### List runtimes that are not in READY status
Find runtimes that are still being created, updated, or that have failed, along with the failure reason where one is reported.

```sql+postgres
select
  agent_runtime_name,
  status,
  failure_reason,
  last_updated_at
from
  aws_bedrock_agentcore_agent_runtime
where
  status <> 'READY';
```

```sql+sqlite
select
  agent_runtime_name,
  status,
  failure_reason,
  last_updated_at
from
  aws_bedrock_agentcore_agent_runtime
where
  status <> 'READY';
```

### Get the protocol and authorizer configuration of a runtime
Inspect how a specific runtime is exposed and how inbound requests are authorized.

```sql+postgres
select
  agent_runtime_name,
  protocol_configuration ->> 'ServerProtocol' as server_protocol,
  authorizer_configuration,
  agent_runtime_artifact
from
  aws_bedrock_agentcore_agent_runtime
where
  agent_runtime_id = 'my_agent-abc123XYZ0';
```

```sql+sqlite
select
  agent_runtime_name,
  json_extract(protocol_configuration, '$.ServerProtocol') as server_protocol,
  authorizer_configuration,
  agent_runtime_artifact
from
  aws_bedrock_agentcore_agent_runtime
where
  agent_runtime_id = 'my_agent-abc123XYZ0';
```

### List runtimes by artifact type
Separate runtimes deployed from a container image from those deployed as a code bundle, and show where each artifact lives.

```sql+postgres
select
  agent_runtime_name,
  region,
  case
    when agent_runtime_artifact ? 'ContainerConfiguration' then 'container'
    else 'code'
  end as artifact_type,
  agent_runtime_artifact -> 'ContainerConfiguration' ->> 'ContainerUri' as container_uri,
  agent_runtime_artifact -> 'CodeConfiguration' ->> 'Runtime' as code_runtime,
  agent_runtime_artifact -> 'CodeConfiguration' -> 'Code' -> 'S3' ->> 'Bucket' as code_bucket
from
  aws_bedrock_agentcore_agent_runtime;
```

```sql+sqlite
select
  agent_runtime_name,
  region,
  case
    when json_extract(agent_runtime_artifact, '$.ContainerConfiguration') is not null then 'container'
    else 'code'
  end as artifact_type,
  json_extract(agent_runtime_artifact, '$.ContainerConfiguration.ContainerUri') as container_uri,
  json_extract(agent_runtime_artifact, '$.CodeConfiguration.Runtime') as code_runtime,
  json_extract(agent_runtime_artifact, '$.CodeConfiguration.Code.S3.Bucket') as code_bucket
from
  aws_bedrock_agentcore_agent_runtime;
```

### List runtimes that accept JWT bearer tokens
Find runtimes using a custom JWT authorizer instead of AWS IAM, along with the issuer and allowed clients.

```sql+postgres
select
  agent_runtime_name,
  region,
  authorizer_configuration -> 'CustomJWTAuthorizer' ->> 'DiscoveryUrl' as discovery_url,
  authorizer_configuration -> 'CustomJWTAuthorizer' -> 'AllowedClients' as allowed_clients,
  authorizer_configuration -> 'CustomJWTAuthorizer' -> 'AllowedAudience' as allowed_audience
from
  aws_bedrock_agentcore_agent_runtime
where
  authorizer_configuration ? 'CustomJWTAuthorizer';
```

```sql+sqlite
select
  agent_runtime_name,
  region,
  json_extract(authorizer_configuration, '$.CustomJWTAuthorizer.DiscoveryUrl') as discovery_url,
  json_extract(authorizer_configuration, '$.CustomJWTAuthorizer.AllowedClients') as allowed_clients,
  json_extract(authorizer_configuration, '$.CustomJWTAuthorizer.AllowedAudience') as allowed_audience
from
  aws_bedrock_agentcore_agent_runtime
where
  json_extract(authorizer_configuration, '$.CustomJWTAuthorizer') is not null;
```
