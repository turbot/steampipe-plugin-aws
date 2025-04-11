---
title: "Steampipe Table: aws_iam_account_password_policy - Query AWS IAM Account Password Policies using SQL"
description: "Allows users to query AWS IAM Account Password Policies to gain insights about password policy details such as minimum password length, password expiration period, and whether it requires at least one number or symbol."
folder: "IAM"
---

# Table: aws_iam_account_password_policy - Query AWS IAM Account Password Policies using SQL

The AWS Identity and Access Management (IAM) Account Password Policy is a resource that allows you to manage the password policy for your AWS account. This includes settings like the minimum password length, whether to require symbols, numbers, or uppercase letters, and whether to allow users to change their own password. It can help you enforce strong password practices in your organization.

## Table Usage Guide

The `aws_iam_account_password_policy` table in Steampipe provides you with information about IAM account password policies within AWS Identity and Access Management (IAM). This table enables you, as a DevOps engineer, to query password policy-specific details, including minimum password length, password expiration period, and whether it requires at least one number or symbol. You can utilize this table to gather insights on password policies, such as password complexity requirements, password rotation policies, and more. The schema outlines the various attributes of the IAM account password policy, including the option for users to change password, hard expiry, and password reuse prevention.

**Important Notes**
- For more information about using a password policy, you can visit [Managing an IAM Password Policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_ManagingPasswordPolicies.html).

## Examples


### List the password policy for the account
Gain insights into your AWS account's security by examining its password policy. This query can help you understand the strength and complexity requirements of your passwords, which can aid in enhancing your account's security.
```sql+postgres
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

```sql+sqlite
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
Determine whether your AWS IAM account password policy mandates the inclusion of at least one uppercase letter, to ensure enhanced security and compliance with CIS v1.1.05 standards.

```sql+postgres
select
  require_uppercase_characters
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  require_uppercase_characters
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one lowercase letter (CIS v1.1.06)
Determine the areas in which your AWS IAM password policy mandates the inclusion of at least one lowercase letter. This query is useful for ensuring your password policy aligns with the CIS v1.1.06 benchmark, enhancing system security.

```sql+postgres
select
  require_lowercase_characters
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  require_lowercase_characters
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires at least one symbol (CIS v1.1.07)
Determine the areas in which your AWS IAM account password policy mandates the inclusion of at least one symbol. This can be useful in enhancing the security of your system by enforcing stronger password requirements.

```sql+postgres
select
  require_symbols
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  require_symbols
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy require at least one number (CIS v1.1.08)
Determine the areas in which your AWS IAM password policy mandates the inclusion of at least one numerical digit, which is a recommended security measure according to CIS v1.1.08.

```sql+postgres
select
  require_numbers
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  require_numbers
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy requires minimum length of 14 or greater (CIS v1.1.09)
Determine the areas in which your AWS IAM account password policy adheres to the CIS v1.1.09 standard, which requires a minimum password length of 14 or greater. This query can help enhance security by ensuring password complexity.

```sql+postgres
select
  minimum_password_length >= 14
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  minimum_password_length >= 14 as 'minimum_password_length'
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy prevents password reuse (CIS v1.1.10)
Determine the areas in which your AWS IAM password policy restricts reusing old passwords. This is crucial to reinforce security measures and mitigate the risk of unauthorized access.

```sql+postgres
select
  password_reuse_prevention
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  password_reuse_prevention
from
  aws_iam_account_password_policy;
```

### Ensure IAM password policy expires passwords within 90 days or less (CIS v1.1.11)
Assess the elements within your AWS IAM account password policy to ensure that passwords expire within a 90-day period. This is crucial for maintaining a robust security posture and aligning with CIS benchmark recommendations.

```sql+postgres
select
  (expire_passwords and max_password_age <= 90)
from
  aws_iam_account_password_policy;
```

```sql+sqlite
select
  (expire_passwords = 1 and max_password_age <= 90)
from
  aws_iam_account_password_policy;
```