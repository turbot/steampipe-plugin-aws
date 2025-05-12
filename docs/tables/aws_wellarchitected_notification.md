---
title: "Steampipe Table: aws_wellarchitected_notification - Query AWS Well-Architected Tool Notifications using SQL"
description: "Allows users to query AWS Well-Architected Tool Notifications for detailed information about each notification."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_notification - Query AWS Well-Architected Tool Notifications using SQL

The AWS Well-Architected Tool Notifications is a feature of AWS Well-Architected Tool that allows you to receive important updates and insights about your workloads. It helps you identify potential risks and areas for improvement in your AWS environment. The tool uses best practices from AWS's Well-Architected Framework to provide actionable recommendations.

## Table Usage Guide

The `aws_wellarchitected_notification` table in Steampipe provides you with information about notifications within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query notification-specific details, including the associated workload, notification type, and status. You can utilize this table to gather insights on notifications, such as notifications associated with a specific workload, the status of notifications, and more. The schema outlines the various attributes of the notification for you, including the workload ID, notification type, status, and associated metadata.

## Examples

### List notifications for workloads where lens version is upgraded
Discover the segments that received notifications due to an upgrade in their lens version. This can help in understanding the distribution and impact of the new lens version across different workloads.

```sql+postgres
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  type = 'LENS_VERSION_UPGRADED';
```

```sql+sqlite
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  type = 'LENS_VERSION_UPGRADED';
```

### List notifications for workloads where lens version is deprecated
Identify instances where notifications are related to workloads operating on deprecated lens versions. This is useful for staying updated on potential system vulnerabilities and planning necessary updates.

```sql+postgres
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  type = 'LENS_VERSION_DEPRECATED';
```

```sql+sqlite
select
  workload_name,
  lens_alias,
  lens_arn,
  current_lens_version,
  latest_lens_version
from
  aws_wellarchitected_notification
where
  type = 'LENS_VERSION_DEPRECATED';
```

### Check if there is a notification for a particular workload
Determine if a specific workload has a notification by comparing the current and latest lens versions. This allows for timely updates and avoids potential issues caused by outdated lenses.

```sql+postgres
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

```sql+sqlite
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