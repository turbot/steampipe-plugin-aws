---
title: "Table: aws_service_discovery_namespace - Query AWS Cloud Map Service Discovery Namespace using SQL"
description: "Allows users to query AWS Cloud Map Service Discovery Namespace to retrieve details about the namespaces in AWS Cloud Map."
---

# Table: aws_service_discovery_namespace - Query AWS Cloud Map Service Discovery Namespace using SQL

The `aws_service_discovery_namespace` table in Steampipe provides information about AWS Cloud Map Service Discovery Namespaces. This table allows DevOps engineers to query namespace-specific details, including namespace type (DNS, HTTP), associated services, and associated metadata. Users can utilize this table to gather insights on namespaces, such as the number of services in each namespace, namespace types, and more. The schema outlines the various attributes of the service discovery namespace, including the namespace ID, ARN, name, type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_service_discovery_namespace` table, you can use the `.inspect aws_service_discovery_namespace` command in Steampipe.

### Key columns:

- `id`: The ID of the namespace. This can be used to join this table with other tables.
- `arn`: The Amazon Resource Name (ARN) of the namespace. This column is useful for joining with other tables that use ARN as a reference.
- `name`: The name of the namespace. This column is important for understanding the purpose of the namespace and can be used to join with other tables that reference the namespace by name.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_namespace;
```

### List private namespaces

```sql
select
  name,
  id,
  arn,
  type,
  service_count
from
  aws_service_discovery_namespace
where
  type ilike '%private%';
```

### List HTTP type namespaces

```sql
select
  name,
  id,
  arn,
  type,
  service_account
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### List namespaces created in the last 30 days

```sql
select
  name,
  id,
  description,
  create_date
from
  aws_service_discovery_namespace
where
  create_date >= now() - interval '30' day;
```

### Get HTTP property details of namespaces

```sql
select
  name,
  id,
  http_properties ->> 'HttpName' as http_name
from
  aws_service_discovery_namespace
where
  type = 'HTTP';
```

### Get private DNS property details of namspaces

```sql
select
  name,
  id,
  dns_properties ->> 'HostedZoneId' as HostedZoneId,
  dns_properties -> 'SOA' ->> 'TTL' as ttl
from
  aws_service_discovery_namespace
where
  type = 'DNS_PRIVATE';
```

### Count namespaces by type

```sql
select
  type,
  count(type)
from
  aws_service_discovery_namespace
group by
  type;
```
