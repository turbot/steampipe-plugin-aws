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
  jsonb_pretty(parameters)
from
  aws_redshift_parameter_group;
```