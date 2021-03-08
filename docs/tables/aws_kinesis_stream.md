# Table: aws_kinesis_stream

Amazon Kinesis Streams is a massively scalable and durable real-time streaming service.

## Examples

### List all the streams in your account

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


### List all the streams that are not in Active state

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


### List all the streams that have no consumers

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
	consumer_count = '0';
```


### List all the streams which are not encrypted

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


### List all the encrypted streams which are not encrypted using customer managed key

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
	And key_id = 'alias/aws/kinesis';
```