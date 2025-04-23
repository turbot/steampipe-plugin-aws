---
title: "Steampipe Table: aws_kinesis_consumer - Query AWS Kinesis Consumers using SQL"
description: "Allows users to query AWS Kinesis Consumers. This table provides information about Kinesis Consumers within AWS Kinesis Data Streams. It enables users to gather insights on consumers such as consumer ARN, creation timestamp, stream ARN and more."
folder: "Kinesis"
---

# Table: aws_kinesis_consumer - Query AWS Kinesis Consumers using SQL

The AWS Kinesis Consumer is an entity within the AWS Kinesis service that allows you to process data from a Kinesis data stream. It retrieves records from the stream and passes them to your application for processing. This enables you to scale your processing resources to accommodate the rate of data flow from your Kinesis stream.

## Table Usage Guide

The `aws_kinesis_consumer` table in Steampipe provides you with information about Kinesis Consumers within AWS Kinesis Data Streams. This table enables you, as a DevOps engineer, to query consumer-specific details, including consumer ARN, creation timestamp, and associated stream ARN. You can utilize this table to gather insights on consumers, such as details of the consumer, consumer status, and more. The schema outlines for you the various attributes of the Kinesis Consumer, including the consumer name, consumer status, consumer ARN, and associated tags.

## Examples

### Basic info
Explore the status of various streaming data consumers to understand their operational state and associated data streams. This is useful for monitoring the health and activity of data streaming services.

```sql+postgres
select
  consumer_name,
  consumer_arn,
  consumer_status,
  stream_arn
from
  aws_kinesis_consumer;
```

```sql+sqlite
select
  consumer_name,
  consumer_arn,
  consumer_status,
  stream_arn
from
  aws_kinesis_consumer;
```


### List consumers which are not in the active state
Determine the areas in which consumers are not currently active, allowing for targeted troubleshooting and system optimization.

```sql+postgres
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where
  consumer_status != 'ACTIVE'
```

```sql+sqlite
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where
  consumer_status != 'ACTIVE'
```