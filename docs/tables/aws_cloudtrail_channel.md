---
title: "Table: aws_cloudtrail_channel - Query AWS CloudTrail Channel using SQL"
description: "Allows users to query AWS CloudTrail Channel data, including trail configurations, status, and associated metadata."
---

# Table: aws_cloudtrail_channel - Query AWS CloudTrail Channel using SQL

The `aws_cloudtrail_channel` table in Steampipe provides information about CloudTrail trails within AWS CloudTrail. This table allows DevOps engineers to query trail-specific details, including trail configurations, status, and associated metadata. Users can utilize this table to gather insights on trails, such as their status, S3 bucket details, encryption status, and more. The schema outlines the various attributes of the CloudTrail trail, including the trail ARN, home region, S3 bucket name, and whether log file validation is enabled.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_channel` table, you can use the `.inspect aws_cloudtrail_channel` command in Steampipe.

### Key columns:

- `title`: The name of the CloudTrail trail. This is a unique identifier and can be used to join this table with other tables that contain trail-specific information.
- `home_region`: The AWS region in which the trail was created. This can be used to filter results based on region.
- `s3_bucket_name`: The name of the S3 bucket where CloudTrail logs are stored. This can be used to join this table with tables that provide information about S3 buckets.

## Examples

### Basic info

```sql
select
  name,
  arn,
  source,
  apply_to_all_regions
from
  aws_cloudtrail_channel;
```

### List channels that are not applied to all regions

```sql
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

### Get advanced event selector details of each channel

```sql
select
  name,
  a ->> 'Name' as advanced_event_selector_name,
  a ->> 'FieldSelectors' as field_selectors
from
  aws_cloudtrail_channel,
  jsonb_array_elements(advanced_event_selectors) as a;
```

