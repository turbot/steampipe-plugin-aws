---
title: "Table: aws_kinesis_stream - Query AWS Kinesis Stream using SQL"
description: "Allows users to query AWS Kinesis Stream data, including stream name, status, creation time, and associated tags."
---

# Table: aws_kinesis_stream - Query AWS Kinesis Stream using SQL

The `aws_kinesis_stream` table in Steampipe provides information about Kinesis streams within AWS Kinesis. This table allows DevOps engineers to query stream-specific details, including the stream name, status, creation time, and associated metadata. Users can utilize this table to gather insights on streams, such as stream health, data throughput, and more. The schema outlines the various attributes of the Kinesis stream, including the stream ARN, creation timestamp, number of shards, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kinesis_stream` table, you can use the `.inspect aws_kinesis_stream` command in Steampipe.

**Key columns**:

- `stream_name`: This is the name of the Kinesis stream. It is a key column because it is the unique identifier for each stream and can be used to join this table with other tables.
- `stream_arn`: This is the Amazon Resource Name (ARN) of the Kinesis stream. It is a key column because it provides a unique identifier across all of AWS, which can be used for joining with other tables.
- `status`: This column indicates the current status of the stream (e.g., CREATING, DELETING, ACTIVE, UPDATING). It is a key column because it provides important information about the stream's lifecycle state.

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
