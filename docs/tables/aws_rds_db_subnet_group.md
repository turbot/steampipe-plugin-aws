# Table: aws_rds_db_subnet_group

A DB subnet group is a collection of subnets (typically private) that are created for a VPC and then designated for DB instances.

## Examples

### DB subnet group basic info

```sql
select
  name,
  status,
  vpc_id
from
  aws_rds_db_subnet_group;
```


### Subnets info of each subnet in subnet group

```sql
select
  name,
  subnet -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  subnet ->> 'SubnetIdentifier' as subnet_identifier,
  subnet -> 'SubnetOutpost' ->> 'Arn' as subnet_outpost,
  subnet ->> 'SubnetStatus' as subnet_status
from
  aws_rds_db_subnet_group
  cross join jsonb_array_elements(subnets) as subnet;
```


### List of subnet group without application tag key

```sql
select
  name,
  tags
from
  aws_rds_db_subnet_group
where
  not tags :: JSONB ? 'application';
```
