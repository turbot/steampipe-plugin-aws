---
title: "Steampipe Table: aws_codeconnections_host - Query AWS CodeConnections hosts using SQL"
description: "Allows users to query AWS CodeConnections hosts, including their provider endpoints, status, VPC configuration, and tags."
folder: "CodeConnections"
---

# Table: aws_codeconnections_host - Query AWS CodeConnections hosts using SQL

AWS CodeConnections hosts represent infrastructure where a supported third-party provider, such as GitHub Enterprise Server, is installed. One host can be used by multiple connections to that provider.

## Table Usage Guide

The `aws_codeconnections_host` table provides host details for each configured AWS region, including provider endpoints, setup status, VPC configuration, and tags.

## Examples

### Basic info

```sql+postgres
select
  name,
  provider_type,
  provider_endpoint,
  status,
  region
from
  aws_codeconnections_host;
```

```sql+sqlite
select
  name,
  provider_type,
  provider_endpoint,
  status,
  region
from
  aws_codeconnections_host;
```

### Find hosts that are not available

```sql+postgres
select
  name,
  arn,
  status,
  status_message
from
  aws_codeconnections_host
where
  status <> 'AVAILABLE';
```

```sql+sqlite
select
  name,
  arn,
  status,
  status_message
from
  aws_codeconnections_host
where
  status <> 'AVAILABLE';
```

### Review VPC settings for hosts

```sql+postgres
select
  name,
  vpc_configuration ->> 'VpcId' as vpc_id,
  vpc_configuration -> 'SubnetIds' as subnet_ids,
  vpc_configuration -> 'SecurityGroupIds' as security_group_ids
from
  aws_codeconnections_host
where
  vpc_configuration is not null;
```

```sql+sqlite
select
  name,
  json_extract(vpc_configuration, '$.VpcId') as vpc_id,
  json_extract(vpc_configuration, '$.SubnetIds') as subnet_ids,
  json_extract(vpc_configuration, '$.SecurityGroupIds') as security_group_ids
from
  aws_codeconnections_host
where
  vpc_configuration is not null;
```
