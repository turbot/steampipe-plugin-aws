# Table: aws_health_affected_entity

Information about an entity that is affected by a Health event.

## Examples

### Basic info

```sql
select
  arn,
  entity_url,
  entity_value,
  event_arn,
  last_updated_time,
  status_code
from
  aws_health_affected_entity;
```

### List affected entities that are unimpaired

```sql
select
  arn,
  entity_url,
  entity_value,
  event_arn,
  last_updated_time,
  status_code
from
  aws_health_affected_entity
where
  status_code = 'UNIMPAIRED';
```

### Get health event details for each entity

```sql
select
  e.arn,
  e.entity_url,
  e.event_arn,
  v.event_start_time,
  v.event_end_time,
  v.event_type_category,
  v.event_type_code
  v.service
from
  aws_health_affected_entity as e,
  aws_health_event as v;
```