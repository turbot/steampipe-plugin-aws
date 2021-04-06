# Table: aws_redshift_parameter_group

A parameter group is a group of parameters that apply to all of the databases that you create in the cluster. 

## Examples

### Basic info

```sql
select
  name,
  description,
  family
from
  aws_redshift_parameter_group;
```


### Get the details of the require_ssl parameter associated with each parameter group

```sql
select
  name,
  p ->> 'ParameterName' as parameter_name,
  p ->> 'ParameterValue' as parameter_value,
  p ->> 'Description' as description,
  p ->> 'Source' as source,
  p ->> 'DataType' as data_type,
  p ->> 'ApplyType' as apply_type,
  p ->> 'IsModifiable' as is_modifiable,
  p ->> 'AllowedValues' as allowed_values,
  p ->> 'MinimumEngineVersion' as minimum_engine_version
from
  aws_redshift_parameter_group,
  jsonb_array_elements(parameters) as p
where
  (
    p ->> 'ParameterName' = 'require_ssl'
    and p ->> 'ParameterValue' = 'false'
  )
  or (
    p ->> 'ParameterName' = 'require_ssl'
    and p ->> 'ParameterValue' = 'true'
  );
```