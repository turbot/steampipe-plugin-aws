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


### Get a specific consumers

```sql
select
  consumer_name,
  consumer_arn,
  stream_arn
from
  aws_kinesis_consumer
where
  consumer_arn = 'arn:aws:kinesis:us-east-1:986325076436:stream/turbot-data-stream/consumer/turbot-consumer:1616584220';
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
  consumer_status = 'CREATING' or consumer_status = 'DELETING';
```


### List of consumers with particular stream

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn,
  stream_arn
from
  aws_kinesis_consumer
where
  stream_arn = 'arn:aws:kinesis:us-east-1:986325076436:stream/turbot-data-stream';
```


### List of active consumers

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where
  consumer_status = 'ACTIVE';
```