# Table: aws_cloudtrail_query

AWS CloudTrail Lake lets you run SQL-based queries on your events. CloudTrail Lake converts existing events in row-based JSON format to Apache ORC format. ORC is a columnar storage format that is optimized for fast retrieval of data.

**Important notes:**

- The table returns queries from the past 7 days.

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
 s.event_data_store_arn = q.event_data_store_arn;
```

## List queries created within the last 30 days

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
  creation_time <= now() - interval '30' day;
```