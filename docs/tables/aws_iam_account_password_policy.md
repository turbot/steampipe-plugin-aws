# Table: aws_iam_account_password_policy

Retrieves the password policy for the AWS account. For more information about using a password policy, go to [Managing an IAM Password Policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingPasswordPolicies.html).

## Examples


### List the password policy for the account
```sql
select
  allow_users_to_change_password,
  expire_passwords,
  hard_expiry,
  max_password_age,
  minimum_password_length,
  password_reuse_prevention,
  require_lowercase_characters,
  require_numbers,
  require_symbols,
  require_uppercase_characters
from
    aws_iam_account_password_policy;
```