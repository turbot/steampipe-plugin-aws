# Table: aws_redshift_parameter_group

A parameter group contains a WLM configuration and a set of cluster parameters. It optimizes query performance.
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


### Details of several parameters which are associated with each parameter group

```sql
select
  name,
  jsonb_pretty(parameters)
from
  aws_redshift_parameter_group;
```