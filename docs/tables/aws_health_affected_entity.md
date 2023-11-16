---
title: "Table: aws_health_affected_entity - Query AWS Health Affected Entities using SQL"
description: "Allows users to query Affected Entities in AWS Health. The `aws_health_affected_entity` table provides comprehensive details about each entity affected by AWS Health events. It can be utilized to gain insights into the health status of AWS resources, allowing for proactive monitoring and maintenance."
---

# Table: aws_health_affected_entity - Query AWS Health Affected Entities using SQL

The `aws_health_affected_entity` table in Steampipe provides detailed information about entities affected by AWS Health events. This table allows system administrators and DevOps engineers to query entity-specific details, including entity ARN, event ARN, status, last updated time, and associated tags. Users can utilize this table to gain insights into the health status of AWS resources, enabling proactive monitoring and maintenance. The schema outlines the various attributes of the affected entity, such as entity ARN, event ARN, entity value, last updated time, status, and tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_health_affected_entity` table, you can use the `.inspect aws_health_affected_entity` command in Steampipe.

### Key columns:

- `entity_arn`: The unique identifier for the affected entity. This can be used to join this table with other tables that contain entity ARN information.
- `event_arn`: The unique identifier for the event affecting the entity. This can be used to join this table with other tables that contain event ARN information.
- `last_updated_time`: The time when the entity was last updated. This can be used to sort or filter entities based on the recency of their updates.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  e.arn,
  e.entity_url,
  e.event_arn,
  v.event_start_time,
  v.event_end_time,
  v.event_type_category,
  v.event_type_code
  v.service
from
  aws_health_affected_entity as e,
  aws_health_event as v;
```