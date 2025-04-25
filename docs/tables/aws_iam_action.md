---
title: "Steampipe Table: aws_iam_action - Query AWS IAM Action using SQL"
description: "Allows users to query IAM Actions in AWS Identity and Access Management (IAM)."
folder: "IAM"
---

# Table: aws_iam_action - Query AWS IAM Action using SQL

The AWS IAM Action is a component of the AWS Identity and Access Management (IAM) service. It allows you to securely control access to AWS services and resources for your users. You can use IAM actions to allow or deny permissions to AWS resources, based on SQL queries, ensuring the right individuals have the appropriate access.

## Table Usage Guide

The `aws_iam_action` table in Steampipe provides you with information about IAM actions within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query action-specific details, including the action name, description, resource types, and condition keys. You can utilize this table to gather insights on actions, such as actions allowed for a specific resource type, actions that support specific condition keys, and more. The schema outlines the various attributes of the IAM action, including the action name, description, resource types, condition keys, and associated metadata.

**Important Notes**
- You can access the list of possible IAM actions in AWS, along with their access levels and descriptions. The data is sourced from [Parliament](https://github.com/duo-labs/parliament).

- When you use the `aws_iam_action` to search for actions in other tables:
  - You might want to use the `policy_std` column instead of `policy`, as the format is standardized including converting action names to lower case.
  - You might want to join on the `action` column in the `aws_iam_action` as it is also converted to lowercase.

## Examples

### List all actions associated with the s3 service
Explore which actions are linked to a specific cloud storage service to gain insights into service-specific permissions and operation capabilities. This can be particularly useful for managing access controls and understanding the scope of service functionalities.

```sql+postgres
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

```sql+sqlite
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
Gain insights into the specific functionality of the 's3:deleteobject' action in AWS IAM. This is useful for understanding the implications of using this action in your AWS environment.

```sql+postgres
select
  description
from
  aws_iam_action
where
  action = 's3:deleteobject';
```

```sql+sqlite
select
  description
from
  aws_iam_action
where
  action = 's3:deleteobject';
```


### List the actions that are included in 's3:d*'
Explore which actions are included within a specific pattern to gain insights into your AWS IAM configuration. This can help in assessing the elements within your security settings and pinpointing specific areas that match the pattern for better management and security compliance.

```sql+postgres
select
  a.action,
  a.description
from
  aws_iam_action as a,
  glob('s3:d*') as action_name
where
  a.action like action_name;
```

```sql+sqlite
select
  a.action,
  a.description
from
  aws_iam_action as a
where
  a.action like 's3:d%';
```

### Get the list of expanded actions granted in a policy
Determine the areas in which specific policy permissions are granted. This is particularly useful when you want to understand the scope of access that's been allowed under a given policy, such as 'AmazonEC2ReadOnlyAccess'.

```sql+postgres
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

```sql+sqlite
select
  a.action,
  a.access_level
from
  aws_iam_policy p,
  json_each(p.policy_std, '$.Statement') as stmt,
  json_each(stmt.value, '$.Action') as action_glob,
  glob(action_glob.value) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  p.name = 'AmazonEC2ReadOnlyAccess'
  and json_extract(stmt.value, '$.Effect') = 'Allow'
order by
  a.action;
```


### List all the actions allowed by managed policies for a Lambda execution role
Discover the permissions granted by managed policies for a specific Lambda execution role. This query is useful for auditing security configurations, ensuring only necessary permissions are allowed.

```sql+postgres
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

```sql+sqlite
select 
  f.name, 
  f.role, 
  a.action, 
  a.access_level, 
  a.description 
from 
  aws_lambda_function as f 
  join aws_iam_role as r on f.role = r.arn 
  join aws_iam_policy as p, 
  json_each(r.attached_policy_arns) as pol_arn 
where 
  pol_arn.value = p.arn 
  and json_extract(p.policy_std, '$.Statement') = 'Allow' 
  and f.name = 'hellopython' 
  and a.action in (
    select 
      value 
    from 
      json_each(p.policy_std, '$.Statement.Action')
  );
```