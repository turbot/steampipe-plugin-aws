# Table: aws_iot_thing_type

AWS IoT Thing Type is a feature in AWS IoT Core that allows you to categorize your IoT devices (Things) into different types based on common characteristics or use cases. A Thing Type is essentially a template that defines attributes and properties common to a particular class or category of devices. By defining Thing Types, you can manage and interact with groups of Things more effectively and consistently.

## Examples

### Basic info

```sql
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date
from
  aws_iot_thing_type;
```

### List deprecated thing types

```sql
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated
from
  aws_iot_thing_type
where
  deprecated;
```

### List thing types created in the last 30 days

```sql
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated,
  searchable_attributes
from
  aws_iot_thing_type
where
  creation_date >= now() - interval '30' day;
```

### List thing types scheduled for deprecation within the next 30 days

```sql
select
  thing_type_name,
  arn,
  thing_type_id,
  creation_date,
  tags,
  deprecation_date
from
  aws_iot_thing_type
where
  deprecation_date <= now() - interval '30' day;
```