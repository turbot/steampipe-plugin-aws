# Table: aws_iot_thing

AWS IoT Things refer to the virtual representations of physical devices or assets in AWS IoT Core, which is a part of Amazon Web Services (AWS) designed for the Internet of Things (IoT). These "Things" can be any connected device, such as sensors, appliances, machinery, or any other objects that can collect and exchange data over the internet or other networks.

## Examples

### Basic info

```sql
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing;
```

### Filter things by attribute name

```sql
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing
where
  attribute_name = 'foo';
```

### List things for a given type name

```sql
select
  thing_name,
  arn,
  thing_id,
  thing_type_name,
  attribute_value
from
  aws_iot_thing
where
  thing_type_name = 'foo';
```
