---
title: "Steampipe Table: aws_iam_policy_attachment - Query AWS IAM Policy Attachments using SQL"
description: "Allows users to query IAM Policy Attachments in AWS to gather information about the relationship between IAM policies and their associated entities (users, groups, and roles)."
folder: "IAM"
---

# Table: aws_iam_policy_attachment - Query AWS IAM Policy Attachments using SQL

The AWS Identity and Access Management (IAM) Policy Attachment is a feature that enables you to attach and detach IAM policies from users, groups, and roles. These policy attachments define what actions are allowed or denied by the attached entities. They are an essential part of managing access permissions in your AWS environment.

## Table Usage Guide

The `aws_iam_policy_attachment` table in Steampipe allows you to query IAM Policy Attachments in AWS to gather information about the relationship between IAM policies and their associated entities (users, groups, and roles). You can use this table to identify which IAM policies are attached to which entities, enabling you to manage and audit access permissions across your AWS environment. The schema outlines the various attributes of the IAM policy attachment, including the policy ARN, policy name, and the associated users, groups, and roles.

**Important Notes**
- Using the `is_attached` column as a filter will help to reduce your query response time.

## Examples

### List attached groups information
Discover the segments that are attached to various policy groups to better manage and organize your AWS IAM policy attachments. This could be beneficial in a real-world scenario where you need to quickly identify and assess the attachments for potential security or configuration issues.

```sql+postgres
select
  policy_arn,
  is_attached,
  policy_groups
from
  aws_iam_policy_attachment
where
  is_attached;
```

```sql+sqlite
select
  policy_arn,
  is_attached,
  policy_groups
from
  aws_iam_policy_attachment
where
  is_attached = 1;
```

### List attached users information
Determine the areas in which user information is attached to policies in your AWS IAM setup. This can be beneficial for auditing and managing user access rights across your AWS environment.

```sql+postgres
select
  policy_arn,
  is_attached,
  policy_users
from
  aws_iam_policy_attachment
where
  is_attached;
```

```sql+sqlite
select
  policy_arn,
  is_attached,
  policy_users
from
  aws_iam_policy_attachment
where
  is_attached = 1;
```

### List users with AdministratorAccess policy
Identify instances where users have been granted 'AdministratorAccess' within your AWS IAM policies. This is useful for auditing security and managing access control across your AWS environment.

```sql+postgres
select
  name as policy_name, 
  policy_arn, 
  jsonb_pretty(policy_users) as policy_users
from
  aws_iam_policy p
  left join aws_iam_policy_attachment a on p.arn = a.policy_arn 
where
  name = 'AdministratorAccess' and a.is_attached;
```

```sql+sqlite
select
  name as policy_name, 
  policy_arn, 
  json_pretty(policy_users) as policy_users
from
  aws_iam_policy p
  left join aws_iam_policy_attachment a on p.arn = a.policy_arn 
where
  name = 'AdministratorAccess' and a.is_attached = 1;
```