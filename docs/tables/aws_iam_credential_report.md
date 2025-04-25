---
title: "Steampipe Table: aws_iam_credential_report - Query AWS IAM Credential Reports using SQL"
description: "Allows users to query AWS IAM Credential Reports, providing a comprehensive overview of the AWS Identity and Access Management (IAM) users, their status, and credential usage."
folder: "IAM"
---

# Table: aws_iam_credential_report - Query AWS IAM Credential Reports using SQL

The AWS IAM Credential Report is a document that provides details about how the AWS Identity and Access Management (IAM) users in your AWS account are accessing AWS services. It lists all your AWS account's users and the status of their various credentials, including passwords, access keys, MFA devices, and signing certificates. This report can help you audit and improve the security of your AWS account.

## Table Usage Guide

The `aws_iam_credential_report` table in Steampipe provides you with information about IAM credential reports within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query user-specific details, including access keys, password status, and MFA device usage. You can utilize this table to gather insights on IAM users, such as inactive users, users with password-enabled login, access key usage, and more. The schema outlines the various attributes of the IAM credential report, including the user name, user creation time, access key details, and password last used date. For more information about the credential report, see [Getting Credential Reports](https://docs.aws.amazon.com/IAM/latest/UserGuide/credential-reports.html) in the IAM User Guide.

**Important Notes**
- You need a valid credential report to exist for this table. To generate one, please run the following AWS CLI command - `aws iam generate-credential-report`.

## Examples

### List users that have logged into the console in the past 90 days
Determine the users who have accessed the console within the past three months. This can be useful for monitoring user activity, identifying potential security risks, or auditing user access for compliance purposes.

```sql+postgres
select
  user_name
from
  aws_iam_credential_report
where
  password_enabled
  and password_last_used > (current_date - interval '90' day);
```

```sql+sqlite
select
  user_name
from
  aws_iam_credential_report
where
  password_enabled = 1
  and password_last_used > date('now','-90 day');
```

### List users that have NOT logged into the console in the past 90 days
Identify users who may have abandoned their accounts or are no longer active by pinpointing those who haven't logged in for the past 90 days. This assists in maintaining secure and efficient user management by flagging potential inactive accounts for review or deletion.

```sql+postgres
select
  user_name,
  password_last_used,
  age(password_last_used)
from
  aws_iam_credential_report
where
  password_enabled
  and password_last_used <= (current_date - interval '90' day)
order by
  password_last_used;
```

```sql+sqlite
select
  user_name,
  password_last_used,
  julianday('now') - julianday(password_last_used)
from
  aws_iam_credential_report
where
  password_enabled
  and julianday('now') - julianday(password_last_used) >= 90
order by
  password_last_used;
```

### List users with console access that have never logged in to the console
Discover the segments of users who have been granted console access but have never utilized it. This can be useful in identifying unnecessary access privileges and enhancing security measures.

```sql+postgres
select
  user_name
from
  aws_iam_credential_report
where
  password_status = 'never_used';
```

```sql+sqlite
select
  user_name
from
  aws_iam_credential_report
where
  password_status = 'never_used';
```

### List access keys older than 90 days
Discover the segments that have access keys older than 90 days to assess potential security risks and ensure timely key rotation. This can help maintain secure access protocols and prevent unauthorized access.

```sql+postgres
select
  user_name,
  access_key_1_last_rotated,
  age(access_key_1_last_rotated) as access_key_1_age,
  access_key_2_last_rotated,
  age(access_key_2_last_rotated) as access_key_2_age
from
  aws_iam_credential_report
where
  access_key_1_last_rotated <= (current_date - interval '90' day)
  or access_key_2_last_rotated <= (current_date - interval '90' day)
order by
  user_name;
```

```sql+sqlite
select
  user_name,
  access_key_1_last_rotated,
  julianday('now') - julianday(access_key_1_last_rotated) as access_key_1_age,
  access_key_2_last_rotated,
  julianday('now') - julianday(access_key_2_last_rotated) as access_key_2_age
from
  aws_iam_credential_report
where
  julianday('now') - julianday(access_key_1_last_rotated) >= 90
  or julianday('now') - julianday(access_key_2_last_rotated) >= 90
order by
  user_name;
```

### List users that have a console password but do not have MFA enabled
Explore which users have an active console password but lack multi-factor authentication. This is useful for identifying potential security vulnerabilities within your AWS IAM user base.

```sql+postgres
select
  user_name,
  mfa_active,
  password_enabled
from
  aws_iam_credential_report
where
  password_enabled
  and not mfa_active;
```

```sql+sqlite
select
  user_name,
  mfa_active,
  password_enabled
from
  aws_iam_credential_report
where
  password_enabled = 1
  and mfa_active = 0;
```

### Check if root login has MFA enabled
Determine if the root account of your AWS IAM service has multifactor authentication (MFA) enabled. This is crucial for enhancing account security and preventing unauthorized access.

```sql+postgres
select
  user_name,
  mfa_active
from
  aws_iam_credential_report
where
  user_name = '<root_account>';
```

```sql+sqlite
select
  user_name,
  mfa_active
from
  aws_iam_credential_report
where
  user_name = '<root_account>';
```