# Table: aws_iam_account_password_policy

The password policy for the AWS account. For more information about using a password policy, go to [Managing an IAM Password Policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingPasswordPolicies.html).

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

### CIS v1 > 1 Identity and Access Management > 1.05 Ensure IAM password policy requires at least one uppercase letter (Scored)
```sql
select
  require_uppercase_characters as cis_v1_1_05
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.06 Ensure IAM password policy requires at least one lowercase letter (Scored)
```sql
select
    require_lowercase_characters  as cis_v1_1_06
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.07 Ensure IAM password policy requires at least one symbol (Scored)
```sql
select
    require_symbols as cis_v1_1_07
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.08 Ensure IAM password policy require at least one number (Scored)
```sql
select
    require_numbers as cis_v1_1_08
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.09 Ensure IAM password policy requires minimum length of 14 or greater (Scored)
```sql
select
    minimum_password_length >= 14 as cis_v1_1_09
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.10 Ensure IAM password policy prevents password reuse (Scored)
```sql
select
    password_reuse_prevention as cis_v1_1_10
from
    aws_iam_account_password_policy;
```

### CIS v1 > 1 Identity and Access Management > 1.11 Ensure IAM password policy expires passwords within 90 days or less (Scored)
```sql
select
    (expire_passwords and max_password_age <= 90) as cis_v1_1_11
from
    aws_iam_account_password_policy;
```