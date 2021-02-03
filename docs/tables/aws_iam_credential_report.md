# Table: aws_iam_credential_report

Retrieves a credential report for the AWS account. For more information about the credential report, see [Getting Credential Reports](https://docs.aws.amazon.com/IAM/latest/UserGuide/credential-reports.html) in the IAM User Guide.

_Please note_: This table requires a valid credential report to exist. To generate, please run the follow AWS CLI command:

`aws iam generate-credential-report`

## Examples

### Who has logged into the console in the past 90 days?
```sql
select user_name
from aws_iam_credential_report
where password_enabled
and password_last_used > (current_date - interval '90' day);
```

### Who has console access and has never logged in?
```sql
select
    user_name
from
    aws_iam_credential_report
where
    password_enabled
and password_last_used is null;
```