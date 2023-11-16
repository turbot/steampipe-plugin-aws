---
title: "Table: aws_cloudtrail_query - Query AWS CloudTrail using SQL"
description: "Allows users to query AWS CloudTrail events for a detailed view of account activity, including actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services."
---

# Table: aws_cloudtrail_query - Query AWS CloudTrail using SQL

The `aws_cloudtrail_query` table in Steampipe provides information about CloudTrail events within AWS. This table allows DevOps engineers to query event-specific details, including the identity of the API caller, the time of the API call, the source IP address of the API caller, and the request parameters made. Users can utilize this table to gather insights on account activity, such as actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services. The schema outlines the various attributes of the CloudTrail event, including the event name, event time, event source, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_query` table, you can use the `.inspect aws_cloudtrail_query` command in Steampipe.

**Key columns**:

- `event_time`: This is the date and time the request was made, in coordinated universal time (UTC). It can be used to filter events based on the time they occurred.
- `event_name`: This is the name of the event that occurred. It can be used to filter events based on their type.
- `event_source`: This is the service that the request was made to. It can be used to filter events based on the AWS service they are associated with.

## Examples

### Basic info

```sql
select
  query_id,
  event_data_store_arn,
  query_status,
  query_status,
  creation_time,
  events_matched,
  events_scanned
from
  aws_cloudtrail_query;
```

### List queries that are failed

```sql
select
  query_id,
  event_data_store_arn,
  query_status,
  creation_time,
  query_string,
  execution_time_in_millis
from
  aws_cloudtrail_query
where
  query_status = 'FAILED';
```

### Get event data store details for the queries

```sql
select
  q.query_id as query_id,
  q.event_data_store_arn as event_data_store_arn,
  s.name as event_data_store_name,
  s.status as event_data_store_status,
  s.multi_region_enabled as multi_region_enabled,
  s.termination_protection_enabled as termination_protection_enabled,
  s.updated_timestamp as event_data_store_updated_timestamp
from
  aws_cloudtrail_query as q,
  aws_cloudtrail_event_data_store as s
where
 s.arn = q.event_data_store_arn;
```

## List queries created within the last 3 days

```sql
select
  query_id,
  event_data_store_arn,
  query_status,
  creation_time,
  query_string,
  execution_time_in_millis
from
  aws_cloudtrail_query
where
  creation_time <= now() - interval '3' day;
```