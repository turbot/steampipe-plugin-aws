---
title: "Table: aws_account - Query AWS Account using SQL"
description: "Allows users to query AWS Account information, including details about the account's status, owner, and associated resources."
---

# Table: aws_account - Query AWS Account using SQL

The `aws_account` table in Steampipe provides information about the AWS Account. This table allows DevOps engineers to query account-specific details, including the account status, owner, and associated resources. Users can utilize this table to gather insights on the AWS account, such as the account's ARN, creation date, email address, and more. The schema outlines the various attributes of the AWS account, including the account ID, account alias, and whether the account is a root account.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_account` table, you can use the `.inspect aws_account` command in Steampipe.

### Key columns:

- `account_id`: This column contains the AWS account ID. It is a unique identifier for the AWS account and can be used to join this table with others that contain AWS account information.
- `account_alias`: This column contains the AWS account alias. It provides a human-readable identifier for the account and can be used for easier identification and querying.
- `is_root`: This column indicates whether the account is a root account. It is useful for understanding the account's level of access and permissions.

## Examples

### Basic AWS account info

```sql
select
  alias,
  arn,
  organization_id,
  organization_master_account_email,
  organization_master_account_id
from
  aws_account
  cross join jsonb_array_elements(account_aliases) as alias;
```

### Organization policy of aws account

```sql
select
  organization_id,
  policy ->> 'Type' as policy_type,
  policy ->> 'Status' as policy_status
from
  aws_account
  cross join jsonb_array_elements(organization_available_policy_types) as policy;
```
