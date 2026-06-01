---
title: "Steampipe Table: aws_ecr_registry - Query AWS ECR Private Registry using SQL"
description: "Allows users to query the AWS ECR private registry, including its replication configuration and registry permissions policy."
folder: "ECR"
---

# Table: aws_ecr_registry - Query AWS ECR Private Registry using SQL

Each AWS account has one Amazon Elastic Container Registry (ECR) private registry per Region. The registry exposes account-wide settings such as the replication configuration (which destinations images are copied to) and the registry permissions policy (which principals can replicate images into the registry).

## Table Usage Guide

The `aws_ecr_registry` table in Steampipe provides you with information about the per-region ECR private registry for your account. It returns a single row per region and exposes the registry's replication configuration and registry policy. Use it to audit cross-region/cross-account replication setup, confirm that a registry policy is (or is not) attached, and understand which accounts can replicate images into your registry.

## Examples

### Basic info
List every regional registry along with the regions it replicates to.

```sql+postgres
select
  registry_id,
  region,
  jsonb_pretty(replication_configuration) as replication_configuration
from
  aws_ecr_registry;
```

```sql+sqlite
select
  registry_id,
  region,
  replication_configuration
from
  aws_ecr_registry;
```

### Find registries with no replication configured
Identify regions where image replication is not configured. This can highlight gaps in disaster-recovery or multi-region image distribution strategies.

```sql+postgres
select
  registry_id,
  region
from
  aws_ecr_registry
where
  replication_configuration is null
  or jsonb_array_length(replication_configuration -> 'Rules') = 0;
```

```sql+sqlite
select
  registry_id,
  region
from
  aws_ecr_registry
where
  replication_configuration is null
  or json_array_length(json_extract(replication_configuration, '$.Rules')) = 0;
```

### List registry replication destinations
Flatten the replication configuration to show each destination region and account that the registry replicates images to.

```sql+postgres
select
  r.registry_id,
  r.region as source_region,
  dest ->> 'Region'     as destination_region,
  dest ->> 'RegistryId' as destination_registry_id
from
  aws_ecr_registry as r,
  jsonb_array_elements(r.replication_configuration -> 'Rules')        as rule,
  jsonb_array_elements(rule -> 'Destinations')                        as dest;
```

```sql+sqlite
select
  r.registry_id,
  r.region as source_region,
  json_extract(dest.value, '$.Region')     as destination_region,
  json_extract(dest.value, '$.RegistryId') as destination_registry_id
from
  aws_ecr_registry as r,
  json_each(json_extract(r.replication_configuration, '$.Rules')) as rule,
  json_each(json_extract(rule.value, '$.Destinations')) as dest;
```

### Find registries that have no registry policy attached
Registries without a registry policy cannot accept cross-account image replication.

```sql+postgres
select
  registry_id,
  region
from
  aws_ecr_registry
where
  policy is null;
```

```sql+sqlite
select
  registry_id,
  region
from
  aws_ecr_registry
where
  policy is null;
```

### List principals granted access by the registry policy
Inspect which principals are granted access by the registry permissions policy.

```sql+postgres
select
  registry_id,
  region,
  stmt ->> 'Sid'      as sid,
  stmt ->> 'Effect'   as effect,
  stmt ->  'Principal' as principal,
  stmt ->  'Action'    as action
from
  aws_ecr_registry as r,
  jsonb_array_elements(r.policy -> 'Statement') as stmt
where
  r.policy is not null;
```

```sql+sqlite
select
  registry_id,
  region,
  json_extract(stmt.value, '$.Sid')       as sid,
  json_extract(stmt.value, '$.Effect')    as effect,
  json_extract(stmt.value, '$.Principal') as principal,
  json_extract(stmt.value, '$.Action')    as action
from
  aws_ecr_registry as r,
  json_each(json_extract(r.policy, '$.Statement')) as stmt
where
  r.policy is not null;
```
