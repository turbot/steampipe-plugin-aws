---
title: "Steampipe Table: aws_glue_connection - Query AWS Glue Connections using SQL"
description: "Allows users to query AWS Glue Connections to retrieve detailed information about each connection."
folder: "Glue"
---

# Table: aws_glue_connection - Query AWS Glue Connections using SQL

The AWS Glue Connection is a component of AWS Glue which allows you to store and retrieve metadata related to your data sources, data targets, and transformations. It facilitates the management of data across multiple data stores by providing a unified view of your data. This enables AWS Glue to connect to your source and target databases, data warehouses, and data lakes for data extraction, transformation, and loading (ETL) processes.

## Table Usage Guide

The `aws_glue_connection` table in Steampipe provides you with information about connections within AWS Glue. This table allows you, as a DevOps engineer, to query connection-specific details, including the connection name, the connection type, the physical connection requirements, and the connection properties. You can utilize this table to gather insights on connections, such as the type of connections, their properties, and the requirements for physical connections. The schema outlines the various attributes of the AWS Glue connection for you, including the catalog ID, creation time, last updated time, match criteria, and associated tags.

## Examples

### Basic info
Explore which AWS Glue connections are currently established to understand their type, creation time, and the region they're in. This can help in managing and optimizing the use of AWS resources.

```sql+postgres
select
  name,
  connection_type,
  creation_time,
  description,
  region
from
  aws_glue_connection;
```

```sql+sqlite
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
This query helps you examine the properties of JDBC connections, including connection URLs and SSL status. It's useful for managing and auditing your database connections, ensuring they are secure and set up correctly.

```sql+postgres
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

```sql+sqlite
select
  name,
  connection_type,
  json_extract(connection_properties, '$.JDBC_CONNECTION_URL') as connection_url,
  json_extract(connection_properties, '$.JDBC_ENFORCE_SSL') as ssl_enabled,
  creation_time
from
  aws_glue_connection
where
  connection_type = 'JDBC';
```

### List mongodb connections with ssl disabled
Identify instances where MongoDB connections have SSL disabled to assess potential security vulnerabilities. This can be useful in maintaining secure data practices by pinpointing the specific connections that may require updating.

```sql+postgres
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

```sql+sqlite
select
  name,
  connection_type,
  json_extract(connection_properties, '$.CONNECTION_URL') as connection_url,
  json_extract(connection_properties, '$.JDBC_ENFORCE_SSL') as ssl_enabled,
  creation_time
from
  aws_glue_connection
where
  connection_type = 'JDBC'
  and json_extract(connection_properties, '$.JDBC_ENFORCE_SSL') = 'false';
```

### List connection vpc details
This query is useful to analyze the details of your AWS Glue connections in relation to their corresponding VPC subnets. It helps in assessing the configuration of physical connection requirements and understanding the link between different AWS resources.

```sql+postgres
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

```sql+sqlite
select
  c.name as connection_name,
  s.vpc_id as vpc_id,
  s.title as subnet_name,
  json_extract(physical_connection_requirements, '$.SubnetId') as subnet_id,
  json_extract(physical_connection_requirements, '$.AvailabilityZone') as availability_zone,
  cidr_block,
  json_extract(physical_connection_requirements, '$.SecurityGroupIdList') as security_group_ids
from
  aws_glue_connection c
  join aws_vpc_subnet s on json_extract(c.physical_connection_requirements, '$.SubnetId') = s.subnet_id;
```