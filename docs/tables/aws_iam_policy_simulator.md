# Table: aws_iam_policy_simulator

The IAM policy simulator allows you to test and troubleshoot IAM policies.

Note that you ***must*** specify a single `action`, `resource_arn`, and `principal_arn` in a where clause in order to use this table.  Also, see the note below on issue relating to a [known issue](https://github.com/turbot/steampipe-postgres-fdw/issues/3) with nested select queries (select where in (select ...)) and joins on tables with required key columns.

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


## NOTE: Issue with nested select queries and joins on tables with required key columns
Currently, there is a [known issue](https://github.com/turbot/steampipe-postgres-fdw/issues/3) with nested select queries (select where in (select ...)) and joins on tables with required key columns. It seems that the qualifiers are not passed to the parent query because the nested query is executed in parallel. We are actively working to resolve this issue.

For example, this works as you would expect:

```sql
select
  *
from
  aws_iam_policy_simulator
where
  principal_arn = 'arn:aws:iam::012345678901:user/mike'
  and resource_arn = '*'
  and action = 's3:DeleteBucket'
```


This SHOULD work but currently doesn't:
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
```
Error: pq: rpc error: code = Internal desc = 'List' call requires an '=' qual for all columns: principal_arn,action,resource_arn
```

This SHOULD ALSO work but currently doesn't:
```sql
select
  decision
from
  aws_iam_policy_simulator
where
  action = 's3:DeleteBucket'
  and resource_arn = '*'
  and principal_arn = (select  name from  aws_iam_user );
```
```
Error: pq: rpc error: code = Internal desc = 'List' call requires an '=' qual for all columns: principal_arn,action,resource_arn
```
