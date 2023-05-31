# Table: aws_docdb_cluster_instance

An Amazon DocumentDB instance is an isolated database environment in the cloud. An instance can contain multiple user-created databases.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  engine,
  engine_version,
  db_instance_class,
  availability_zone
from
  aws_docdb_cluster_instance;
```

### List instances which are publicly accessible

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  engine,
  engine_version,
  db_instance_class,
  availability_zone
from
  aws_docdb_cluster_instance
where
  publicly_accessible;
```

### Get DB subnet group info of each instance

```sql
select
  db_subnet_group_arn,
  db_subnet_group_name,
  db_subnet_group_description,
  db_subnet_group_status
from
  aws_docdb_cluster_instance;
```

### Get VPC and subnet info of each instance

```sql
select
  db_instance_identifier as attached_vpc,
  vsg ->> 'VpcSecurityGroupId' as vpc_security_group_id,
  vsg ->> 'Status' as status,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub -> 'SubnetOutpost' ->> 'Arn' as subnet_outpost,
  sub ->> 'SubnetStatus' as subnet_status
from
  aws_docdb_cluster_instance
  cross join jsonb_array_elements(vpc_security_groups) as vsg
  cross join jsonb_array_elements(subnets) as sub;
```

### List instances with unecrypted storage

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  not storage_encrypted;
```

### List instances with cloudwatch logs disabled

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  enabled_cloudwatch_logs_exports is null;
```

### Get endpoint info for each instance

```sql
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_docdb_cluster_instance;
```
