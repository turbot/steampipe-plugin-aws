---
title: "Table: aws_kinesis_consumer - Query AWS Kinesis Consumers using SQL"
description: "Allows users to query AWS Kinesis Consumers. This table provides information about Kinesis Consumers within AWS Kinesis Data Streams. It enables users to gather insights on consumers such as consumer ARN, creation timestamp, stream ARN and more."
---

# Table: aws_kinesis_consumer - Query AWS Kinesis Consumers using SQL

The `aws_kinesis_consumer` table in Steampipe provides information about Kinesis Consumers within AWS Kinesis Data Streams. This table allows DevOps engineers to query consumer-specific details, including consumer ARN, creation timestamp, and associated stream ARN. Users can utilize this table to gather insights on consumers, such as details of the consumer, consumer status, and more. The schema outlines the various attributes of the Kinesis Consumer, including the consumer name, consumer status, consumer ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kinesis_consumer` table, you can use the `.inspect aws_kinesis_consumer` command in Steampipe.

**Key columns**:

- `consumer_arn`: The Amazon Resource Name (ARN) of the Kinesis Consumer. This can be used to join with other tables that contain Kinesis Consumer ARNs.
- `consumer_name`: The name of the Kinesis Consumer. This can be used to join with other tables that contain Kinesis Consumer names.
- `stream_arn`: The ARN of the Kinesis Data Stream that the consumer is associated with. This can be used to join with other tables that contain Kinesis Data Stream ARNs.

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


### List consumers which are not in the active state

```sql
select
  consumer_name,
  consumer_status,
  consumer_arn
from
  aws_kinesis_consumer
where
  consumer_status != 'ACTIVE'
```
