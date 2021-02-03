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


