---
title: "Steampipe Table: aws_iam_group - Query AWS IAM Group using SQL"
description: "Allows users to query AWS IAM Group data such as group name, path, and ARN. This table provides information about IAM groups within AWS Identity and Access Management (IAM)."
folder: "IAM"
---

# Table: aws_iam_group - Query AWS IAM Group using SQL

The AWS Identity and Access Management (IAM) Group is a feature that allows you to manage user access to AWS services and resources. With IAM Groups, you can specify permissions for multiple users, which can make it easier to manage the permissions for those users. IAM Groups are not truly identities because they cannot be identified as Principals in a resource's policy.

## Table Usage Guide

The `aws_iam_group` table in Steampipe provides you with information about IAM groups within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query group-specific details, including group name, path, and ARN. You can utilize this table to gather insights on groups, such as group membership, group policy attachments, and more. The schema outlines the various attributes of the IAM group for you, including the group ARN, creation date, group ID, and associated metadata.

## Examples

### User details associated with each IAM group
Explore which users are associated with each IAM group, including their user details such as permissions boundary, last password usage and creation date. This can be useful for auditing user access and ensuring appropriate permissions are in place.

```sql+postgres
select
  name as group_name,
  iam_user ->> 'UserName' as user_name,
  iam_user ->> 'UserId' as user_id,
  iam_user ->> 'PermissionsBoundary' as permission_boundary,
  iam_user ->> 'PasswordLastUsed' as password_last_used,
  iam_user ->> 'CreateDate' as user_create_date
from
  aws_iam_group
  cross join jsonb_array_elements(users) as iam_user;
```

```sql+sqlite
select
  name as group_name,
  json_extract(iam_user, '$.UserName') as user_name,
  json_extract(iam_user, '$.UserId') as user_id,
  json_extract(iam_user, '$.PermissionsBoundary') as permission_boundary,
  json_extract(iam_user, '$.PasswordLastUsed') as password_last_used,
  json_extract(iam_user, '$.CreateDate') as user_create_date
from
  aws_iam_group,
  json_each(users) as iam_user;
```


### List all the users in each group having Administrator access
Discover the segments that include users with Administrator access across different groups. This is beneficial for auditing purposes, allowing for a quick overview of who has high-level access and potential control within your system.

```sql+postgres
select
  name as group_name,
  iam_user ->> 'UserName' as user_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_group
  cross join jsonb_array_elements(users) as iam_user,
  jsonb_array_elements_text(attached_policy_arns) as attachments
where
  split_part(attachments, '/', 2) = 'AdministratorAccess';
```

```sql+sqlite
select
  name as group_name,
  json_extract(iam_user, '$.UserName') as user_name,
  substr(attachments, instr(attachments, '/')+1) as attached_policies
from
  aws_iam_group
  cross join json_each(users) as iam_user,
  json_each(attached_policy_arns) as attachments
where
  substr(attachments, instr(attachments, '/')+1) = 'AdministratorAccess';
```


### List the policies attached to each IAM group
Discover the segments that are associated with each IAM group in terms of their attached policies. This can be useful in understanding the permissions and access levels of different groups within your AWS environment.

```sql+postgres
select
  name as group_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_group
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments;
```

```sql+sqlite
select
  name as group_name,
  substr(json_each.value, instr(json_each.value, '/') + 1) as attached_policies
from
  aws_iam_group,
  json_each(attached_policy_arns);
```


### Find groups that have inline policies
Determine the areas in which certain groups have inline policies in place, enabling you to better understand and manage your AWS IAM group permissions and security.
```sql+postgres
select
  name as group_name,
  inline_policies
from
  aws_iam_group
where 
  inline_policies is not null;
```

```sql+sqlite
select
  name as group_name,
  inline_policies
from
  aws_iam_group
where 
  inline_policies is not null;
```