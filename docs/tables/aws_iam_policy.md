---
title: "Table: aws_iam_policy - Query AWS IAM Policy using SQL"
description: "Allows users to query AWS IAM Policies, providing detailed information about each policy, including permissions, attachment, and associated metadata."
---

# Table: aws_iam_policy - Query AWS IAM Policy using SQL

The `aws_iam_policy` table in Steampipe provides information about IAM policies within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query policy-specific details, including permissions, attachments, and associated metadata. Users can utilize this table to gather insights on policies, such as policies with wildcard permissions, verification of policy documents, and more. The schema outlines the various attributes of the IAM policy, including the policy ARN, creation date, update date, attached entities, and policy default version ID.

**Note** that the `policy` and `policy_std` columns require additional calls - You can greatly decrease your query time by NOT selecting those columns when you don't need them.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_policy` table, you can use the `.inspect aws_iam_policy` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the policy. This is a unique identifier that can be used to join this table with other tables.
- `policy_name`: The name of the policy. This can be used to filter the policies based on specific naming conventions.
- `default_version_id`: The ID of the policy's default version. This can be useful when joining with a table of policy versions to get the details of the default policy version.

## Examples

### List customer-defined policies

```sql
select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed;
```

### List customer-defined policies with a path prefix

```sql
select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed
  and path = '/turbot/';
```

### Find attached customer-managed policies

```sql
select
  name,
  arn,
  permissions_boundary_usage_count
from
  aws_iam_policy
where
  is_attached;
```

### Find unused customer-managed policies

```sql
select
  name,
  attachment_count,
  permissions_boundary_usage_count
from
  aws_iam_policy
where
  not is_aws_managed
  and not is_attached
  and permissions_boundary_usage_count = 0;
```

### Find policy statements that grant Full Control (*:*) access

```sql
select
  name,
  arn,
  action,
  s ->> 'Effect' as effect
from
  aws_iam_policy,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Action') as action
where
  action in ('*', '*:*')
  and s ->> 'Effect' = 'Allow';
```

### Find policy statements that grant service level full access

```sql
select
  name,
  arn,
  action,
  s ->> 'Effect' as effect
from
  aws_iam_policy,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Action') as action
where
  s ->> 'Effect' = 'Allow'
  and (
    action = '*'
    or action like '%:*'
  );
```

### Expand wildcards to list all actions granted by a policy

```sql
select
  a.action,
  a.access_level,
  a.description
from
  aws_iam_policy p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action_glob,
  glob(action_glob) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  p.name = 'AmazonEC2ReadOnlyAccess'
  and stmt ->> 'Effect' = 'Allow'
order by
  a.action;
```
