---
title: "Table: aws_organizations_account - Query AWS Organizations Account using SQL"
description: "Allows users to query AWS Organizations Account and provides information about each AWS account that is a member of an organization in AWS Organizations."
---

# Table: aws_organizations_account - Query AWS Organizations Account using SQL

The `aws_organizations_account` table in Steampipe provides information about each AWS account that is a member of an organization in AWS Organizations. This table allows DevOps engineers to query account-specific details, including account status, joined method, email, and associated metadata. Users can utilize this table to gather insights on accounts, such as accounts with specific statuses, the method used by the accounts to join the organization, and more. The schema outlines the various attributes of the AWS account, including the account ID, ARN, email, joined method, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_organizations_account` table, you can use the `.inspect aws_organizations_account` command in Steampipe.

### Key columns:

- `account_id` - The unique identifier (ID) of the AWS account that is a member of an organization. This can be used to join with other tables that require account ID.
- `arn` - The Amazon Resource Number (ARN) of the account. It is useful for joining with other tables that use ARN to reference AWS resources.
- `email` - The email address associated with the AWS account. This can be useful for identifying the owner of the account or for sending notifications.

## Examples

### Basic info

```sql
select
  id,
  arn,
  email,
  joined_method,
  joined_timestamp,
  name,
  status,
  tags
from
  aws_organizations_account;
```

### List suspended accounts

```sql
select
  id,
  name,
  arn,
  email,
  joined_method,
  joined_timestamp,
  status
from
  aws_organizations_account
where
  status = 'SUSPENDED';
```
