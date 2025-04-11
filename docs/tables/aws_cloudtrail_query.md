---
title: "Steampipe Table: aws_cloudtrail_query - Query AWS CloudTrail using SQL"
description: "Allows users to query AWS CloudTrail events for a detailed view of account activity, including actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_query - Query AWS CloudTrail using SQL

The AWS CloudTrail is a service that enables governance, compliance, operational auditing, and risk auditing of your AWS account. It allows you to log, continuously monitor, and retain account activity related to actions across your AWS infrastructure. With CloudTrail, you can conduct security analysis, track changes to your AWS resources, and aid in compliance reporting.

## Table Usage Guide

The `aws_cloudtrail_query` table in Steampipe provides you with information about CloudTrail events within AWS. This table allows you, as a DevOps engineer, to query event-specific details, including the identity of the API caller, the time of the API call, the source IP address of the API caller, and the request parameters made. You can utilize this table to gather insights on account activity, such as actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services. The schema outlines the various attributes of the CloudTrail event for you, including the event name, event time, event source, and associated tags.

## Examples

### Basic info
Gain insights into the status and efficiency of your AWS CloudTrail queries, including the number of events matched and scanned, to optimize resource usage and improve query performance. This can be particularly useful for troubleshooting and auditing purposes.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which AWS CloudTrail queries have failed to gain insights into potential issues or bottlenecks within your system.

```sql+postgres
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

```sql+sqlite
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
Explore the relationship between specific queries and their corresponding event data stores in AWS CloudTrail, providing insights into the status, multi-region capability, and termination protection of these data stores.

```sql+postgres
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

```sql+sqlite
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
Identify AWS CloudTrail queries that have been created within the last three days, allowing you to monitor recent query activity and understand their execution times.

```sql+postgres
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

```sql+sqlite
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
  creation_time <= datetime('now', '-3 day');
```