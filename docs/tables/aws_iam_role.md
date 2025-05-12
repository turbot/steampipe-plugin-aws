---
title: "Steampipe Table: aws_iam_role - Query AWS Identity and Access Management (IAM) Roles using SQL"
description: "Allows users to query IAM Roles to gain insights into their permissions, trust policies, and associated metadata."
folder: "IAM"
---

# Table: aws_iam_role - Query AWS Identity and Access Management (IAM) Roles using SQL

The AWS Identity and Access Management (IAM) Roles are a feature of your AWS account that you can use to delegate permissions to AWS services or users. They enable trusted entities to carry out operations on your behalf, without sharing your root user credentials. Using IAM Roles, one can define a set of permissions to access the resources that a user or service needs to perform tasks.

## Table Usage Guide

The `aws_iam_role` table in Steampipe provides you with information about IAM roles within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query role-specific details, including permissions, trust policies, and associated metadata. You can utilize this table to gather insights on roles, such as roles with wildcard permissions, trust relationships between roles, verification of trust policies, and more. The schema outlines the various attributes of the IAM role for you, including the role ARN, creation date, attached policies, and associated tags.

## Examples

### List IAM roles that have an inline policy.
Determine the areas in which AWS IAM roles have an inline policy configured. This can be useful to assess potential security risks, as inline policies can grant additional permissions to roles.

```sql+postgres
select
  name,
  create_date
from
  aws_iam_role
where
  inline_policies is not null;
```

```sql+sqlite
select
  name,
  create_date
from
  aws_iam_role
where
  inline_policies is not null;
```

### List the attached policies for each IAM role.
Identify the policies associated with each Identity Access Management (IAM) role. This helps in understanding the permissions and access rights granted to each role, aiding in security and access management.


```sql+postgres
select
  name,
  description,
  split_part(policy, '/', 3) as attached_policy
from
  aws_iam_role
  cross join jsonb_array_elements_text(attached_policy_arns) as policy;
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List IAM roles with their associated permission boundaries.
Discover the segments that have IAM roles and the associated permission boundaries. This is useful for understanding the limitations and permissions associated with each role in your AWS environment.

```sql+postgres
select
  name,
  description,
  permissions_boundary_arn,
  permissions_boundary_type
from
  aws_iam_role;
```

```sql+sqlite
select
  name,
  description,
  permissions_boundary_arn,
  permissions_boundary_type
from
  aws_iam_role;
```

### List IAM roles that have policies allowing all (\*) actions.
Identify instances where IAM roles have policies that permit all actions. This can be useful in auditing security settings to ensure that no roles have overly broad permissions, which could pose a security risk.Use this query to identify which AWS IAM roles and their respective policies allow all actions, in order to assess potential security concerns.


```sql+postgres
select
  r.name as role_name,
  p.name as policy_name
from
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as policy_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action
where
  policy_arn = p.arn
  and stmt ->> 'Effect' = 'Allow'
  and action = '*'
order by
  r.name;
```

```sql+sqlite
select
  r.name as role_name,
  p.name as policy_name
from
  aws_iam_role as r,
  aws_iam_policy as p,
  json_each(r.attached_policy_arns) as policy_arn,
  json_each(p.policy_std) as stmt,
  json_each(stmt.value) as action
where
  policy_arn = p.arn
  and json_extract(stmt.value, '$.Effect') = 'Allow'
  and action.value = '*'
order by
  r.name;
```

### Find all iam policy actions with wildcards for a given role.
Discover the segments within your IAM policy that contain wildcard actions for a specific role. This is useful in identifying potential security risks, as these wildcard actions can grant broader permissions than intended.

```sql+postgres
select
  r.name as role_name,
  p.name as policy_name,
  stmt ->> 'Sid' as statement,
  action
from
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as policy_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action
where
  r.name = 'owner'
  and policy_arn = p.arn
  and (
    action like '%*%'
    or action like '%?%'
  );
```

```sql+sqlite
select
  r.name as role_name,
  p.name as policy_name,
  json_extract(stmt.value, '$.Sid') as statement,
  action.value as action
from
  aws_iam_role as r
  join json_each(r.attached_policy_arns) as policy_arn
  join aws_iam_policy as p
  join json_each(p.policy_std) as stmt
  join json_each(json_extract(stmt.value, '$.Action')) as action
where
  r.name = 'owner'
  and policy_arn.value = p.arn
  and (
    action.value like '%*%'
    or action.value like '%?%'
  );
```

### Identify actions that grant elevated privileges to a specific IAM role.
Determine the areas in which specific IAM roles are granted elevated privileges. This is useful for auditing security measures and ensuring that roles do not have unnecessary access rights.

```sql+postgres
select
  r.name,
  a.action,
  a.access_level,
  a.description
from
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action_glob,
  glob(action_glob) as action_regex
  join aws_iam_action as a on a.action like action_regex
where
  pol_arn = p.arn
  and stmt ->> 'Effect' = 'Allow'
  and r.name = 'AWSServiceRoleForRDS'
  and access_level not in ('List', 'Read')
order by
  action;
```

```sql+sqlite
select
  r.name,
  a.action,
  a.access_level,
  a.description
from
  aws_iam_role as r,
  json_each(r.attached_policy_arns) as pol_arn,
  aws_iam_policy as p,
  json_each(json_extract(p.policy_std, '$.Statement')) as stmt,
  json_each(json_extract(stmt.value, '$.Action')) as action_glob,
  aws_iam_action as a
where
  pol_arn.value = p.arn
  and json_extract(stmt.value, '$.Effect') = 'Allow'
  and r.name = 'AWSServiceRoleForRDS'
  and a.action glob action_glob.value
  and a.access_level not in ('List', 'Read')
order by
  a.action;
```

### Compare permission actions between two roles.
Analyze the differences in permission actions between two specific roles in an AWS environment. This is beneficial for understanding discrepancies in access levels, which is essential for maintaining security and proper role allocation.

```sql+postgres
with roles as (
  select
    name,
    attached_policy_arns
  from
    aws_iam_role
  where
    name in ('AWSServiceRoleForSSO', 'AWSServiceRoleForRDS')
),
policies as (
  select
    name,
    arn,
    policy_std
  from
    aws_iam_policy
),
role1_permissions as (
  select
    r.name,
    a.action,
    a.access_level,
    a.description
  from
    roles as r,
    jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
    policies as p,
    jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
    jsonb_array_elements_text(stmt -> 'Action') as action_glob,
    glob (action_glob) as action_regex
    join aws_iam_action a on a.action like action_regex
  where
    pol_arn = p.arn
    and stmt ->> 'Effect' = 'Allow'
    and r.name = 'AWSServiceRoleForSSO'
),
role2_permissions as (
  select
    r.name,
    a.action,
    a.access_level,
    a.description
  from
    roles as r,
    jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
    policies as p,
    jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
    jsonb_array_elements_text(stmt -> 'Action') as action_glob,
    glob (action_glob) as action_regex
    join aws_iam_action a on a.action like action_regex
  where
    pol_arn = p.arn
    and stmt ->> 'Effect' = 'Allow'
    and r.name = 'AWSServiceRoleForRDS'
)
select
  *
from
  role2_permissions
where
  action not in ( select action from role1_permissions)
order by
  action;
```

```sql+sqlite
Error: SQLite does not support JSON operations equivalent to jsonb_array_elements_text and glob function in the given PostgreSQL query.
```

### Identify roles using wildcard principals in their trust policy and those roles trusted by them.
Determine the areas in which roles are using wildcard principals in their trust policy and identify those roles that are trusted by them. This is particularly useful for assessing security vulnerabilities and understanding the levels of access within your AWS environment.

[Refer here](https://twitter.com/nathanwallace/status/1442574375857922048?s=20)

```sql+postgres
select
  maintenance.name,
  admin.name,
  jsonb_pretty(maintenance_stmt),
  jsonb_pretty(admin_stmt)
from
  -- use the account to get the organization_id
  aws_account as a,
  -- check any role as the "maintenance-role"
  aws_iam_role as maintenance,
  -- Combine via join with any role as the "admin-role"
  aws_iam_role as admin,
  jsonb_array_elements(maintenance.assume_role_policy_std -> 'Statement') as maintenance_stmt,
  jsonb_array_elements(admin.assume_role_policy_std -> 'Statement') as admin_stmt
where
  -- maintenance role can be assumed by any AWS principal
  maintenance_stmt -> 'Principal' -> 'AWS' ? '*'
  -- maintenance role principal must be in same account
  and maintenance_stmt -> 'Condition' -> 'StringEquals' -> 'aws:principalorgid' ? a.organization_id
  -- admin role specifically allow maintenance role
  and admin_stmt -> 'Principal' -> 'AWS' ? maintenance.arn;
```

```sql+sqlite
select
  maintenance.name,
  admin.name,
  json_pretty(maintenance_stmt),
  json_pretty(admin_stmt)
from
  aws_account as a,
  aws_iam_role as maintenance,
  aws_iam_role as admin,
  json_each(maintenance.assume_role_policy_std, '$.Statement') as maintenance_stmt,
  json_each(admin.assume_role_policy_std, '$.Statement') as admin_stmt
where
  json_extract(maintenance_stmt.value, '$.Principal.AWS') = '*'
  and json_extract(maintenance_stmt.value, '$.Condition.StringEquals."aws:principalorgid"') = a.organization_id
  and json_extract(admin_stmt.value, '$.Principal.AWS') = maintenance.arn;
```

### List the roles that might allow other roles/users to bypass their assigned IAM permissions.
Determine the areas in which certain roles may potentially allow other users to circumvent their assigned IAM permissions. This query is useful for identifying potential security risks and ensuring that permissions are correctly assigned within your AWS environment.

```sql+postgres
select
  r.name,
  stmt
from
  aws_iam_role as r,
  jsonb_array_elements(r.assume_role_policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Principal' -> 'AWS') as trust
where
  trust = '*'
  or trust like 'arn:aws:iam::%:role/%'
```

```sql+sqlite
select
  r.name,
  json_extract(r.assume_role_policy_std, '$.Statement') as stmt
from
  aws_iam_role as r,
  json_each(r.assume_role_policy_std, '$.Statement'),
  json_each(json_extract(stmt, '$.Principal.AWS')) as trust
where
  trust = '*'
  or trust like 'arn:aws:iam::%:role/%'
```

### Verify the Trust policy of Role has validation conditions when used with GitHub Actions
This query is designed to assess the trust policy of a role when interacting with GitHub Actions. It helps identify any roles that may be missing specific condition checks, providing a valuable tool for enhancing security and ensuring proper configuration.

```sql+postgres
select
  iam.arn as resource,
  iam.description,
  iam.assume_role_policy_std,
  case
    when pstatement -> 'Condition' -> 'StringLike' -> 'token.actions.githubusercontent.com:sub' is not null
    or pstatement -> 'Condition' -> 'StringEquals' -> 'token.actions.githubusercontent.com:sub' is not null then 'ok'
    else 'alarm'
  end as status,
  case
    when pstatement -> 'Condition' -> 'StringLike' -> 'token.actions.githubusercontent.com:sub' is not null
    or pstatement -> 'Condition' -> 'StringEquals' -> 'token.actions.githubusercontent.com:sub' is not null then iam.arn || ' Condition Check Exists'
    else iam.arn || ' Missing Condition Check'
  end as reason
from
  aws_iam_role as iam,
  jsonb_array_elements(iam.assume_role_policy_std -> 'Statement') as pstatement
where
  pstatement -> 'Action' ?& array [ 'sts:assumerolewithwebidentity' ]
  and (pstatement -> 'Principal' -> 'Federated') :: text like '%token.actions.githubusercontent.com%'
order by
  status asc
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```