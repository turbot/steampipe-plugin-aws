---
title: "Table: aws_iam_role - Query AWS Identity and Access Management (IAM) Roles using SQL"
description: "Allows users to query IAM Roles to gain insights into their permissions, trust policies, and associated metadata."
---

# Table: aws_iam_role - Query AWS Identity and Access Management (IAM) Roles using SQL

The `aws_iam_role` table in Steampipe provides information about IAM roles within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query role-specific details, including permissions, trust policies, and associated metadata. Users can utilize this table to gather insights on roles, such as roles with wildcard permissions, trust relationships between roles, verification of trust policies, and more. The schema outlines the various attributes of the IAM role, including the role ARN, creation date, attached policies, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_role` table, you can use the `.inspect aws_iam_role` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the role. This can be used to join with other tables that also contain IAM role ARNs.

- `role_name`: The name of the role. This can be used to join with other tables that also contain IAM role names.

- `create_date`: The date and time, in ISO 8601 date-time format, when the role was created. This can be useful in determining the age of the role.

## Examples

### List IAM roles that have an inline policy.
Use this query to identify AWS IAM roles that have inline policies to help administrators identify roles that might be using inline policies instead of the recommended managed policies.

```sql
select
  name,
  create_date
from
  aws_iam_role
where
  inline_policies is not null;
```

### List the attached policies for each IAM role.
Use this query to determine which policies are attached to each AWS IAM role and highlights the importance of understanding the permissions granted to those roles.

```sql
select
  name,
  description,
  split_part(policy, '/', 3) as attached_policy
from
  aws_iam_role
  cross join jsonb_array_elements_text(attached_policy_arns) as policy;
```

### List IAM roles with their associated permission boundaries.
Use this query to list AWS IAM roles with their descriptions and associated permissions boundaries to better manage and understand role permissions.

```sql
select
  name,
  description,
  permissions_boundary_arn,
  permissions_boundary_type
from
  aws_iam_role;
```

### List IAM roles that have policies allowing all (\*) actions.
Use this query to identify which AWS IAM roles and their respective policies allow all actions, in order to assess potential security concerns.

```sql
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

### Find all iam policy actions with wildcards for a given role.
Use this query to identify AWS IAM policy actions that use wildcard characters for any role to ensure policy configurations are not overly permissive.

```sql
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

### Identify actions that grant elevated privileges to a specific IAM role.
Use this query to identify which actions permit an IAM role (e.g., "AWSServiceRoleForRDS") in AWS to execute tasks beyond basic list and read functions, aiding in recognizing and addressing potential security concerns.

```sql
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

### Compare permission actions between two roles.
Use this query to identify which permission actions and their associated access levels are unique to a specified IAM role, to understand differences in permissions when compared to another specific IAM role.

```sql
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

### Identify roles using wildcard principals in their trust policy and those roles trusted by them.
Use this query to locate AWS IAM roles with an open trust policy within the same organization and identify other roles that trust them.

[Refer here](https://twitter.com/nathanwallace/status/1442574375857922048?s=20)

```sql
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

### List the roles that might allow other roles/users to bypass their assigned IAM permissions.
Use this query to determine which AWS IAM roles can be potentially assumed by any user or another role, highlighting potential security concerns for unauthorized access or privilege escalation

```sql
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

### Verify the Trust policy of Role has validation conditions when used with GitHub Actions
Use this query to evaluate AWS IAM roles and ascertain if they include validation conditions when invoked via GitHub Actions. Specifically, it checks for the presence of conditions related to the token.actions.githubusercontent.com domain within the trust policy of the role. If such conditions exist, it will label the role as 'ok'; otherwise, it will be labeled as 'alarm'. Additionally, the query provides a reason for the assigned status based on whether the condition check exists or is missing.

```sql
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
