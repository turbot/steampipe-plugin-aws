---
title: "Steampipe Table: aws_health_event - Query AWS Health Events using SQL"
description: "Allows users to query AWS Health Events to retrieve information about events that affect your AWS services and accounts."
folder: "Health"
---

# Table: aws_health_event - Query AWS Health Events using SQL

The AWS Health Events service provides personalized information about events that can affect your AWS infrastructure, guiding your through scheduled changes, and accelerating troubleshooting. It aggregates and consolidates alerts and notifications from multiple AWS services into a single, easy-to-search interface. AWS Health Events also integrates with AWS Organizations to provide a centralized view of health events across all your accounts.

## Table Usage Guide

The `aws_health_event` table in Steampipe provides you with information about AWS Health Events. These events give you timely information about service disruptions, scheduled changes, and other important AWS-related events that can affect your services and accounts. This table allows you, as a DevOps engineer, to query event-specific details, including event type, start time, end time, and affected services. You can utilize this table to monitor your AWS services, understand the impact of AWS events, and plan necessary actions accordingly. The schema outlines the various attributes of the AWS Health Event for you, including the event ARN, event type category, service, region, start time, and end time.

## Examples

### Basic info
Explore which events have occurred within your AWS services, pinpointing their specific locations and durations. This can help in assessing the overall health and performance of your AWS infrastructure.

```sql+postgres
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

```sql+sqlite
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
Gain insights into upcoming events across your AWS services to better prepare and manage your resources. This query is useful for proactive planning and risk mitigation.

```sql+postgres
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

```sql+sqlite
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
Explore the various events related to the EC2 service by analyzing their start and end times, status, and other relevant details. This can help determine any patterns or irregularities in the service's operation, aiding in proactive management and troubleshooting.

```sql+postgres
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

```sql+sqlite
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
Gain insights into the specific events occurring within a particular availability zone, enabling you to monitor and manage service health more effectively.

```sql+postgres
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

```sql+sqlite
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