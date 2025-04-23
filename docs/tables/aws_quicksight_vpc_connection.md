---
title: "Steampipe Table: aws_quicksight_vpc_connection - Query AWS QuickSight VPC Connections using SQL"
description: "Allows users to query AWS QuickSight VPC Connections to retrieve details about VPC connections used for secure data access."
folder: "QuickSight"
---

# Table: aws_quicksight_vpc_connection - Query AWS QuickSight VPC Connections using SQL

AWS QuickSight VPC Connection is a feature that allows you to connect QuickSight to your data sources in a VPC. This connection enables secure access to data sources within your VPC for visualization and analysis in QuickSight. VPC connections help maintain data security by providing a private network connection between QuickSight and your data sources.

## Table Usage Guide

The `aws_quicksight_vpc_connection` table in Steampipe provides you with information about VPC connections within AWS QuickSight. This table allows you, as a data analyst or administrator, to query VPC connection-specific details, including configuration, status, and security settings. You can use this table to gather insights on VPC connections, such as their current status, the VPC and security groups they're connected to, security group configurations, and more. The schema outlines the various attributes of the QuickSight VPC connection for you to query, including the connection ID, ARN, creation time, and associated network details.

## Examples

### Basic info
Retrieve basic information about QuickSight VPC connections.

```sql+postgres
select
  name,
  vpc_connection_id,
  status,
  availability_status,
  vpc_id
from
  aws_quicksight_vpc_connection;
```

```sql+sqlite
select
  name,
  vpc_connection_id,
  status,
  availability_status,
  vpc_id
from
  aws_quicksight_vpc_connection;
```

### List recently created VPC connections
Identify VPC connections that were created within the last 30 days.

```sql+postgres
select
  name,
  vpc_connection_id,
  created_time,
  status
from
  aws_quicksight_vpc_connection
where
  created_time > now() - interval '30 days';
```

```sql+sqlite
select
  name,
  vpc_connection_id,
  created_time,
  status
from
  aws_quicksight_vpc_connection
where
  created_time > datetime('now', '-30 days');
```

### Get VPC connections with their associated security groups
Retrieve VPC connections along with their associated security groups for security assessment.

```sql+postgres
select
  name,
  vpc_connection_id,
  vpc_id,
  security_group_ids
from
  aws_quicksight_vpc_connection;
```

```sql+sqlite
select
  name,
  vpc_connection_id,
  vpc_id,
  security_group_ids
from
  aws_quicksight_vpc_connection;
```

### List VPC connections by availability status
Group VPC connections by their availability status to understand the current operational state.

```sql+postgres
select
  availability_status,
  count(*) as connection_count
from
  aws_quicksight_vpc_connection
group by
  availability_status;
```

```sql+sqlite
select
  availability_status,
  count(*) as connection_count
from
  aws_quicksight_vpc_connection
group by
  availability_status;
```

### Find VPC connections with DNS resolvers configured
Identify VPC connections that have DNS resolvers configured for name resolution.

```sql+postgres
select
  name,
  vpc_connection_id,
  dns_resolvers
from
  aws_quicksight_vpc_connection
where
  dns_resolvers is not null
  and jsonb_array_length(dns_resolvers) > 0;
```

```sql+sqlite
select
  name,
  vpc_connection_id,
  dns_resolvers
from
  aws_quicksight_vpc_connection
where
  dns_resolvers is not null;
```
