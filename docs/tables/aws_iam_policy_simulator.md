# Table: aws_iam_policy_simulator

The IAM policy simulator allows you to test and troubleshoot IAM policies.

## Examples

### Does user bob have s3:DeleteBucket on any resource?
```sql
select
  decision
from
  aws_iam_policy_simulator
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```

### Which users have have s3:DeleteBucket on any resource?
```sql
select
  u.name,
  decision
from
  aws_iam_policy_simulator p,
  aws_iam_user u
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and p.principal_arn = u.arn;
```
