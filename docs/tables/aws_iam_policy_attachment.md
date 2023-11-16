---
title: "Table: aws_iam_policy_attachment - Query AWS IAM Policy Attachments using SQL"
description: "Allows users to query IAM Policy Attachments in AWS to gather information about the relationship between IAM policies and their associated entities (users, groups, and roles)."
---

# Table: aws_iam_policy_attachment - Query AWS IAM Policy Attachments using SQL

The `aws_iam_policy_attachment` table in Steampipe allows users to query IAM Policy Attachments in AWS to gather information about the relationship between IAM policies and their associated entities (users, groups, and roles). This table can be used to identify which IAM policies are attached to which entities, enabling users to manage and audit access permissions across their AWS environment. The schema outlines the various attributes of the IAM policy attachment, including the policy ARN, policy name, and the associated users, groups, and roles.

**Note:** Using `is_attached` column as filter, will help to reduce the query response time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_policy_attachment` table, you can use the `.inspect aws_iam_policy_attachment` command in Steampipe.

**Key columns**:

- `policy_arn`: The Amazon Resource Name (ARN) of the IAM policy. This can be used to join with the `aws_iam_policy` table to get more details about the policy.
- `policy_name`: The name of the IAM policy. This provides a human-readable way to identify policies.
- `users`, `groups`, `roles`: These columns list the associated entities for each policy. These can be used to join with the `aws_iam_user`, `aws_iam_group`, and `aws_iam_role` tables respectively to get more details about the entities.

## Examples

### List attached groups information

```sql
select
  policy_arn,
  is_attached,
  policy_groups
from
  aws_iam_policy_attachment
where
  is_attached;
```

### List attached users information

```sql
select
  policy_arn,
  is_attached,
  policy_users
from
  aws_iam_policy_attachment
where
  is_attached;
```

### List users with AdministratorAccess policy

```sql
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
