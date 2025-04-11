---
title: "Steampipe Table: aws_cloudwatch_log_destination - Query AWS CloudWatch Log Destinations using SQL"
description: "Allows users to query AWS CloudWatch Log Destinations, providing information about destination configurations for vended log delivery."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_destination - Query AWS CloudWatch Log Destinations using SQL

AWS CloudWatch Log Destination represents a destination configuration for vended log delivery. These destinations define where logs can be sent, such as S3 buckets, CloudWatch Logs, or Firehose delivery streams. Destinations are part of CloudWatch Logs Vended Log Delivery, which enables near real-time log delivery from AWS services like CloudFront.

## Table Usage Guide

The `aws_cloudwatch_log_destination` table in Steampipe provides you with information about destination configurations for vended log delivery in AWS CloudWatch Logs. This table allows you, as a DevOps engineer, security analyst, or cloud administrator, to query destination-specific details, including the destination type, ARN, and configuration. You can utilize this table to gather insights on log delivery destinations, such as identifying all S3 buckets configured as log destinations, understanding which CloudWatch Log groups are receiving logs, and monitoring your overall logging architecture. The schema outlines the various attributes of the delivery destination for you, including the destination name, ARN, type, creation time, and associated tags.

## Examples

### Basic info
Explore the configuration details of CloudWatch Log delivery destinations in your AWS account to understand where logs can be sent.

```sql+postgres
select
  destination_name,
  arn,
  role_arn,
  target_arn,
  creation_time
from
  aws_cloudwatch_log_destination;
```

```sql+sqlite
select
  destination_name,
  arn,
  role_arn,
  target_arn,
  creation_time
from
  aws_cloudwatch_log_destination;
```

### Classify destinations by target type
Determine the types of destinations used for log delivery by classifying based on the target ARN.

```sql+postgres
select
  destination_name,
  case
    when target_arn like 'arn:aws:s3:%' then 'S3'
    when target_arn like 'arn:aws:firehose:%' then 'Firehose'
    when target_arn like 'arn:aws:logs:%' then 'CloudWatch Logs'
    else 'Other'
  end as target_type,
  target_arn
from
  aws_cloudwatch_log_destination
order by
  target_type;
```

```sql+sqlite
select
  destination_name,
  case
    when target_arn like 'arn:aws:s3:%' then 'S3'
    when target_arn like 'arn:aws:firehose:%' then 'Firehose'
    when target_arn like 'arn:aws:logs:%' then 'CloudWatch Logs'
    else 'Other'
  end as target_type,
  target_arn
from
  aws_cloudwatch_log_destination
order by
  target_type;
```

### List S3 buckets configured as log destinations
Identify S3 buckets that are used as targets for log delivery in your AWS environment.

```sql+postgres
select
  d.destination_name,
  d.arn as destination_arn,
  d.target_arn,
  b.name as bucket_name,
  b.region,
  b.creation_date
from
  aws_cloudwatch_log_destination d
join
  aws_s3_bucket b on d.target_arn like concat('arn:aws:s3:::', b.name, '%')
where
  d.target_arn like 'arn:aws:s3:%';
```

```sql+sqlite
select
  d.destination_name,
  d.arn as destination_arn,
  d.target_arn,
  b.name as bucket_name,
  b.region,
  b.creation_date
from
  aws_cloudwatch_log_destination d
join
  aws_s3_bucket b on d.target_arn like 'arn:aws:s3:::' || b.name || '%'
where
  d.target_arn like 'arn:aws:s3:%';
```

### List recently created destinations
Track recent changes in your logging setup by viewing recently created log destinations.

```sql+postgres
select
  destination_name,
  target_arn,
  role_arn,
  creation_time
from
  aws_cloudwatch_log_destination
where
  creation_time > now() - interval '7 days'
order by
  creation_time desc;
```

```sql+sqlite
select
  destination_name,
  target_arn,
  role_arn,
  creation_time
from
  aws_cloudwatch_log_destination
where
  creation_time > datetime('now', '-7 days')
order by
  creation_time desc;
```

### Trace full delivery pipeline from source to destination
Get end-to-end visibility of the log delivery pipeline across source, delivery, and destination.

```sql+postgres
select
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id,
  dst.destination_name,
  dst.target_arn
from
  aws_cloudwatch_log_delivery_source s
join
  aws_cloudwatch_log_delivery d on s.name = d.delivery_source_name
join
  aws_cloudwatch_log_destination dst on d.delivery_destination_arn = dst.arn
order by
  s.service, s.name;
```

```sql+sqlite
select
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id,
  dst.destination_name,
  dst.target_arn
from
  aws_cloudwatch_log_delivery_source s
join
  aws_cloudwatch_log_delivery d on s.name = d.delivery_source_name
join
  aws_cloudwatch_log_destination dst on d.delivery_destination_arn = dst.arn
order by
  s.service, s.name;
```