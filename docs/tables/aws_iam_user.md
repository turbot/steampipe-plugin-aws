---
title: "Table: aws_iam_user - Query AWS IAM User using SQL"
description: "Allows users to query AWS IAM User data, providing details such as user ID, name, path, creation date, and more. This table is useful for security audits, policy enforcement, and operational troubleshooting."
---

# Table: aws_iam_user - Query AWS IAM User using SQL

The `aws_iam_user` table in Steampipe provides information about IAM users within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query user-specific details, including user ID, name, path, and creation date. Users can utilize this table to gather insights on user permissions, access keys, and associated metadata. The schema outlines the various attributes of the IAM user, including the user ARN, creation date, attached policies, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_user` table, you can use the `.inspect aws_iam_user` command in Steampipe.

**Key columns**:

- `user_name`: The name of the IAM user. This can be used to join with other tables that contain user-specific information.
- `arn`: The Amazon Resource Name (ARN) of the user. This is a unique identifier that can be used to join with any other table that contains ARN information.
- `user_id`: The unique ID assigned by AWS to the user. This can be used to join with other tables that contain user ID information.

## Examples

### Basic IAM user info

```sql
select
  name,
  user_id,
  path,
  create_date,
  password_last_used
from
  aws_iam_user;
```

### Groups details to which the IAM user belongs

```sql
select
  name as user_name,
  iam_group ->> 'GroupName' as group_name,
  iam_group ->> 'GroupId' as group_id,
  iam_group ->> 'CreateDate' as create_date
from
  aws_iam_user
  cross join jsonb_array_elements(groups) as iam_group;
```

### List all the users having Administrator access

```sql
select
  name as user_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_user
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments
where
  split_part(attachments, '/', 2) = 'AdministratorAccess';
```

### List all the users for whom MFA is not enabled

```sql
select
  name,
  user_id,
  mfa_enabled
from
  aws_iam_user
where
  not mfa_enabled;
```

### List the policies attached to each IAM user

```sql
select
  name as user_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_user
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments;
```

### Find users that have inline policies

```sql
select
  name as user_name,
  inline_policies
from
  aws_iam_user
where
  inline_policies is not null;
```
