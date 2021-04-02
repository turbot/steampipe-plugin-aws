# Table: aws_kinesis_firehose_delivery_stream

The AWS Kinesis Firehose Delivery Stream resource delivers real-time streaming data to an Amazon Simple Storage Service (Amazon S3), Amazon Redshift or Amazon Elasticsearch Service (Amazon ES) destination.

## Examples

### Basic info

```sql
select
  delivery_stream_name,
  arn,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream;
```


### List inactive delivery streams

```sql
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_status != 'ACTIVE';
```


### List delivery streams that are not encrypted

```sql
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type,
  delivery_stream_encryption_configuration ->> 'Status' as encryption_status
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_encryption_configuration ->> 'Status' = 'DISABLED';
```


### List delivery streams for a specific delivery stream type

```sql
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_type = 'DirectPut';
```


### List delivery streams with at least one failure

```sql
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  delivery_stream_type,
  failure_description
from
  aws_kinesis_firehose_delivery_stream
where
  failure_description is not null;
```
