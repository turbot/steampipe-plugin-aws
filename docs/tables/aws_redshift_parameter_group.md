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


### Details of the parameters associated with each parameter group

```sql
select
  name,
  p ->> 'Description' as description,
  p ->> 'ParameterName' as parameter_name,
  p ->> 'ParameterValue' as parameter_value,
  p ->> 'Source' as source,
  p ->> 'DataType' as data_type,
  p ->> 'ApplyType' as apply_type
from
  aws_redshift_parameter_group,
  jsonb_array_elements(parameters) as p;
```