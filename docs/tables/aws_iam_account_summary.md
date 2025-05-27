---
title: "Steampipe Table: aws_iam_account_summary - Query AWS Identity and Access Management (IAM) Account Summary using SQL"
description: "Allows users to query AWS IAM Account Summary to get a detailed overview of the account's IAM usage and resource consumption."
folder: "IAM"
---

# Table: aws_iam_account_summary - Query AWS Identity and Access Management (IAM) Account Summary using SQL

The AWS Identity and Access Management (IAM) Account Summary provides an overview of your AWS security settings including users, groups, roles, and policies in your account. This service is useful for auditing and monitoring purposes, allowing you to ensure your account is secure and compliant with your organization's policies. It provides a user-friendly SQL interface for querying your IAM settings.

## Table Usage Guide

The `aws_iam_account_summary` table in Steampipe provides you with information about the AWS IAM Account Summary. This table allows you, as a DevOps engineer, to query IAM usage and resource consumption details, including users, groups, roles, policies, and more. You can utilize this table to gather insights on IAM usage, such as the number of users, roles, and policies, and verify the usage against AWS service limits. The schema outlines the various attributes of the IAM Account Summary, including the summary map and account ID.

**Important Notes**
- The number and size of IAM resources in your AWS account are limited. For more information, see [IAM and STS Quotas](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-quotas.html) in the IAM User Guide.

## Examples

### List the IAM summary for the account 
Analyze the general overview of your AWS Identity and Access Management (IAM) to gain insights into user access and permissions within your account. This could be beneficial in identifying potential security risks or for general account management.

```sql+postgres
select
  *
from
  aws_iam_account_summary;
```

```sql+sqlite
select
  *
from
  aws_iam_account_summary;
```

### Ensure MFA is enabled for the "root" account (CIS v1.1.13)
Determine the areas in which Multi-Factor Authentication (MFA) is activated for the primary account to enhance security measures as per CIS v1.1.13 guidelines.

```sql+postgres
select
  account_mfa_enabled
from
  aws_iam_account_summary;
```

```sql+sqlite
select
  account_mfa_enabled
from
  aws_iam_account_summary;
```




### Summary report - Total number of IAM resources in the account by type
Determine the distribution of different types of Identity and Access Management (IAM) resources in your AWS account. This can help you understand the composition of your IAM resources and manage them more effectively.

```sql+postgres
select
  users,
  groups,
  roles,
  policies
from
  aws_iam_account_summary;
```

```sql+sqlite
select
  users,
  groups,
  roles,
  policies
from
  aws_iam_account_summary;
```