# Table: aws_iam_policy_simulator

The IAM policy simulator allows you to test and troubleshoot IAM policies.

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



### Check if user has s3:DeleteBucket on any resource including details of any policy granting or denying access

```sql
select
  decision,
  jsonb_pretty(matched_statements)
from
  aws_iam_policy_simulator
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and principal_arn = 'arn:aws:iam::012345678901:user/bob';
```
