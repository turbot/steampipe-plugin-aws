# Table: aws_iam_policy

An IAM Policy is an AWS Identity and Access Management (IAM) Managed Policy

Note that the `policy` and `policy_std` columns require additional calls - You can greatly decrease your query time by NOT selecting those columns when you don't need them.

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

### Find policy statements that grant Full Control (_:_) access

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
