# Table: aws_wellarchitected_notification

A Notification indicates that a new version of a Well-Architected lens is available.

## Examples

### List notifications for workloads where lens version is upgraded

```sql
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  notification_type = 'LENS_VERSION_UPGRADED';
```

### List notifications for workloads where lens version is deprecated

```sql
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  notification_type = 'LENS_VERSION_DEPRECATED';
```

### Check if there is a notification for a particular workload

```sql
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  workload_id = '123451c59cebcd4612f1f858bf75566';
```