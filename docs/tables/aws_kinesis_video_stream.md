---
title: "Steampipe Table: aws_kinesis_video_stream - Query AWS Kinesis Video Streams using SQL"
description: "Allows users to query Kinesis Video Streams to obtain metadata about each stream, including the stream's ARN, creation time, status, and other information."
folder: "Kinesis"
---

# Table: aws_kinesis_video_stream - Query AWS Kinesis Video Streams using SQL

The AWS Kinesis Video Streams service enables you to securely ingest, process, and store video at any scale for applications that power robots, smart cities, industrial automation, and more. It provides SDKs that you can use to stream video from devices to AWS for playback, storage, analytics, machine learning, and other processing. It also automatically provisions and elastically scales all the infrastructure needed to ingest streaming video data from millions of devices.

## Table Usage Guide

The `aws_kinesis_video_stream` table in Steampipe provides you with information about the Kinesis Video Streams within AWS Kinesis. This table allows you, as a DevOps engineer, to query stream-specific details, including the stream ARN, creation time, status, and more. You can utilize this table to gather insights on streams, such as stream status, data retention period, and the version of the stream. The schema outlines the various attributes of the Kinesis Video Stream for you, including the stream name, ARN, version, status, creation time, and data retention period.

## Examples

### Basic info
Discover the segments that are currently active in your AWS Kinesis video stream, including their creation times and regions. This can aid in monitoring the status and distribution of your video streams across different regions.

```sql+postgres
select
  stream_name,
  stream_arn,
  status,
  creation_time,
  region
from
  aws_kinesis_video_stream;
```

```sql+sqlite
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
Determine the areas in which video streams are not in an active state. This is beneficial in identifying potential issues or disruptions in your video streaming service.

```sql+postgres
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

```sql+sqlite
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
Identify unencrypted video streams to assess potential security vulnerabilities. This query can be used to pinpoint areas where stronger encryption methods may need to be implemented for improved data protection.

```sql+postgres
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

```sql+sqlite
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
  substr(kms_key_id, instr(kms_key_id, ':')+5) = 'alias/aws/kinesisvideo';
```


### List video streams with data retention period < 7 days (represented as hours in query)
Determine the areas in which video streams have a data retention period of less than a week. This is useful for identifying potential data loss risks due to short retention periods.

```sql+postgres
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

```sql+sqlite
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