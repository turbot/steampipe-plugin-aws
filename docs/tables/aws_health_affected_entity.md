---
title: "Steampipe Table: aws_health_affected_entity - Query AWS Health Affected Entities using SQL"
description: "Allows users to query Affected Entities in AWS Health. The `aws_health_affected_entity` table provides comprehensive details about each entity affected by AWS Health events. It can be utilized to gain insights into the health status of AWS resources, allowing for proactive monitoring and maintenance."
folder: "Health"
---

# Table: aws_health_affected_entity - Query AWS Health Affected Entities using SQL

AWS Health Affected Entities are items that are impacted by events. These entities can include AWS accounts, services, and other AWS resources. They provide valuable information to help you analyze and respond to service health notifications, making it easier to manage events and improve the resiliency of your applications.

## Table Usage Guide

The `aws_health_affected_entity` table in Steampipe provides you with detailed information about entities affected by AWS Health events. This table allows you, as a system administrator or DevOps engineer, to query entity-specific details, including entity ARN, event ARN, status, last updated time, and associated tags. You can utilize this table to gain insights into the health status of AWS resources, enabling proactive monitoring and maintenance. The schema outlines the various attributes of the affected entity for you, such as entity ARN, event ARN, entity value, last updated time, status, and tags.

## Examples

### Basic info
Explore which AWS health events have impacted your resources. This can help you understand the status and last updated time of each affected entity, providing valuable insights for incident response and recovery.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which entities within the AWS Health service are functioning without impairment. This information can be useful for ensuring system stability and identifying areas of consistent performance.

```sql+postgres
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

```sql+sqlite
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

### Get health event details of each entity
Explore the health status of various entities to gain insights into any potential issues or disruptions. This is particularly useful for proactive problem management and maintaining optimal system health.

```sql+postgres
select
  e.arn,
  e.entity_url,
  e.event_arn,
  v.event_type_category,
  v.event_type_code,
  v.service
from
  aws_health_affected_entity as e,
  aws_health_event as v;
```

```sql+sqlite
select
  e.arn,
  e.entity_url,
  e.event_arn,
  v.event_type_category,
  v.event_type_code,
  v.service
from
  aws_health_affected_entity as e,
  aws_health_event as v;
```