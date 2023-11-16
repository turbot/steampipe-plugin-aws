---
title: "Table: aws_iam_credential_report - Query AWS IAM Credential Reports using SQL"
description: "Allows users to query AWS IAM Credential Reports, providing a comprehensive overview of the AWS Identity and Access Management (IAM) users, their status, and credential usage."
---

# Table: aws_iam_credential_report - Query AWS IAM Credential Reports using SQL

The `aws_iam_credential_report` table in Steampipe provides information about IAM credential reports within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query user-specific details, including access keys, password status, and MFA device usage. Users can utilize this table to gather insights on IAM users, such as inactive users, users with password-enabled login, access key usage, and more. The schema outlines the various attributes of the IAM credential report, including the user name, user creation time, access key details, and password last used date. For more information about the credential report, see [Getting Credential Reports](https://docs.aws.amazon.com/IAM/latest/UserGuide/credential-reports.html) in the IAM User Guide.

_Please note_: This table requires a valid credential report to exist. To generate, please run the follow AWS CLI command:

`aws iam generate-credential-report`

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_credential_report` table, you can use the `.inspect aws_iam_credential_report` command in Steampipe.

**Key columns**:

- `user`: The name of the IAM user. This column is essential as it can be used to join with other tables that contain user-specific information.
- `access_key_1_active`: Indicates whether the first access key is active. This column is crucial for understanding the status and usage of access keys.
- `mfa_active`: Indicates whether the MFA is active for the user. This column provides insights into the security measures taken by the user.

## Examples

### List users that have logged into the console in the past 90 days

```sql
select
  user_name
from
  aws_iam_credential_report
where
  password_enabled
  and password_last_used > (current_date - interval '90' day);
```

### List users that have NOT logged into the console in the past 90 days

```sql
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

### List users with console access that have never logged in to the console

```sql
select
  user_name
from
  aws_iam_credential_report
where
  password_status = 'never_used';
```

### List access keys older than 90 days

```sql
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

### List users that have a console password but do not have MFA enabled

```sql
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

### Check if root login has MFA enabled

```sql
select
  user_name,
  mfa_active
from
  aws_iam_credential_report
where
  user_name = '<root_account>';
```
