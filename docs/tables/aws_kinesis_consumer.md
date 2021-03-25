# Table: aws_kinesis_consumer

A consumer, known as an Amazon Kinesis Data Streams application, is an application that you build to read and process data records from Kinesis data streams.

## Examples

### List all consumers

```sql
select
  consumer_name,
  consumer_arn,
  consumer_status,
  stream_arn
from
  aws_kinesis_consumer;
```


### Get details of specific consumer

```sql
select
  consumer_name,
  consumer_arn,
  stream_arn
from
  aws_kinesis_consumer
where
  consumer_arn = 'arn:aws:kinesis:us-east-1:986250123456:stream/my-data-stream/consumer/my-consumer:1616584220';
```


### List of consumers which can't read data

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where
  consumer_status != 'ACTIVE';
```


### List of consumers with a particular stream

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn,
  stream_arn
from
  aws_kinesis_consumer
where
  stream_arn = 'arn:aws:kinesis:us-east-1:986250123456:stream/my-data-stream';
```