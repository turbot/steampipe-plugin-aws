---
title: "Steampipe Table: aws_codeconnections_repository_link - Query AWS CodeConnections repository links using SQL"
description: "Allows users to query AWS CodeConnections repository links and their associated connections, repositories, owners, encryption keys, and tags."
folder: "CodeConnections"
---

# Table: aws_codeconnections_repository_link - Query AWS CodeConnections repository links using SQL

AWS CodeConnections repository links associate external Git repositories with connections so that Git sync can monitor and synchronize repository changes.

## Table Usage Guide

The `aws_codeconnections_repository_link` table provides repository link details for each configured AWS region, including the repository, provider, connection, encryption key, and tags.

## Examples

### Basic info

```sql+postgres
select
  repository_name,
  repository_link_id,
  provider_type,
  owner_id,
  region
from
  aws_codeconnections_repository_link;
```

```sql+sqlite
select
  repository_name,
  repository_link_id,
  provider_type,
  owner_id,
  region
from
  aws_codeconnections_repository_link;
```

### Find repository links without a customer managed encryption key

```sql+postgres
select
  repository_name,
  arn,
  connection_arn
from
  aws_codeconnections_repository_link
where
  encryption_key_arn is null;
```

```sql+sqlite
select
  repository_name,
  arn,
  connection_arn
from
  aws_codeconnections_repository_link
where
  encryption_key_arn is null;
```

### Join repository links to their connections

```sql+postgres
select
  r.repository_name,
  r.owner_id,
  c.name as connection_name,
  c.connection_status
from
  aws_codeconnections_repository_link as r
  join aws_codeconnections_connection as c on c.arn = r.connection_arn;
```

```sql+sqlite
select
  r.repository_name,
  r.owner_id,
  c.name as connection_name,
  c.connection_status
from
  aws_codeconnections_repository_link as r
  join aws_codeconnections_connection as c on c.arn = r.connection_arn;
```
