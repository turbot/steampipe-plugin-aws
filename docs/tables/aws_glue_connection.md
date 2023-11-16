---
title: "Table: aws_glue_connection - Query AWS Glue Connections using SQL"
description: "Allows users to query AWS Glue Connections to retrieve detailed information about each connection."
---

# Table: aws_glue_connection - Query AWS Glue Connections using SQL

The `aws_glue_connection` table in Steampipe provides information about connections within AWS Glue. This table allows DevOps engineers to query connection-specific details, including the connection name, the connection type, the physical connection requirements, and the connection properties. Users can utilize this table to gather insights on connections, such as the type of connections, their properties, and the requirements for physical connections. The schema outlines the various attributes of the AWS Glue connection, including the catalog ID, creation time, last updated time, match criteria, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_connection` table, you can use the `.inspect aws_glue_connection` command in Steampipe.

### Key columns:

- `name`: The name of the connection. This can be used to join this table with other tables that require connection name.
- `connection_type`: The type of the connection (JDBC, SFTP, etc.). This can provide insights into the type of connections used.
- `physical_connection_requirements`: The physical connection requirements. This can be useful in understanding the requirements for establishing the connection.

## Examples

### Basic info

```sql
select
  name,
  connection_type,
  creation_time,
  description,
  region
from
  aws_glue_connection;
```

### List connection properties for JDBC connections

```sql
select
  name,
  connection_type,
  connection_properties ->> 'JDBC_CONNECTION_URL' as connection_url,
  connection_properties ->> 'JDBC_ENFORCE_SSL' as ssl_enabled,
  creation_time
from
  aws_glue_connection
where
  connection_type = 'JDBC';
```

### List mongodb connections with ssl disabled

```sql
select
  name,
  connection_type,
  connection_properties ->> 'CONNECTION_URL' as connection_url,
  connection_properties ->> 'JDBC_ENFORCE_SSL' as ssl_enabled,
  creation_time
from
  aws_glue_connection
where
  connection_type = 'JDBC'
  and connection_properties ->> 'JDBC_ENFORCE_SSL' = 'false';
```

### List connection vpc details

```sql
select
  c.name as connection_name,
  s.vpc_id as vpc_id,
  s.title as subnet_name,
  physical_connection_requirements ->> 'SubnetId' as subnet_id,
  physical_connection_requirements ->> 'AvailabilityZone' as availability_zone,
  cidr_block,
  physical_connection_requirements ->> 'SecurityGroupIdList' as security_group_ids
from
  aws_glue_connection c
  join aws_vpc_subnet s on physical_connection_requirements ->> 'SubnetId' = s.subnet_id;
```