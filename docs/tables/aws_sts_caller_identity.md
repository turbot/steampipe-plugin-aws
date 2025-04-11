---
title: "Steampipe Table: aws_sts_caller_identity - Query AWS Security Token Service Caller Identity using SQL"
description: "Allows users to query AWS Security Token Service Caller Identity to retrieve details about the IAM user or role whose credentials are used to call the operation."
folder: "STS"
---

# Table: aws_sts_caller_identity - Query AWS Security Token Service Caller Identity using SQL

The AWS Security Token Service (STS) Caller Identity is a resource that provides details about the IAM user or role whose credentials are used to call the operation. It returns the AWS account ID number of the account that owns or contains the calling entity, along with the AWS Access Key ID used to make the call. This service is particularly useful for auditing and tracking purposes, ensuring that all actions within an AWS environment can be traced back to their origin.

## Table Usage Guide

The `aws_sts_caller_identity` table in Steampipe provides you with information about the AWS Security Token Service (STS) Caller Identity. This table allows you to query details about the IAM user or role whose credentials are used to call the operation. The schema outlines for you the various attributes of the STS Caller Identity, including the user ARN, user ID, and account ID.

## Examples

### Basic info
Explore which AWS services are currently being accessed by users, providing a way to monitor usage and manage access permissions effectively. This can be particularly useful for identifying unusual activity or potential security risks.

```sql+postgres
select
  arn,
  user_id,
  title,
  account_id,
  akas
from
  aws_sts_caller_identity;
```

```sql+sqlite
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
This query is useful to identify the specific users who were created using the 'AssumeRole' function within your AWS account. Understanding this information can help maintain security and control over user access and permissions.

```sql+postgres
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

```sql+sqlite
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
Determine the specifics of a user account created through federation, including when it was created and the last time the password was used. This information can be useful for auditing purposes, helping to identify potential security risks or irregularities.

```sql+postgres
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

```sql+sqlite
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