# Table: aws_dax_parameter_group

Amazon DynamoDB Accelerator (DAX) Parameter Group is a named set of parameters that are applied to every node in the cluster. The parameters are the type of configuration that will be applied to the cluster during the creation.

## Examples

### Basic info

```sql
select
  parameter_name,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type
from
  aws_dax_parameter;
```

### Count parameters by parameter group

```sql
select
  parameter_group_name,
  count(parameter_name) number_of_parameters
from
  aws_dax_parameter
group by
  parameter_group_name;
```

### List modifiable parameters

```sql
select
  parameter_name,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type,
  is_modifiable
from
  aws_dax_parameter
where
  is_modifiable = 'TRUE';
```

### List parameters which are not user define

```sql
select
  parameter_name,
  change_type,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type,
  source
from
  aws_dax_parameter
where
  source <> 'user';
  ```