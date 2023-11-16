---
title: "Table: aws_iam_group - Query AWS IAM Group using SQL"
description: "Allows users to query AWS IAM Group data such as group name, path, and ARN. This table provides information about IAM groups within AWS Identity and Access Management (IAM)."
---

# Table: aws_iam_group - Query AWS IAM Group using SQL

The `aws_iam_group` table in Steampipe provides information about IAM groups within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query group-specific details, including group name, path, and ARN. Users can utilize this table to gather insights on groups, such as group membership, group policy attachments, and more. The schema outlines the various attributes of the IAM group, including the group ARN, creation date, group ID, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_group` table, you can use the `.inspect aws_iam_group` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) specifying the group. This can be used to join with other tables where group ARN is required.
- `group_name`: The friendly name that identifies the group. This can be used to join with other tables where group name is required.
- `group_id`: The stable and unique string identifying the group. This can be used to join with other tables where group ID is required.

## Examples

### User details associated with each IAM group

```sql
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


### List all the users in each group having Administrator access

```sql
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


### List the policies attached to each IAM group

```sql
select
  name as group_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_group
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments;
```


### Find groups that have inline policies
```sql
select
  name as group_name,
  inline_policies
from
  aws_iam_group
where 
  inline_policies is not null;
```