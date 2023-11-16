---
title: "Table: aws_iam_action - Query AWS IAM Action using SQL"
description: "Allows users to query IAM Actions in AWS Identity and Access Management (IAM)."
---

# Table: aws_iam_action - Query AWS IAM Action using SQL

The `aws_iam_action` table in Steampipe provides information about IAM actions within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query action-specific details, including the action name, description, resource types, and condition keys. Users can utilize this table to gather insights on actions, such as actions allowed for a specific resource type, actions that support specific condition keys, and more. The schema outlines the various attributes of the IAM action, including the action name, description, resource types, condition keys, and associated metadata.

The list of possible IAM actions in AWS, along with their access levels and descriptions. The data is sourced from [Parliament](https://github.com/duo-labs/parliament).

When using the `aws_iam_action` to search for actions in other tables:
- You probably want to use the `policy_std` column instead of `policy`, as the format is standardized including converting action names to lower case.
- You probably want to join on the `action` column in the `aws_iam_action` as it is also converted to lowercase.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_action` table, you can use the `.inspect aws_iam_action` command in Steampipe.

Key columns:

- `name`: The name of the action. This is the primary key of the table and can be used to join this table with other tables.
- `resource_types`: The types of resources that the action can apply to. This can be useful when querying for actions that are applicable to specific resource types.
- `condition_keys`: The condition keys that the action supports. This can be useful when querying for actions that support specific condition keys.

## Examples

### List all actions associated with the s3 service
```sql
select
  action,
  description
from
  aws_iam_action
where
  prefix = 's3'
order by
  action;
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


### List the actions that are included in 's3:d*'
```sql
select
  a.action,
  a.description
from
  aws_iam_action as a,
  glob('s3:d*') as action_name
where
  a.action like action_name;
```

### Get the list of expanded actions granted in a policy
```sql
select
  a.action,
  a.access_level
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


### List all the actions allowed by managed policies for a Lambda execution role
```sql
select
  f.name,
  f.role,
  a.action,
  a.access_level,
  a.description
from 
  aws_lambda_function as f,
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action_glob,
  glob(action_glob) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  f.role = r.arn
  and pol_arn = p.arn 
  and stmt ->> 'Effect' = 'Allow'
  and f.name = 'hellopython';
```