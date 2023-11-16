---
title: "Table: aws_iam_account_password_policy - Query AWS IAM Account Password Policies using SQL"
description: "Allows users to query AWS IAM Account Password Policies to gain insights about password policy details such as minimum password length, password expiration period, and whether it requires at least one number or symbol."
---

# Table: aws_iam_account_password_policy - Query AWS IAM Account Password Policies using SQL

The `aws_iam_account_password_policy` table in Steampipe provides information about IAM account password policies within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query password policy-specific details, including minimum password length, password expiration period, and whether it requires at least one number or symbol. Users can utilize this table to gather insights on password policies, such as password complexity requirements, password rotation policies, and more. The schema outlines the various attributes of the IAM account password policy, including the allow users to change password, hard expiry, and password reuse prevention.

For more information about using a password policy, go to [Managing an IAM Password Policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingPasswordPolicies.html).

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_account_password_policy` table, you can use the `.inspect aws_iam_account_password_policy` command in Steampipe.

### Key columns:

- `title`: The title of the password policy. This is useful for identifying the specific password policy.
- `minimum_password_length`: This column is important as it gives information about the minimum password length required by the policy.
- `max_password_age`: This column is useful for understanding the password expiration period set by the policy.

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
  require_uppercase_characters
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one lowercase letter (CIS v1.1.06)
```sql
select
  require_lowercase_characters
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one symbol (CIS v1.1.07)
```sql
select
  require_symbols
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy require at least one number (CIS v1.1.08)
```sql
select
  require_numbers
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires minimum length of 14 or greater (CIS v1.1.09)
```sql
select
  minimum_password_length >= 14
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy prevents password reuse (CIS v1.1.10)
```sql
select
  password_reuse_prevention
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy expires passwords within 90 days or less (CIS v1.1.11)
```sql
select
  (expire_passwords and max_password_age <= 90)
from
  aws_iam_account_password_policy;
```