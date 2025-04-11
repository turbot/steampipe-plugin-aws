---
title: "Steampipe Table: aws_iam_user - Query AWS IAM User using SQL"
description: "Allows users to query AWS IAM User data, providing details such as user ID, name, path, creation date, and more. This table is useful for security audits, policy enforcement, and operational troubleshooting."
folder: "IAM"
---

# Table: aws_iam_user - Query AWS IAM User using SQL

The AWS Identity and Access Management (IAM) User is a resource that represents an individual or application that interacts with AWS. It contains the name, credentials, and permissions to access AWS resources. IAM Users enable the security best practice of granting least privilege, which means granting only the permissions required to perform a task.

## Table Usage Guide

The `aws_iam_user` table in Steampipe provides you with information about IAM users within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query user-specific details, including user ID, name, path, and creation date. You can utilize this table to gather insights on user permissions, access keys, and associated metadata. The schema outlines the various attributes of the IAM user, including the user ARN, creation date, attached policies, and associated tags for you.

## Examples

### Basic IAM user info
Discover the segments that provide details about users in your AWS IAM, including when they were created and when they last used their password. This can be useful for auditing user activity and maintaining security compliance.

```sql+postgres
select
  name,
  user_id,
  path,
  create_date,
  password_last_used
from
  aws_iam_user;
```

```sql+sqlite
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
Determine the affiliations of individual IAM users by identifying the groups they are a part of, providing insights into user access and permissions management within your AWS environment.

```sql+postgres
select
  name as user_name,
  iam_group ->> 'GroupName' as group_name,
  iam_group ->> 'GroupId' as group_id,
  iam_group ->> 'CreateDate' as create_date
from
  aws_iam_user
  cross join jsonb_array_elements(groups) as iam_group;
```

```sql+sqlite
select
  name as user_name,
  json_extract(iam_group, '$.GroupName') as group_name,
  json_extract(iam_group, '$.GroupId') as group_id,
  json_extract(iam_group, '$.CreateDate') as create_date
from
  aws_iam_user,
  json_each(groups) as iam_group;
```

### List all the users having Administrator access
This query helps identify users who have been granted administrator access in an AWS environment. It's useful for auditing user permissions and ensuring only authorized individuals have such high-level access.

```sql+postgres
select
  name as user_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_user
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments
where
  split_part(attachments, '/', 2) = 'AdministratorAccess';
```

```sql+sqlite
select
  name as user_name,
  json_extract(attachments, '$[1]') as attached_policies
from
  aws_iam_user
  cross join json_each(attached_policy_arns) as attachments
where
  json_extract(attachments, '$[1]') = 'AdministratorAccess';
```

### List all the users for whom MFA is not enabled
Discover the users who have not enabled multi-factor authentication, allowing you to identify potential security risks and ensure all accounts are adequately protected.

```sql+postgres
select
  name,
  user_id,
  mfa_enabled
from
  aws_iam_user
where
  not mfa_enabled;
```

```sql+sqlite
select
  name,
  user_id,
  mfa_enabled
from
  aws_iam_user
where
  mfa_enabled = 0;
```

### List the policies attached to each IAM user
Determine the areas in which specific access controls are applied by identifying the policies attached to each user in your AWS IAM service. This can help ensure appropriate security measures are in place and assist in auditing user access rights.

```sql+postgres
select
  name as user_name,
  split_part(attachments, '/', 2) as attached_policies
from
  aws_iam_user
  cross join jsonb_array_elements_text(attached_policy_arns) as attachments;
```

```sql+sqlite
select
  name as user_name,
  json_extract(attachments, '$[1]') as attached_policies
from
  aws_iam_user,
  json_each(attached_policy_arns) as attachments;
```

### Find users that have inline policies
Identify instances where AWS IAM users have inline policies attached to their accounts. This is useful for security audits, as inline policies can grant or deny permissions to AWS services and resources.

```sql+postgres
select
  name as user_name,
  inline_policies
from
  aws_iam_user
where
  inline_policies is not null;
```

```sql+sqlite
select
  name as user_name,
  inline_policies
from
  aws_iam_user
where
  inline_policies is not null;
```