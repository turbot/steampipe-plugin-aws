# Table: aws_iot_thing_group

An IoT Thing Group in AWS IoT Core is a logical grouping of IoT devices (Things) that allows you to manage and interact with these devices collectively. Thing Groups are part of the device management capabilities in AWS IoT Core, and they provide a way to organize your fleet of devices for easier management and control. This concept is especially useful in large-scale IoT deployments where handling devices individually is impractical.

## Examples

### Basic info

```sql
select
  group_name,
  thing_group_id,
  thing_group_description,
  arn,
  creation_date,
  parent_group_name
from
  aws_iot_thing_group;
```

### Filter thing groups by parent group

```sql
select
  group_name,
  thing_group_id,
  creation_date,
  parent_group_name,
  version
from
  aws_iot_thing_group
where
  parent_group_name = 'foo';
```

### List thing groups created in the last 30 days

```sql
select
  group_name,
  thing_group_id,
  parent_group_name,
  creation_date,
  status
from
  aws_iot_thing_group
where
  creation_date >= now() - interval '30' day;
```

### List dynamic groups

```sql
select
  group_name,
  thing_group_id,
  query_string,
  query_version,
  status
from
  aws_iot_thing_group
where
  status <> '';
```
