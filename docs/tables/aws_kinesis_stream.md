# Table: aws_kinesis_stream

A stream captures and transports data records that are continuously emitted from different data sources or producers. Scale-out within a stream is explicitly supported by means of shards, which are uniquely identified groups of data records in a stream.

## Examples

### Basic info

```sql
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream;
```


### List streams that are not active

```sql
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  stream_status != 'ACTIVE';
```


### List streams that have no consumers

```sql
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  consumer_count = 0;
```


### List streams that are not encrypted

```sql
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type = 'NONE';
```


### List streams that are not encrypted using CMK

```sql
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type != 'NONE'
  and key_id = 'alias/aws/kinesis';
```
