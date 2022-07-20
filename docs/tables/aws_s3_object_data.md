# Table: aws_s3_object_data

The `aws_s3_object_data` table provides access to the data stored in a S3 object.

The data is serialized into a string if it contains valid UTF8 bytes, otherwise it is encoded into Base64, as defined in [RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648).

To list objects, you must mention the `key` and the name of the container `bucket` which contains the objects.

> Note: Using this table adds to cost to your monthly bill from AWS. Optimizations have been put in place to minimize the impact as much as possible. Please refer to [AWS S3 Pricing](https://aws.amazon.com/s3/pricing/) to understand the cost implications.

## Examples

### Basic info

```sql
select
  key,
  bucket,
  content_type
from
  aws_s3_object_data
where
  bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### Get a server side encrypted object

```sql
select
  key,
  bucket,
  content_type,
  data
from
  aws_s3_object_data
where
  bucket = 'user_uploads'
  and key = 'avatar_9ac3097c-1e56-4108-b92e-226a3f4caeb8'
  and sse_customer_key = 'K01iUWVUaFdtWnE0dDd3OXokQyZGKUpATmNSZlVqWG4=';
```

### Parse object data into `jsonb`

```sql
select
  key,
  bucket,
  data::jsonb
from
  aws_s3_object_data
where
  bucket = 'logs'
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/db_logs.json.gz';
```

### Process `jsonb` data in objects

```sql
select
  event ->> 'level' as level,
  event ->> 'severity' as severity,
  event ->> 'message' as event_message,
  event ->> 'data' as event_data,
  event ->> 'timestamp' as timestamp 
from
  aws_s3_object_data,
  jsonb_array_elements((data::jsonb) -> 'events') as event 
where
  bucket = 'logs' 
  and key = 'logs/application_logs/2020/11/04/14/40/dashboard/auth_logs.json.gz' 
  and event ->> 'level' = 'error';
```

### Export binary `data` by converting back from `base64`

```sql
select
  lo_export(decode(data, 'base64'), '/app_data/user_data/avatars/' || key || ".png") 
from
  aws_s3_object_data 
where
  bucket = 'user_uploads' 
  and key = 'avatar_9ac3097c-1e56-4108-b92e-226a3f4caeb8' 
  and sse_customer_key = 'K01iUWVUaFdtWnE0dDd3OXokQyZGKUpATmNSZlVqWG4=';
```
