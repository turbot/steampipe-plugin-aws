---
title: "Steampipe Table: aws_organizations_account - Query AWS Organizations Account using SQL"
description: "Allows users to query AWS Organizations Account and provides information about each AWS account that is a member of an organization in AWS Organizations."
folder: "Organizations"
---

# Table: aws_organizations_account - Query AWS Organizations Account using SQL

The AWS Organizations Account is a resource within AWS Organizations service that allows you to centrally manage and govern your environment as you grow and scale your AWS resources. By using AWS Organizations Account, you can create, invite, and manage accounts, set up and apply policies, and consolidate your billing. This helps you to automate AWS account creation and management, and control access to your AWS services.

## Table Usage Guide

The `aws_organizations_account` table in Steampipe provides you with information about each AWS account that is a member of an organization in AWS Organizations. This table allows you, as a DevOps engineer, to query account-specific details, including account status, joined method, email, and associated metadata. You can utilize this table to gather insights on accounts, such as accounts with specific statuses, the method used by the accounts to join the organization, and more. The schema outlines the various attributes of the AWS account, including the account ID, ARN, email, joined method, and status for you.

**Important Notes**
- The `account_id` column in this table is the account ID from which the API calls are being made (often the management account). To get the described member account's ID, query the `id` column.

## Examples

### Basic info
Explore the membership details of your AWS organization accounts, including their status and the method they joined with. This can help in understanding account utilization and managing user access.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are suspended within your organization's account. This is particularly useful for auditing and compliance purposes, allowing you to identify and address any potential issues or risks associated with these accounts.

```sql+postgres
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

```sql+sqlite
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