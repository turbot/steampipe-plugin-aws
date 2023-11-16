---
title: "Table: aws_kinesis_firehose_delivery_stream - Query AWS Kinesis Firehose Delivery Stream using SQL"
description: "Allows users to query AWS Kinesis Firehose Delivery Stream data, providing detailed information about each delivery stream in the AWS account."
---

# Table: aws_kinesis_firehose_delivery_stream - Query AWS Kinesis Firehose Delivery Stream using SQL

The `aws_kinesis_firehose_delivery_stream` table in Steampipe provides information about each Kinesis Firehose Delivery Stream within AWS. This table allows DevOps engineers to query delivery stream-specific details, including the delivery stream name, status, creation time, and associated metadata. Users can utilize this table to gather insights on delivery streams, such as the status of the stream, its creation time, and the destinations configured for the stream. The schema outlines the various attributes of the Kinesis Firehose Delivery Stream, including the delivery stream ARN, delivery stream type, delivery stream status, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kinesis_firehose_delivery_stream` table, you can use the `.inspect aws_kinesis_firehose_delivery_stream` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the delivery stream. This can be used to join this table with other tables that rely on ARN for identification.
- `name`: The name of the delivery stream. This is a unique identifier for the delivery stream and can be used to join with other tables that require the stream name.
- `delivery_stream_status`: The current status of the delivery stream. This helps in identifying the operational status of the delivery stream.

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
