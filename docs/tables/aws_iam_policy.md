---
title: "Steampipe Table: aws_iam_policy - Query AWS IAM Policy using SQL"
description: "Allows users to query AWS IAM Policies, providing detailed information about each policy, including permissions, attachment, and associated metadata."
folder: "IAM"
---

# Table: aws_iam_policy - Query AWS IAM Policy using SQL

The AWS Identity and Access Management (IAM) Policy is a resource that allows you to manage permissions and control access to AWS services and resources. With IAM policies, you can specify who is allowed and denied access, and what actions they can or cannot perform. These policies help you secure your AWS resources, ensure compliance with your security policies, and manage access across your entire AWS environment.

## Table Usage Guide

The `aws_iam_policy` table in Steampipe provides you with information about IAM policies within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query policy-specific details, including permissions, attachments, and associated metadata. You can utilize this table to gather insights on policies, such as policies with wildcard permissions, verification of policy documents, and more. The schema outlines the various attributes of the IAM policy for you, including the policy ARN, creation date, update date, attached entities, and policy default version ID.

**Important Notes**
- The `policy` and `policy_std` columns require additional calls - You can greatly decrease your query time by NOT selecting those columns when you don't need them.

## Examples

### List customer-defined policies
Determine the areas in which custom policies, as defined by the user, are implemented within the AWS IAM service. This query is useful for auditing security measures and ensuring that AWS resources are governed by appropriate, user-defined policies.

```sql+postgres
select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed;
```

```sql+sqlite
select
  name,
  arn
from
  aws_iam_policy
where
  is_aws_managed = 0;
```

### List customer-defined policies with a path prefix
Explore the custom policies within a specific path prefix to understand their names and resources, which is particularly beneficial for managing and organizing security controls in a streamlined manner. This allows for efficient monitoring and modification of policies that are not managed by AWS, hence offering enhanced control over your security infrastructure.

```sql+postgres
select
  name,
  arn
from
  aws_iam_policy
where
  not is_aws_managed
  and path = '/turbot/';
```

```sql+sqlite
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
Discover the segments that are utilizing customer-managed policies within your AWS environment. This allows you to better manage your resources and understand which policies are attached, enhancing overall security and governance.

```sql+postgres
select
  name,
  arn,
  permissions_boundary_usage_count
from
  aws_iam_policy
where
  is_attached;
```

```sql+sqlite
select
  name,
  arn,
  permissions_boundary_usage_count
from
  aws_iam_policy
where
  is_attached = 1;
```

### Find unused customer-managed policies
Determine the areas in which customer-managed policies are not being utilized. This is beneficial in identifying potential areas of cost reduction and improving security by eliminating unnecessary permissions.

```sql+postgres
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

```sql+sqlite
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
This example helps identify policies that potentially grant unrestricted access, allowing for a comprehensive review of security settings. It aids in enhancing security by pinpointing areas where permissions may be overly broad.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### Find policy statements that grant service level full access
Explore which policy statements allow full service level access. This can be useful for maintaining security standards by identifying policies that may potentially grant excessive permissions.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  json_extract(s.value, '$.Action') as action,
  json_extract(s.value, '$.Effect') as effect
from
  aws_iam_policy,
  json_each(policy_std, 'Statement') as s
where
  json_extract(s.value, '$.Effect') = 'Allow'
  and (
    json_extract(s.value, '$.Action') = '*'
    or json_extract(s.value, '$.Action') like '%:*'
  );
```

### Expand wildcards to list all actions granted by a policy
Identify all actions permitted by a specific policy. This is particularly useful for understanding the scope of permissions given to a particular policy, thereby aiding in effective access management.

```sql+postgres
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

```sql+sqlite
select
  a.action,
  a.access_level,
  a.description
from
  aws_iam_policy p,
  json_each(p.policy_std, '$.Statement') as stmt,
  json_each(stmt.value, '$.Action') as action_glob,
  glob(action_glob.value) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  p.name = 'AmazonEC2ReadOnlyAccess'
  and json_extract(stmt.value, '$.Effect') = 'Allow'
order by
  a.action;
```

