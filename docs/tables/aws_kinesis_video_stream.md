# Table: aws_kinesis_video_stream

Amazon Kinesis Video Streams makes it easy to securely stream video from connected devices to AWS for analytics, machine learning (ML), playback, and other processing.

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
