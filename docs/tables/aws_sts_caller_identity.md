---
title: "Table: aws_sts_caller_identity - Query AWS Security Token Service Caller Identity using SQL"
description: "Allows users to query AWS Security Token Service Caller Identity to retrieve details about the IAM user or role whose credentials are used to call the operation."
---

# Table: aws_sts_caller_identity - Query AWS Security Token Service Caller Identity using SQL

The `aws_sts_caller_identity` table in Steampipe provides information about the AWS Security Token Service (STS) Caller Identity. This table allows users to query details about the IAM user or role whose credentials are used to call the operation. The schema outlines the various attributes of the STS Caller Identity, including the user ARN, user ID, and account ID.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sts_caller_identity` table, you can use the `.inspect aws_sts_caller_identity` command in Steampipe.

### Key columns:

- `account`: The AWS account ID number of the account that owns the resource. Useful for correlating with other tables that also use account ID.
- `arn`: The Amazon Resource Name (ARN) that identifies the caller. This can be used to join with other tables that use ARN.
- `user_id`: The unique identifier for the entity that made the call. This can be used to join with other tables that use user ID.

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