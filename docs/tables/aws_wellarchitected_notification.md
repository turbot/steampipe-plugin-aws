---
title: "Table: aws_wellarchitected_notification - Query AWS Well-Architected Tool Notifications using SQL"
description: "Allows users to query AWS Well-Architected Tool Notifications for detailed information about each notification."
---

# Table: aws_wellarchitected_notification - Query AWS Well-Architected Tool Notifications using SQL

The `aws_wellarchitected_notification` table in Steampipe provides information about notifications within AWS Well-Architected Tool. This table allows DevOps engineers to query notification-specific details, including the associated workload, notification type, and status. Users can utilize this table to gather insights on notifications, such as notifications associated with a specific workload, the status of notifications, and more. The schema outlines the various attributes of the notification, including the workload ID, notification type, status, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_notification` table, you can use the `.inspect aws_wellarchitected_notification` command in Steampipe.

**Key columns**:

- `workload_id`: This is the identifier for the workload associated with the notification. It can be used to join this table with the `aws_wellarchitected_workload` table for more detailed insights on the workload.
- `notification_type`: This column indicates the type of notification. It can be used to filter notifications based on their type.
- `status`: The status of the notification. It can be useful to filter notifications based on their status.

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