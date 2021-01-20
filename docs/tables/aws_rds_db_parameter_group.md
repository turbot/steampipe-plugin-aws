# Table: aws_rds_db_parameter_group

A DB parameter group acts as a container for engine configuration values that are applied to one or more DB instances.

## Examples

### List of DB parameter group and corresponding parameter group family

```sql
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_parameter_group;
```


### Parameters info of each parameter group.

```sql
select
  name,
  db_parameter_group_family,
  pg ->> 'ParameterName' as parameter_name,
  pg ->> 'ParameterValue' as parameter_value,
  pg ->> 'AllowedValues' as allowed_values,
  pg ->> 'ApplyType' as apply_type,
  pg ->> 'IsModifiable' as is_modifiable,
  pg ->> 'DataType' as data_type,
  pg ->> 'Description' as description,
  pg ->> 'MinimumEngineVersion' as minimum_engine_version
from
  aws_rds_db_parameter_group
  cross join jsonb_array_elements(parameters) as pg;
```