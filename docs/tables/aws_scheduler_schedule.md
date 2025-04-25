---
title: "Steampipe Table: aws_scheduler_schedule - Query AWS EventBridge Schedules using SQL"
description: "Allows users to query AWS EventBridge Scheduler schedules, including metadata, schedule expressions, target configurations, and other details."
folder: "Scheduler"
---

# Table: aws_scheduler_schedule - Query AWS EventBridge Schedules using SQL

The AWS EventBridge Scheduler is a fully managed service that enables automated task execution at specific times or intervals. It allows you to define schedules for running tasks, managing workflows, and automating processes across AWS services. Schedules can be one-time or recurring, using expressions such as `rate`, `cron`, or `at`.

## Table Usage Guide

The `aws_scheduler_schedule` table in Steampipe provides information about schedules configured in AWS EventBridge Scheduler. This table allows DevOps engineers, system administrators, and cloud professionals to query schedule-specific details, such as schedule expressions, start and end times, flexible time windows, target configurations, and states.

The schema outlines key attributes of the schedule, including its name, ARN, associated group, creation date, schedule state, and metadata. You can use this table to monitor schedules, identify configuration details, and optimize automation workflows.

## Examples

### Basic info
Retrieve the name, ARN, group name, and creation date of all schedules in your AWS account.

```sql+postgres
select
  name,
  arn,
  group_name,
  creation_date
from
  aws_scheduler_schedule;
```

```sql+sqlite
select
  name,
  arn,
  group_name,
  creation_date
from
  aws_scheduler_schedule;
```

### List schedules with a specific state (Enabled)
Identify all schedules that are currently enabled.

```sql+postgres
select
  name,
  arn,
  state,
  schedule_expression
from
  aws_scheduler_schedule
where
  state = 'ENABLED';
```

```sql+sqlite
select
  name,
  arn,
  state,
  schedule_expression
from
  aws_scheduler_schedule
where
  state = 'ENABLED';
```

### List schedules with flexible time windows
Fetch schedules that have flexible time windows configured.

```sql+postgres
select
  name,
  flexible_time_window,
  schedule_expression
from
  aws_scheduler_schedule
where
  flexible_time_window is not null;
```

```sql+sqlite
select
  name,
  flexible_time_window,
  schedule_expression
from
  aws_scheduler_schedule
where
  flexible_time_window is not null;
```

### Retrieve schedules with start and end dates
Identify schedules that have both start and end dates defined.

```sql+postgres
select
  name,
  start_date,
  end_date,
  state
from
  aws_scheduler_schedule
where
  start_date is not null
  and end_date is not null;
```

```sql+sqlite
select
  name,
  start_date,
  end_date,
  state
from
  aws_scheduler_schedule
where
  start_date is not null
  and end_date is not null;
```

### List schedules grouped by their schedule group
Fetch schedules and group them based on their schedule group name.

```sql+postgres
select
  group_name,
  count(*) as total_schedules
from
  aws_scheduler_schedule
group by
  group_name;
```

```sql+sqlite
select
  group_name,
  count(*) as total_schedules
from
  aws_scheduler_schedule
group by
  group_name;
```

### Analyze target configurations for specific schedules
Retrieve the target configurations, including the action and resource details, for specific schedules.

```sql+postgres
select
  name,
  target,
  action_after_completion
from
  aws_scheduler_schedule
where
  name = 'my-schedule';
```

```sql+sqlite
select
  name,
  target,
  action_after_completion
from
  aws_scheduler_schedule
where
  name = 'my-schedule';
```

### Fetch schedules created in the last 30 days
Identify newly created schedules within the past month for auditing or review.

```sql+postgres
select
  name,
  creation_date,
  state,
  group_name
from
  aws_scheduler_schedule
where
  creation_date > now() - interval '30 days';
```

```sql+sqlite
select
  name,
  creation_date,
  state,
  group_name
from
  aws_scheduler_schedule
where
  creation_date > datetime('now','-30 days');
```

### Retrieve schedules with specific time zones
Fetch schedules where the schedule expression is evaluated in a specific time zone.

```sql+postgres
select
  name,
  schedule_expression,
  schedule_expression_timezone
from
  aws_scheduler_schedule
where
  schedule_expression_timezone = 'UTC';
```

```sql+sqlite
select
  name,
  schedule_expression,
  schedule_expression_timezone
from
  aws_scheduler_schedule
where
  schedule_expression_timezone = 'UTC';
```
