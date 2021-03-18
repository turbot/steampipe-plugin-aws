# Table: aws_kinesis_consumer

A consumer, known as an Amazon Kinesis Data Streams application, is an application that you build to read and process data records from Kinesis data streams.

## Examples

### Basic info

```sql
select
  consumer_name,
  consumer_arn,
  consumer_status,
  stream_arn
from
  aws_kinesis_consumer;
```

### List of consumers which can't read data

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where consumer_status = 'CREATING' or consumer_status = 'DELETING';
```