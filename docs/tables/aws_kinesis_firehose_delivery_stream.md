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


### List firehose delivery streams which are not active

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


### List firehose delivery streams which are not encrypted

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


### List firehose delivery streams which are not updated in last 7 days

```sql
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  last_update_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  last_update_timestamp < (current_date - interval '7' day);
```


### List firehose delivery streams for a specific delivery stream type

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


### List firehose delivery streams with failure description

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