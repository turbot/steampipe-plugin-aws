---
title: "Steampipe Table: aws_cloudwatch_log_delivery - Query AWS CloudWatch Log Deliveries using SQL"
description: "Allows users to query AWS CloudWatch Log Deliveries, providing information about delivery configurations for vended log delivery."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_delivery - Query AWS CloudWatch Log Deliveries using SQL

AWS CloudWatch Log Delivery represents a delivery configuration for vended log delivery, connecting a delivery source to a destination. These deliveries define how logs from AWS services like CloudFront are delivered to destinations such as S3 buckets, CloudWatch Logs, or Firehose delivery streams, enabling near real-time log analysis and storage.

## Table Usage Guide

The `aws_cloudwatch_log_delivery` table in Steampipe provides you with information about delivery configurations for vended log delivery in AWS CloudWatch Logs. This table allows you, as a DevOps engineer, security analyst, or cloud administrator, to query delivery-specific details, including the associated source, destination, and status. You can utilize this table to gather insights on log delivery configurations, such as identifying where specific services are sending their logs, determining which destinations are receiving logs, and monitoring the status of your log delivery pipelines. The schema outlines the various attributes of the delivery configuration for you, including the delivery ID, source name, destination ARN, creation time, and associated tags.

## Examples

### Basic info
Gain insights into your log delivery configurations to understand how logs are routed within your AWS environment.

```sql+postgres
select
  id,
  arn,
  delivery_source_name,
  delivery_destination_arn,
  delivery_destination_type
from
  aws_cloudwatch_log_delivery;
```

```sql+sqlite
select
  id,
  arn,
  delivery_source_name,
  delivery_destination_arn,
  delivery_destination_type
from
  aws_cloudwatch_log_delivery;
```

### Count deliveries by destination type
Analyze the distribution of deliveries across different destination types to understand your log delivery targets.

```sql+postgres
select
  delivery_destination_type,
  count(*) as delivery_count
from
  aws_cloudwatch_log_delivery
group by
  delivery_destination_type;
```

```sql+sqlite
select
  delivery_destination_type,
  count(*) as delivery_count
from
  aws_cloudwatch_log_delivery
group by
  delivery_destination_type;
```

### Identify CloudFront deliveries to CloudWatch Logs
Find all CloudFront sources configured to deliver logs to CloudWatch Log groups.

```sql+postgres
select
  d.id as delivery_id,
  s.name as source_name,
  s.service,
  d.delivery_destination_type,
  split_part(d.delivery_destination_arn, ':', 7) as log_group_name
from
  aws_cloudwatch_log_delivery d
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
where
  s.service = 'cloudfront'
  and d.delivery_destination_type = 'CloudWatchLogs';
```

```sql+sqlite
select
  d.id as delivery_id,
  s.name as source_name,
  s.service,
  d.delivery_destination_type,
  substr(
    d.delivery_destination_arn,
    instr(d.delivery_destination_arn, 'log-group:') + 10
  ) as log_group_name
from
  aws_cloudwatch_log_delivery d
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
where
  s.service = 'cloudfront'
  and d.delivery_destination_type = 'CloudWatchLogs';
```

### Get S3 bucket details for vended log deliveries
List the S3 buckets used as delivery destinations for CloudWatch logs.

```sql+postgres
select
  d.id as delivery_id,
  s.name as source_name,
  s.service,
  b.name as bucket_name,
  b.region,
  b.creation_date
from
  aws_cloudwatch_log_delivery d
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
join
  aws_s3_bucket b on b.arn = d.delivery_destination_arn
where
  d.delivery_destination_type = 'S3';
```

```sql+sqlite
select
  d.id as delivery_id,
  s.name as source_name,
  s.service,
  b.name as bucket_name,
  b.region,
  b.creation_date
from
  aws_cloudwatch_log_delivery d
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
join
  aws_s3_bucket b on b.arn = d.delivery_destination_arn
where
  d.delivery_destination_type = 'S3';
```

### Review full CloudFront logging configuration
Compare logging targets for CloudFront distributions, including S3, CloudWatch Logs, and Firehose.

```sql+postgres
select
  cf.id as distribution_id,
  cf.domain_name,
  s.name as source_name,
  d.id as delivery_id,
  d.delivery_destination_type,
  case
    when d.delivery_destination_type = 'S3' then b.name
    when d.delivery_destination_type = 'CloudWatchLogs' then split_part(d.delivery_destination_arn, ':', 7)
    when d.delivery_destination_type = 'Firehose' then split_part(d.delivery_destination_arn, '/', 2)
    else 'Unknown'
  end as destination_name
from
  aws_cloudfront_distribution cf
join
  aws_cloudwatch_log_delivery_source s
    on s.resource_arns @> to_jsonb(array[cf.arn])
    and s.service = 'cloudfront'
join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name
left join
  aws_s3_bucket b
    on b.arn = d.delivery_destination_arn
    and d.delivery_destination_type = 'S3';
```

```sql+sqlite
select
  cf.id as distribution_id,
  cf.domain_name,
  s.name as source_name,
  d.id as delivery_id,
  d.delivery_destination_type,
  case
    when d.delivery_destination_type = 'S3' then b.name
    when d.delivery_destination_type = 'CloudWatchLogs' then split(d.delivery_destination_arn, ':')[7]
    when d.delivery_destination_type = 'Firehose' then split(d.delivery_destination_arn, '/')[2]
    else 'Unknown'
  end as destination_name
from
  aws_cloudwatch_log_delivery_source s
join
  json_each(s.resource_arns) as ra
    on ra.value = cf.arn
join
  aws_cloudfront_distribution cf
    on true
join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name
left join
  aws_s3_bucket b
    on b.arn = d.delivery_destination_arn
    and d.delivery_destination_type = 'S3'
where
  s.service = 'cloudfront';
```