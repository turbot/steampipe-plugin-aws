---
title: "Table: aws_health_event - Query AWS Health Events using SQL"
description: "Allows users to query AWS Health Events to retrieve information about events that affect your AWS services and accounts."
---

# Table: aws_health_event - Query AWS Health Events using SQL

The `aws_health_event` table in Steampipe provides information about AWS Health Events. AWS Health Events provide timely information about service disruptions, scheduled changes, and other important AWS-related events that can affect your services and accounts. This table allows DevOps engineers to query event-specific details, including event type, start time, end time, and affected services. Users can utilize this table to monitor their AWS services, understand the impact of AWS events, and plan necessary actions accordingly. The schema outlines the various attributes of the AWS Health Event, including the event ARN, event type category, service, region, start time, and end time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_health_event` table, you can use the `.inspect aws_health_event` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the event. This can be used to join this table with other tables that contain event ARN.
- `service`: The AWS service that is affected by the event. This can be used to join this table with other tables that contain service information.
- `event_type_category`: The category of the event (for example, 'issue', 'scheduledChange', or 'accountNotification'). This can be useful for filtering events based on their type.

## Examples

### Basic info

```sql
select
  arn,
  availability_zone,
  start_time,
  end_time,
  event_type_category,
  event_type_code,
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
  event_type_category,
  event_type_code,
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
  event_type_category,
  event_type_code,
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
  event_type_category,
  event_type_code,
  event_scope_code,
  status_code,
  service
from
  aws_health_event
where
  availability_zone = 'us-east-1a';
```
