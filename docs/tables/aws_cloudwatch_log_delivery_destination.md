---
title: "Steampipe Table: aws_cloudwatch_log_delivery_destination - Query AWS CloudWatch Log Delivery Destinations using SQL"
description: "Allows users to query AWS CloudWatch Log Delivery Destinations, providing information about destination configurations for vended log delivery."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_delivery_destination - Query AWS CloudWatch Log Delivery Destinations using SQL

AWS CloudWatch Log Delivery Destination represents a destination configuration for vended log delivery. These destinations define where logs can be sent, such as S3 buckets, CloudWatch Logs, or Firehose delivery streams. Delivery destinations are part of CloudWatch Logs Vended Log Delivery, which enables near real-time log delivery from AWS services like CloudFront.

## Table Usage Guide

The `aws_cloudwatch_log_delivery_destination` table in Steampipe provides you with information about destination configurations for vended log delivery in AWS CloudWatch Logs. This table allows you, as a DevOps engineer, security analyst, or cloud administrator, to query destination-specific details, including the destination type, ARN, output format, and associated policy. You can utilize this table to gather insights on log delivery destinations, such as identifying all S3 buckets configured as log destinations, understanding output formats in use, and managing delivery policies. The schema outlines the various attributes of the delivery destination for you, including the name, ARN, destination resource ARN, delivery destination type, output format, and associated tags.

## Examples

### Basic info
Explore the configuration details of CloudWatch Log delivery destinations in your AWS account to understand where logs can be sent.

```sql+postgres
select
  name,
  arn,
  destination_resource_arn,
  delivery_destination_type,
  output_format
from
  aws_cloudwatch_log_delivery_destination;
```

```sql+sqlite
select
  name,
  arn,
  destination_resource_arn,
  delivery_destination_type,
  output_format
from
  aws_cloudwatch_log_delivery_destination;
```

### Count destinations by type
Count the number of log delivery destinations grouped by their destination type to get a quick summary of log delivery architecture.

```sql+postgres
select
  delivery_destination_type,
  count(*) as destination_count
from
  aws_cloudwatch_log_delivery_destination
group by
  delivery_destination_type
order by
  destination_count desc;
```

```sql+sqlite
select
  delivery_destination_type,
  count(*) as destination_count
from
  aws_cloudwatch_log_delivery_destination
group by
  delivery_destination_type
order by
  destination_count desc;
```

### Get S3 buckets used as log destinations
Retrieve a list of S3 buckets used as destinations for vended log delivery to validate where your logs are being stored.

```sql+postgres
select
  d.name as destination_name,
  d.arn as destination_arn,
  d.destination_resource_arn,
  b.name as bucket_name,
  b.region,
  d.output_format
from
  aws_cloudwatch_log_delivery_destination d
join
  aws_s3_bucket b on d.destination_resource_arn like concat('arn:aws:s3:::', b.name, '%')
where
  d.delivery_destination_type = 'S3';
```

```sql+sqlite
select
  d.name as destination_name,
  d.arn as destination_arn,
  d.destination_resource_arn,
  b.name as bucket_name,
  b.region,
  d.output_format
from
  aws_cloudwatch_log_delivery_destination d
join
  aws_s3_bucket b on d.destination_resource_arn like 'arn:aws:s3:::' || b.name || '%'
where
  d.delivery_destination_type = 'S3';
```

### Count destinations by output format
Analyze how many destinations use each output format to understand log formatting standards across destinations.

```sql+postgres
select
  output_format,
  delivery_destination_type,
  count(*) as destination_count
from
  aws_cloudwatch_log_delivery_destination
group by
  output_format,
  delivery_destination_type
order by
  output_format,
  destination_count desc;
```

```sql+sqlite
select
  output_format,
  delivery_destination_type,
  count(*) as destination_count
from
  aws_cloudwatch_log_delivery_destination
group by
  output_format,
  delivery_destination_type
order by
  output_format,
  destination_count desc;
```

### View destinations with defined policies
Identify destinations that have explicit policies defined to assess access control configurations.

```sql+postgres
select
  name,
  delivery_destination_type,
  policy
from
  aws_cloudwatch_log_delivery_destination
where
  policy is not null;
```

```sql+sqlite
select
  name,
  delivery_destination_type,
  policy
from
  aws_cloudwatch_log_delivery_destination
where
  policy is not null;
```

### Identify which sources are using each destination
Identify the source destination mappings.

```sql+postgres
select
  dest.name as destination_name,
  dest.delivery_destination_type,
  dest.output_format,
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id
from
  aws_cloudwatch_log_delivery_destination dest
join
  aws_cloudwatch_log_delivery d on dest.arn = d.delivery_destination_arn
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
order by
  dest.name, s.service;
```

```sql+sqlite
select
  dest.name as destination_name,
  dest.delivery_destination_type,
  dest.output_format,
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id
from
  aws_cloudwatch_log_delivery_destination dest
join
  aws_cloudwatch_log_delivery d on dest.arn = d.delivery_destination_arn
join
  aws_cloudwatch_log_delivery_source s on d.delivery_source_name = s.name
order by
  dest.name, s.service;
```