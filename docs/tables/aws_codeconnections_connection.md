---
title: "Steampipe Table: aws_codeconnections_connection - Query AWS CodeConnections connections using SQL"
description: "Allows users to query AWS CodeConnections connections, including their providers, status, hosts, and tags."
folder: "CodeConnections"
---

# Table: aws_codeconnections_connection - Query AWS CodeConnections connections using SQL

AWS CodeConnections connections link AWS services such as CodePipeline to external source repositories. The connection remains pending until its provider handshake is completed.

## Table Usage Guide

The `aws_codeconnections_connection` table provides connection details for each configured AWS region, including the provider, current status, associated host, and tags.

## Examples

### Basic info

```sql+postgres
select
  name,
  provider_type,
  connection_status,
  region
from
  aws_codeconnections_connection;
```

```sql+sqlite
select
  name,
  provider_type,
  connection_status,
  region
from
  aws_codeconnections_connection;
```

### Find connections that are not available

```sql+postgres
select
  name,
  arn,
  provider_type,
  connection_status
from
  aws_codeconnections_connection
where
  connection_status <> 'AVAILABLE';
```

```sql+sqlite
select
  name,
  arn,
  provider_type,
  connection_status
from
  aws_codeconnections_connection
where
  connection_status <> 'AVAILABLE';
```

### List connections for a particular provider

```sql+postgres
select
  name,
  arn,
  owner_account_id
from
  aws_codeconnections_connection
where
  provider_type = 'GitHub';
```

```sql+sqlite
select
  name,
  arn,
  owner_account_id
from
  aws_codeconnections_connection
where
  provider_type = 'GitHub';
```
