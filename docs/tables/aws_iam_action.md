# Table: aws_iam_action

The list of possible IAM actions in AWS, along with their access levels and descriptions. The data is sourced from [Parliament](https://github.com/duo-labs/parliament).

## Examples

### List all actions associated with the s3 service
```sql
select
    action
from
    aws_iam_action
where
    prefix = 's3'
order by action;
```

### Get a description for the s3:deleteobject action
```sql
select
    description
from
    aws_iam_action
where
    action = 's3:deleteobject';
```

### Get the list of expanded actions for a role
```sql
select a.action, a.access_level
from aws_iam_policy p,
     jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
     jsonb_array_elements_text(stmt -> 'Action') as action_glob,
     glob(action_glob) as action_regex
         join aws_iam_action a ON a.action LIKE action_regex
where p.name = 'owner'
order by a.action;
```
