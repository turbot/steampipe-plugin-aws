---
title: "Steampipe Table: aws_shield_drt_access - Query AWS Shield Advanced SRT Access settings using SQL"
description: "Allows users to query AWS Shield Advanced SRT Access settings and retrieve information about the Shield Response Team's access to your AWS resources."
---

# Table: aws_shield_drt_access - Query AWS Shield Advanced SRT Access settings using SQL

AWS Shield Advanced is a DDoS protection service from AWS. The SRT Access settings allow you to configure the access permissions for the AWS Shield Advanced Shield Response Team (SRT) to the resources in your account.

## Table Usage Guide

The `aws_shield_drt_access` table in Steampipe allows you to query the AWS Shield Advanced SRT Access settings and retrieve information about the IAM role and S3 Buckets the SRT should have access to. For more details about the individual fields, please refer to the [AWS Shield Advanced API documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeDRTAccess.html).

## Examples

### Basic info

```sql+postgres
select
  role_arn,
  jsonb_pretty(log_bucket_list) as log_bucket_list
from
  aws_shield_protection;
```

```sql+sqlite
select
  role_arn,
  json_pretty(log_bucket_list) as log_bucket_list
from
  aws_shield_protection;
```
