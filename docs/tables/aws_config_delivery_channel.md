---
title: "Steampipe Table: aws_config_delivery_channel - Query AWS Config Delivery Channels using SQL"
description: "Allows users to query AWS Config Delivery Channels"
folder: "Config"
---

# Table: aws_config_delivery_channel - Query AWS Config Delivery Channels using SQL

The AWS Config Delivery Channel is a feature that enables AWS Config to deliver configuration snapshots and configuration change notifications to specified destinations. It plays a key role in ensuring that your configuration data is stored securely and notifications are sent promptly for compliance or operational purposes.

## Table Usage Guide

The `aws_config_delivery_channel` table in Steampipe provides insights into the Delivery Channels associated with AWS Config. This table enables DevOps engineers, security analysts, and cloud administrators to query delivery channel details such as the destination S3 bucket, SNS topic for notifications, and delivery status. Use this table to ensure your configuration change data is being delivered correctly and troubleshoot delivery-related issues.

## Examples

### Retrieve basic delivery channel information
Get a detailed view of your AWS Config Delivery Channels, including their destinations and notification settings.

```sql+postgres
select
  name,
  s3_bucket_name,
  s3_key_prefix,
  sns_topic_arn,
  delivery_frequency,
  status,
  title,
  akas
from
  aws_config_delivery_channel;
```

```sql+sqlite
select
  name,
  s3_bucket_name,
  s3_key_prefix,
  sns_topic_arn,
  delivery_frequency,
  status,
  title,
  akas
from
  aws_config_delivery_channel;
```

### List delivery channels without SNS topic configured
Identify delivery channels that do not have an SNS topic configured for notifications. This can help ensure you have proper alerting mechanisms in place.

```sql+postgres
select
  name,
  s3_bucket_name,
  sns_topic_arn
from
  aws_config_delivery_channel
where
  sns_topic_arn is null;
```

```sql+sqlite
select
  name,
  s3_bucket_name,
  sns_topic_arn
from
  aws_config_delivery_channel
where
  sns_topic_arn is null;
```

### Check delivery channels with delivery failures
Discover delivery channels with failed deliveries to address issues in your AWS Config setup.

```sql+postgres
select
  name,
  status ->> 'LastStatus' as last_status,
  status ->> 'LastStatusChangeTime' as last_status_change_time,
  status ->> 'LastErrorCode' as last_error_code,
  status ->> 'LastErrorMessage' as last_error_message
from
  aws_config_delivery_channel
where
  (status ->> 'LastStatus') = 'FAILURE';
```

```sql+sqlite
select
  name,
  json_extract(status, '$.LastStatus') as last_status,
  json_extract(status, '$.LastStatusChangeTime') as last_status_change_time,
  json_extract(status, '$.LastErrorCode') as last_error_code,
  json_extract(status, '$.LastErrorMessage') as last_error_message
from
  aws_config_delivery_channel
where
  json_extract(status, '$.LastStatus') = 'FAILURE';
```

### List delivery channels sending to a specific S3 bucket
Query the delivery channels that are configured to send data to a particular S3 bucket.

```sql+postgres
select
  name,
  s3_bucket_name,
  sns_topic_arn,
  delivery_frequency
from
  aws_config_delivery_channel
where
  s3_bucket_name = 'test-bucket-delivery-channel';
```

```sql+sqlite
select
  name,
  s3_bucket_name,
  sns_topic_arn,
  delivery_frequency
from
  aws_config_delivery_channel
where
  s3_bucket_name = 'test-bucket-delivery-channel';
```

### Analyze delivery frequency of all channels
Get an overview of how often your delivery channels send data, ensuring they align with organizational requirements.

```sql+postgres
select
  name,
  delivery_frequency,
  s3_bucket_name,
  sns_topic_arn
from
  aws_config_delivery_channel;
```

```sql+sqlite
select
  name,
  delivery_frequency,
  s3_bucket_name,
  sns_topic_arn
from
  aws_config_delivery_channel;
```
