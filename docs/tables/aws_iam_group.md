# Table: aws_iam_group

An IAM group is a collection of IAM users. Groups let you specify permissions for multiple users, which makes it easier to manage the permissions for those users.

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