# Table: aws_sts_caller_identity

The details about the AWS IAM user or role whose credentials are used to call the operation.

## Examples

### Basic info

```sql
select
  arn,
  user_id,
  title,
  account_id,
  akas
from
  aws_sts_caller_identity;
```

### Get the details of the user created with AssumeRole

```sql
select
  caller_identity.arn,
  caller_identity.user_id,
  caller_identity.title,
  caller_identity.account_id,
  u.name,
  u.create_date,
  u.password_last_used
from
  aws_sts_caller_identity as caller_identity,
  aws_iam_user as u
where
  caller_identity.user_id = u.user_id
  and caller_identity.arn like '%assumed%';
```

### Get the details of the user created with GetFederationToken

```sql
select
  caller_identity.arn,
  caller_identity.user_id,
  caller_identity.title,
  caller_identity.account_id,
  u.name,
  u.create_date,
  u.password_last_used
from
  aws_sts_caller_identity as caller_identity,
  aws_iam_user as u
where
  caller_identity.user_id = u.user_id
  and caller_identity.arn like '%federated%';
```