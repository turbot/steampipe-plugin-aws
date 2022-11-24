# Table: aws_health_event

AWS Health provides ongoing visibility into your resource performance and the availability of your AWS services and accounts. You can use AWS Health events to learn how service and resource changes might affect your applications running on AWS. 

## Examples

### Basic info

```sql
select
  arn,
  availability_zone,
  start_time,
  end_time,
  event_scope_code,
  service,
  region
from
  aws_health_event;
```

### List upcoming events

```sql
select
  arn,
  start_time,
  end_time,
  event_scope_code,
  status_code,
  service
from
  aws_health_event
where
  status_code = 'upcoming';
```

### List event details for the EC2 service

```sql
select
  arn,
  start_time,
  end_time,
  event_scope_code,
  status_code,
  service
from
  aws_health_event
where
  service = 'EC2';
```

### List event details for an availability zone

```sql
select
  arn,
  availability_zone,
  start_time,
  end_time,
  event_scope_code,
  status_code,
  service
from
  aws_health_event
where
  availability_zone = 'us-east-1a';
```
