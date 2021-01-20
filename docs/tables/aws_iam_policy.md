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
  arn not like 'arn:aws:iam::aws:policy%';
```


### List AWS-defined policies

```sql
select
  name,
  arn
from
  aws_iam_policy
where
  arn like 'arn:aws:iam::aws:policy%';
```


### Find unused (unattached) customer-managed policies

```sql
select
  name,
  attachment_count,
  permissions_boundary_usage_count
from
  aws_iam_policy
where
  arn not like 'arn:aws:iam::aws:policy%'
  and attachment_count + permissions_boundary_usage_count = 0;
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
