# Table: aws_rds_db_option_group

An option group can specify features, called options, that are available for a particular Amazon RDS DB instance. Options can have settings that specify how the option works. When a DB instance associated with an option group, the specified options and option settings are enabled for that DB instance.

## Examples

### Basic parameter group info

```sql
select
  name,
  description,
  engine_name,
  major_engine_version,
  vpc_id
from
  aws_rds_db_option_group;
```


### List of option groups which can be applied to both VPC and non-VPC instances

```sql
select
  name,
  description,
  engine_name,
  allows_vpc_and_non_vpc_instance_memberships
from
  aws_rds_db_option_group
where
  allows_vpc_and_non_vpc_instance_memberships;
```


### Option details of each option group

```sql
select
  name,
  option ->> 'OptionName' as option_name,
  option -> 'Permanent' as Permanent,
  option -> 'Persistent' as Persistent,
  option -> 'VpcSecurityGroupMemberships' as vpc_security_group_membership,
  option -> 'Port' as Port
from
  aws_rds_db_option_group
  cross join jsonb_array_elements(options) as option;
```