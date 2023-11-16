---
title: "Table: aws_kinesis_video_stream - Query AWS Kinesis Video Streams using SQL"
description: "Allows users to query Kinesis Video Streams to obtain metadata about each stream, including the stream's ARN, creation time, status, and other information."
---

# Table: aws_kinesis_video_stream - Query AWS Kinesis Video Streams using SQL

The `aws_kinesis_video_stream` table in Steampipe provides information about the Kinesis Video Streams within AWS Kinesis. This table allows DevOps engineers to query stream-specific details, including the stream ARN, creation time, status, and more. Users can utilize this table to gather insights on streams, such as stream status, data retention period, and the version of the stream. The schema outlines the various attributes of the Kinesis Video Stream, including the stream name, ARN, version, status, creation time, and data retention period.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kinesis_video_stream` table, you can use the `.inspect aws_kinesis_video_stream` command in Steampipe.

**Key columns**:

- `stream_name`: The name of the stream. This column can be used to join with other tables that require stream name.
- `stream_arn`: The Amazon Resource Name (ARN) of the stream. This column can be used to join with other tables that require stream ARN.
- `status`: The current status of the stream. This column can be useful for filtering streams based on their status.

## Examples

### Basic info

```sql
select
  stream_name,
  stream_arn,
  status,
  creation_time,
  region
from
  aws_kinesis_video_stream;
```


### List video streams that are not in Active state

```sql
select
  stream_name,
  stream_arn,
  status,
  creation_time,
  region
from
  aws_kinesis_video_stream
where
  status != 'ACTIVE';
```


### List video streams which are not encrypted using CMK

```sql
select
  stream_name,
  stream_arn,
  status,
  kms_key_id,
  creation_time,
  region
from
  aws_kinesis_video_stream
where
  split_part(kms_key_id, ':', 6) = 'alias/aws/kinesisvideo';
```


### List video streams with data retention period < 7 days (represented as hours in query)

```sql
select
  stream_name,
  stream_arn,
  status,
  creation_time,
  data_retention_in_hours,
  region
from
  aws_kinesis_video_stream
where
  data_retention_in_hours < 168;
```
