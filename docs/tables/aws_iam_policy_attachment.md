# Table: aws_iam_policy_attachment

This table provide information about the attached user(s), role(s), and group(s) to a Managed IAM Policy

Note that the `is_attached` column will help to provide information for only the attached ones. You can greatly decrease your query time by using **where is_attached**.

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