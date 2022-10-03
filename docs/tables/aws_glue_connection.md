# Table: aws_glue_connection

An AWS Glue connection is a Data Catalog object that stores connection information for a particular data store. Connections store login credentials, URI strings, virtual private cloud (VPC) information, and more. Creating connections in the Data Catalog saves the effort of having to specify all connection details every time you create a crawler or job. You can use connections for both sources and targets.

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