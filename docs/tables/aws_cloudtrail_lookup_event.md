---
title: "Steampipe Table: aws_cloudtrail_lookup_event - Query AWS CloudTrail Lookup Events using SQL"
description: "Allows users to query AWS CloudTrail Lookup Events, providing information about each trail event within AWS CloudTrail. The table can be used to retrieve details such as the event time, event name, resources involved, and much more."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_lookup_event - Query AWS CloudTrail Lookup Events using SQL

AWS CloudTrail Lookup Events is a feature within AWS CloudTrail, a service that provides a record of actions taken by a user, role, or an AWS service in AWS. This feature specifically allows you to look up and retrieve information about the events recorded by CloudTrail.

## Table Usage Guide

The `aws_cloudtrail_lookup_event` table in Steampipe provides you with information about each trail event within AWS CloudTrail. This table allows you, as a DevOps engineer, to query event-specific details, including event time, event name, resources involved, and more. You can utilize this table to gather insights on trail events, such as event source, user identity, and request parameters. The schema outlines the various attributes of the CloudTrail event for you, including the event ID, event version, read only, and associated tags.

**Important notes:**

- For improved performance, it is advised that you use the optional qual `start_time` and `end_time` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use CloudWatch filters. Optional quals are supported for the following columns:
  - `read_only`
  - `event_id`
  - `event_name`
  - `event_source`
  - `resource_name`
  - `resource_type`
  - `access_key_id`
  - `start_time`
  - `end_time`
  - `username`

## Examples

### List events that occurred over the last five minutes

This query is useful for gaining insights into recent activity within your AWS environment. It provides a quick overview of the events that have taken place in the last five minutes, which can be particularly useful for immediate incident response or real-time monitoring.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  username,
  jsonb_pretty(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  start_time = now() - interval '5 minutes'
  and end_time = now();
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  username,
  json(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  start_time = datetime('now', '-5 minutes')
  and end_time = datetime('now');
```

### List all action events, i.e., not ReadOnly that occurred over the last hour

Explore which action events have occurred in the last hour on AWS Cloudtrail. This is useful for identifying recent activities that have potentially altered your system.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  username,
  jsonb_pretty(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  start_time = now()
  and end_time = now() - interval '1 hour'
  and read_only = 'true'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  username,
  json(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  start_time = datetime('now')
  and end_time = datetime('now', '-1 hour')
  and read_only = 'true'
order by
  event_time asc;
```

### List events for a specific service (IAM) that occurred over the last hour

This query allows users to monitor recent activity for a specific service, in this case, AWS's Identity and Access Management (IAM). It is particularly useful for security audits, as it provides a chronological overview of events, including who initiated them and what actions were taken, over the last hour.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  jsonb_pretty(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  event_source = 'iam.amazonaws.com'
  and event_time >= now() - interval '1 hour';
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  json(cloud_trail_event) as cloud_trail_event
from
  aws_cloudtrail_lookup_event
where
  event_source = 'iam.amazonaws.com'
  and event_time >= datetime('now', '-1 hour');
```
