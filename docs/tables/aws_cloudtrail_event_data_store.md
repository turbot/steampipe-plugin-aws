---
title: "Table: aws_cloudtrail_event_data_store - Query AWS CloudTrail Event Data using SQL"
description: "Allows users to query AWS CloudTrail Event Data, providing information about API activity in AWS accounts. This includes details about API calls, logins, and other events captured by AWS CloudTrail."
---

# Table: aws_cloudtrail_event_data_store - Query AWS CloudTrail Event Data using SQL

The `aws_cloudtrail_event_data_store` table in Steampipe provides information about API activity in AWS accounts. This includes details about API calls, logins, and other events captured by AWS CloudTrail. This table allows DevOps engineers to query event-specific details, including event names, event sources, and related metadata. Users can utilize this table to gather insights on API activity, such as identifying unusual API calls, tracking login activity, and monitoring changes to AWS resources. The schema outlines the various attributes of the CloudTrail event, including the event ID, event time, event name, and user identity.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_event_data_store` table, you can use the `.inspect aws_cloudtrail_event_data_store` command in Steampipe.

**Key columns**:

- `event_id`: The unique identifier for the event. This can be used to join this table with other tables to get detailed information about specific events.
- `event_name`: The name of the event. This can be useful for filtering events of a specific type.
- `event_source`: The service that the request was made to. This can help in identifying which AWS service the event is associated with.

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
