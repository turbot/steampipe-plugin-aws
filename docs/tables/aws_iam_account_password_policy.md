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

### Ensure IAM password policy requires at least one uppercase letter (CIS v1.1.05)
```sql
select
  require_uppercase_characters as cis_v1_1_05
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one lowercase letter (CIS v1.1.06)
```sql
select
    require_lowercase_characters  as cis_v1_1_06
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one symbol (CIS v1.1.07)
```sql
select
    require_symbols as cis_v1_1_07
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy require at least one number (CIS v1.1.08)
```sql
select
    require_numbers as cis_v1_1_08
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy requires minimum length of 14 or greater (CIS v1.1.09)
```sql
select
    minimum_password_length >= 14 as cis_v1_1_09
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy prevents password reuse (CIS v1.1.10)
```sql
select
    password_reuse_prevention as cis_v1_1_10
from
    aws_iam_account_password_policy;
```

### Ensure IAM password policy expires passwords within 90 days or less (CIS v1.1.11)
```sql
select
    (expire_passwords and max_password_age <= 90) as cis_v1_1_11
from
    aws_iam_account_password_policy;
```