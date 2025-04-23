---
title: "Steampipe Table: aws_shield_drt_access - Query AWS Shield Advanced SRT Access settings using SQL"
description: "Allows users to query AWS Shield Advanced SRT Access settings and retrieve information about the Shield Response Team's access to your AWS resources."
folder: "Shield"
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
  log_bucket_list
from
  aws_shield_drt_access;
```

```sql+sqlite
select
  role_arn,
  log_bucket_list
from
  aws_shield_drt_access;
```

### Check if the SRT role has the correct permissions

```sql+postgres
select
  role.arn,
  role.name,
  trust_policy_statement -> 'Principal' -> 'Service' ? 'drt.shield.amazonaws.com' as can_be_assumed_by_shield,
  role.attached_policy_arns ? 'arn:aws:iam::aws:policy/service-role/AWSShieldDRTAccessPolicy' as has_shield_drt_access_policy
from
  aws_shield_drt_access
join
  aws_iam_role as role on role.arn = aws_shield_drt_access.role_arn,
  jsonb_array_elements(role.assume_role_policy_std -> 'Statement') as trust_policy_statement;
```

```sql+sqlite
select
  role.arn,
  role.name,
  trust_policy_statement -> 'Principal' -> 'Service' ? 'drt.shield.amazonaws.com' as can_be_assumed_by_shield,
  role.attached_policy_arns ? 'arn:aws:iam::aws:policy/service-role/AWSShieldDRTAccessPolicy' as has_shield_drt_access_policy
from
  aws_shield_drt_access
join
  aws_iam_role as role on role.arn = aws_shield_drt_access.role_arn,
  json_each(role.assume_role_policy_std -> 'Statement') as trust_policy_statement;
```
