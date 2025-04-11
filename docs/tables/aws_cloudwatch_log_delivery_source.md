---
title: "Steampipe Table: aws_cloudwatch_log_delivery_source - Query AWS CloudWatch Log Delivery Sources using SQL"
description: "Allows users to query AWS CloudWatch Log Delivery Sources, providing information about source configurations for vended log delivery."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_delivery_source - Query AWS CloudWatch Log Delivery Sources using SQL

AWS CloudWatch Log Delivery Source represents a source configuration for vended log delivery. These sources define AWS services that can send logs to destinations like S3 buckets, CloudWatch Logs, or Firehose delivery streams. Delivery sources are part of CloudWatch Logs Vended Log Delivery, which enables services like CloudFront to deliver logs in near real-time.

## Table Usage Guide

The `aws_cloudwatch_log_delivery_source` table in Steampipe provides you with information about source configurations for vended log delivery in AWS CloudWatch Logs. This table allows you, as a DevOps engineer, security analyst, or cloud administrator, to query source-specific details, including the associated service, log type, and resource ARNs. You can utilize this table to gather insights on log delivery sources, such as identifying all CloudFront distributions configured for vended logs, understanding which services are configured to send logs, and monitoring your overall logging architecture. The schema outlines the various attributes of the delivery source for you, including the source name, ARN, service, creation time, and associated tags.

## Examples

### Basic info
Explore the configuration details of CloudWatch Log delivery sources in your AWS account to understand what services are sending logs to which destinations.

```sql+postgres
select
  name,
  arn,
  service,
  log_type
from
  aws_cloudwatch_log_delivery_source;
```

```sql+sqlite
select
  name,
  arn,
  service,
  log_type
from
  aws_cloudwatch_log_delivery_source;
```

### List CloudFront delivery sources
Identify CloudFront distributions that are configured as log delivery sources.

```sql+postgres
select
  name,
  service,
  log_type,
  resource_arns
from
  aws_cloudwatch_log_delivery_source
where
  service = 'cloudfront';
```

```sql+sqlite
select
  name,
  service,
  log_type,
  resource_arns
from
  aws_cloudwatch_log_delivery_source
where
  service = 'cloudfront';
```

### Count sources by service
Count the number of delivery sources grouped by their originating AWS service.

```sql+postgres
select
  service,
  count(*) as source_count
from
  aws_cloudwatch_log_delivery_source
group by
  service
order by
  source_count desc;
```

```sql+sqlite
select
  service,
  count(*) as source_count
from
  aws_cloudwatch_log_delivery_source
group by
  service
order by
  source_count desc;
```

### Get destinations for each delivery source
Find out where each delivery source is sending its logs by joining with the delivery configuration.

```sql+postgres
select
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id,
  d.delivery_destination_type
from
  aws_cloudwatch_log_delivery_source s
left join
  aws_cloudwatch_log_delivery d
on
  s.name = d.delivery_source_name;
```

```sql+sqlite
select
  s.name as source_name,
  s.service,
  s.log_type,
  d.id as delivery_id,
  d.delivery_destination_type
from
  aws_cloudwatch_log_delivery_source s
left join
  aws_cloudwatch_log_delivery d
on
  s.name = d.delivery_source_name;
```

### Find CloudFront log sources delivering to S3
Get CloudFront sources that are sending vended logs to S3 destinations.

```sql+postgres
select
  s.name as source_name,
  cf.id as distribution_id,
  cf.domain_name,
  d.id as delivery_id,
  d.delivery_destination_type,
  s3.name as s3_bucket_name
from
  aws_cloudwatch_log_delivery_source s
join
  aws_cloudfront_distribution cf
    on s.resource_arns @> to_jsonb(array[cf.arn])
left join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name
left join
  aws_s3_bucket s3
    on s3.arn = d.delivery_destination_arn
where
  s.service = 'cloudfront'
  and d.delivery_destination_type = 'S3';
```

```sql+sqlite
select
  s.name as source_name,
  cf.id as distribution_id,
  cf.domain_name,
  d.id as delivery_id,
  d.delivery_destination_type,
  s3.name as s3_bucket_name
from
  aws_cloudwatch_log_delivery_source s
join
  aws_cloudfront_distribution cf
    on cf.arn in (
      select value
      from json_each(s.resource_arns)
    )
left join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name
left join
  aws_s3_bucket s3
    on s3.arn = d.delivery_destination_arn
where
  s.service = 'cloudfront'
  and d.delivery_destination_type = 'S3';
```

### Compare CloudFront standard vs vended logging
Compare CloudFront distributions with standard logging and those with vended log delivery to evaluate logging configurations.

```sql+postgres
select
  cf.id as distribution_id,
  cf.domain_name,
  case
    when cf.logging->>'Enabled' = 'true' then 'Standard'
    else 'Not configured'
  end as standard_logging_status,
  case
    when s.name is not null then 'Configured'
    else 'Not configured'
  end as vended_logging_status,
  s.name as source_name,
  d.delivery_destination_type
from
  aws_cloudfront_distribution cf
left join
  aws_cloudwatch_log_delivery_source s
    on s.resource_arns @> to_jsonb(array[cf.arn])
    and s.service = 'cloudfront'
left join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name;
```

```sql+sqlite
select
  cf.id as distribution_id,
  cf.domain_name,
  case
    when cf.logging ->> 'Enabled' = 'true' then 'Standard'
    else 'Not configured'
  end as standard_logging_status,
  case
    when s.name is not null then 'Configured'
    else 'Not configured'
  end as vended_logging_status,
  s.name as source_name,
  d.delivery_destination_type
from
  aws_cloudfront_distribution cf
left join
  aws_cloudwatch_log_delivery_source s
    on s.service = 'cloudfront'
left join
  json_each(s.resource_arns) as ra
    on ra.value = cf.arn
left join
  aws_cloudwatch_log_delivery d
    on s.name = d.delivery_source_name;
```