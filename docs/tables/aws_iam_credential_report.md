# Table: aws_iam_credential_report

Retrieves a credential report for the AWS account. For more information about the credential report, see [Getting Credential Reports](https://docs.aws.amazon.com/IAM/latest/UserGuide/credential-reports.html) in the IAM User Guide.

_Please note_: This table requires a valid credential report to exist. To generate, please run the follow AWS CLI command:

`aws iam generate-credential-report`

## Examples

### List Users that have logged into the console in the past 90 days

```sql
select
  user_name
from
  aws_iam_credential_report
where
  password_enabled
  and password_last_used > (current_date - interval '90' day);
```

### Report of users that have NOT logged into the console in the past 90 days?

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

### List of users with console access that have never logged in to the console

```sql
select
    user_name
from
    aws_iam_credential_report
where
    password_never_used;
```

### Find Access Keys older than 90 days

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

### Find users that have a console password but do not have MFA enabled

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
