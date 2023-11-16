---
title: "Table: aws_iam_account_summary - Query AWS Identity and Access Management (IAM) Account Summary using SQL"
description: "Allows users to query AWS IAM Account Summary to get a detailed overview of the account's IAM usage and resource consumption."
---

# Table: aws_iam_account_summary - Query AWS Identity and Access Management (IAM) Account Summary using SQL

The `aws_iam_account_summary` table in Steampipe provides information about the AWS IAM Account Summary. This table allows DevOps engineers to query IAM usage and resource consumption details, including users, groups, roles, policies, and more. Users can utilize this table to gather insights on IAM usage, such as the number of users, roles, and policies, and verify the usage against AWS service limits. The schema outlines the various attributes of the IAM Account Summary, including the summary map and account ID.

The number and size of IAM resources in an AWS account are limited. For more information, see [IAM and STS Quotas](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-quotas.html) in the IAM User Guide.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_account_summary` table, you can use the `.inspect aws_iam_account_summary` command in Steampipe.

**Key columns**:

- `account_id`: This is the AWS account ID. This is a key column as it uniquely identifies the AWS account.
- `summary_map`: This is a map of IAM resource types (users, groups, roles, policies) to the count of each resource. This is a key column as it provides a summary of the IAM usage.
- `users_quota`: This is the maximum number of IAM users allowed for the AWS account. This is a key column as it provides information about the service limit for the number of IAM users.

## Examples

### List the IAM summary for the account 
```sql
select
  *
from
  aws_iam_account_summary;
```

### Ensure MFA is enabled for the "root" account (CIS v1.1.13)
```sql
select
  account_mfa_enabled
from
  aws_iam_account_summary;
```




### Summary report - Total number of IAM resources in the account by type
```sql
select
  users,
  groups,
  roles,
  policies
from
  aws_iam_account_summary;
```