# Table: aws_rds_db_instance

A DB instance is an isolated database environment running in the cloud.

## Examples

### List of DB instances which are publicly accessible

```sql
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  publicly_accessible
from
  aws_rds_db_instance
where
  publicly_accessible;
```


### List of DB instances which are not authenticated through IAM users and roles.

```sql
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  publicly_accessible
from
  aws_rds_db_instance
where
  publicly_accessible;
```


### VPC and subnet info associated with DB instance

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
  aws_rds_db_instance
  cross join jsonb_array_elements(vpc_security_groups) as vsg
  cross join jsonb_array_elements(subnets) as sub;
```


### List of RDS instances where deletion protection is disabled

```sql
select
  db_instance_identifier,
  class,
  engine,
  engine_version,
  deletion_protection
from
  aws_rds_db_instance
where
  not deletion_protection;
```


### List of RDS instances whose storage is not encrypted

```sql
select
  db_instance_identifier,
  class,
  allocated_storage,
  deletion_protection
from
  aws_rds_db_instance
where
  not storage_encrypted;
```


### Endpoint info of each db instances

```sql
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_rds_db_instance;
```