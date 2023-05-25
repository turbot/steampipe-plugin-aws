# Table: aws_iam_role

An IAM role is an AWS Identity and Access Management (IAM) entity with permissions to make AWS service requests.

## Examples

### List of IAM roles with no inline policy

```sql
select
  name,
  create_date
from
  aws_iam_role
where
  inline_policies is null;
```


### List the policies attached to the roles

```sql
select
  name,
  description,
  split_part(policy, '/', 3) as attached_policy
from
  aws_iam_role
  cross join jsonb_array_elements_text(attached_policy_arns) as policy;
```


### Permission boundary information for each role

```sql
select
  name,
  description,
  permissions_boundary_arn,
  permissions_boundary_type
from
  aws_iam_role;
```

### Find all roles that allow *
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

### Find any roles that allow wildcard actions
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

### List higher-level permissions for any specific role
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

### List all actions (with level) in role2, not in role1
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

### List role with wildcard principal in trust policy(maintenance-role) and role(admin-role) that have trust relationship with maintenance-role
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
  allaws.aws_iam_role as iam,
  jsonb_array_elements(iam.assume_role_policy_std -> 'Statement') as pstatement
where
  pstatement -> 'Action' ?& array [ 'sts:assumerolewithwebidentity' ]
  and (pstatement -> 'Principal' -> 'Federated') :: text like '%token.actions.githubusercontent.com%'
order by
  status asc
```
