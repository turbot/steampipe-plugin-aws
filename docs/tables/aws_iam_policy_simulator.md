# Table: aws_iam_policy_simulator

The IAM policy simulator allows you to test and troubleshoot IAM policies.

Note that you ***must*** specify a single `action`, `resource_arn`, and `principal_arn` in a where or join clause in order to use this table.  


## Examples

### Check if user has s3:DeleteBucket on any resource
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


### Check if user has 'ec2:terminateinstances' on any resource including details of any policy granting or denying access

```sql
select
  decision,
  jsonb_pretty(matched_statements)
from
  aws_iam_policy_simulator
where
  action = 'ec2:terminateinstances'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```

### For all users in the account, check whether they have `sts:AssumeRole` on all roles.
```sql
select
  u.name,
  decision
from
  aws_iam_policy_simulator p,
  aws_iam_user u
where
  action = 'sts:AssumeRole'
  and resource_arn = '*'
  and p.principal_arn = u.arn;
```
