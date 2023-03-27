# Table: aws_iam_user

An AWS Identity and Access Management (IAM) user is an entity that you create in AWS to represent the person or application that uses it to interact with.

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
