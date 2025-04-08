---
title: "Steampipe Table: aws_cloudtrail_channel - Query AWS CloudTrail Channel using SQL"
description: "Allows users to query AWS CloudTrail Channel data, including trail configurations, status, and associated metadata."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_channel - Query AWS CloudTrail Channel using SQL

The AWS CloudTrail is a service that enables governance, compliance, operational auditing, and risk auditing of your AWS account. It helps you to log, continuously monitor, and retain account activity related to actions across your AWS infrastructure. The CloudTrail Channel specifically, allows you to manage the delivery of CloudTrail event log files to your specified S3 bucket and CloudWatch Logs log group.

## Table Usage Guide

The `aws_cloudtrail_channel` table in Steampipe provides you with information about CloudTrail trails within AWS CloudTrail. This table allows you, as a DevOps engineer, to query trail-specific details, including trail configurations, status, and associated metadata. You can utilize this table to gather insights on trails, such as their status, S3 bucket details, encryption status, and more. The schema outlines the various attributes of the CloudTrail trail for you, including the trail ARN, home region, S3 bucket name, and whether log file validation is enabled.

## Examples

### Basic info
Analyze the settings of your AWS CloudTrail channels to understand whether they are applied to all regions. This is beneficial to ensure consistent logging and monitoring across your entire AWS environment.

```sql+postgres
select
  name,
  arn,
  source,
  apply_to_all_regions
from
  aws_cloudtrail_channel;
```

```sql+sqlite
select
  name,
  arn,
  source,
  apply_to_all_regions
from
  aws_cloudtrail_channel;
```

### List channels that are not applied to all regions
Identify the AWS Cloudtrail channels which are not configured to apply to all regions. This can be useful for auditing regional compliance or identifying potential gaps in log coverage.

```sql+postgres
select
  name,
  arn,
  source,
  apply_to_all_regions,
  advanced_event_selectors
from
  aws_cloudtrail_channel
where
  not apply_to_all_regions;
```

```sql+sqlite
select
  name,
  arn,
  source,
  apply_to_all_regions,
  advanced_event_selectors
from
  aws_cloudtrail_channel
where
  apply_to_all_regions = 0;
```

### Get advanced event selector details of each channel
Determine the specific event selector details associated with each AWS CloudTrail channel. This query is useful for analyzing channel configurations and identifying any potential areas for optimization or troubleshooting.

```sql+postgres
select
  name,
  a ->> 'Name' as advanced_event_selector_name,
  a ->> 'FieldSelectors' as field_selectors
from
  aws_cloudtrail_channel,
  jsonb_array_elements(advanced_event_selectors) as a;
```

```sql+sqlite
select
  name,
  json_extract(a.value, '$.Name') as advanced_event_selector_name,
  json_extract(a.value, '$.FieldSelectors') as field_selectors
from
  aws_cloudtrail_channel,
  json_each(advanced_event_selectors) as a;
```