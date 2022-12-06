# Table: aws_cloudtrail_event_data_store

AWS CloudTrail event data stores retain data for up to seven years, or 2557 days. By default, event data is retained for 2557 days, and termination protection is enabled for an event data store.

## Examples

### Basic info

```sql
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store;
```

### List event data stores which are not enabled

```sql
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  status <> 'ENABLED';
```

### List event data stores with termination protection disabled

```sql
select
  name,
  arn,
  status,
  created_timestamp,
  multi_region_enabled,
  organization_enabled,
  termination_protection_enabled
from
  aws_cloudtrail_event_data_store
where
  not termination_protection_enabled;
```
