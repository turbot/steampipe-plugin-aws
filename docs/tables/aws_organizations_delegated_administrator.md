---
title: "Steampipe Table: aws_organizations_delegated_administrator - Query AWS Organizations Delegated Administrators using SQL"
description: "Allows users to query AWS Organizations Delegated Administrators and provides information about each account delegated administrative privileges within an AWS Organization."
---

# Table: aws_organizations_delegated_administrator - Query AWS Organizations Delegated Administrators using SQL

The AWS Organizations Delegated Administrator is a resource within AWS Organizations service that allows you to query details about your delegated administrator accounts, including their status, the date delegation was enabled, and associated metadata. This is useful for managing and auditing delegated administration within your organization.

## Table Usage Guide

The `aws_organizations_delegated_administrator` table in Steampipe lets you query details of accounts designated as delegated administrators in AWS Organizations. This table allows you, as a DevOps engineer, to identify accounts with specific statuses, view the date delegation was enabled, and gather other relevant account information. The schema outlines various attributes of each delegated administrator account.

**Important Notes:**
- The table provides details about the *delegated administrator* accounts, not the management account making the API calls.
- The `account_id` column in this table is the account ID from which the API calls are being made (often the management account). To get the described member account's ID, query the `delegated_account_id` column.

## Examples

### Basic info
Retrieve basic information about all delegated administrators in your organization.

```sql+postgres
select
  delegated_account_id,
  arn,
  email,
  joined_method,
  joined_timestamp,
  name,
  status,
  delegation_enabled_date
from
  aws_organizations_delegated_administrator;
```

```sql+sqlite
select
  delegated_account_id,
  arn,
  email,
  joined_method,
  joined_timestamp,
  name,
  status,
  delegation_enabled_date
from
  aws_organizations_delegated_administrator;
```

### List delegated administrators with a specific status
Identify delegated administrators with a particular status (e.g., ACTIVE, SUSPENDED, PENDING_CLOSURE).

```sql+postgres
select
  delegated_account_id,
  name,
  arn,
  email,
  joined_method,
  joined_timestamp,
  status,
  delegation_enabled_date
from
  aws_organizations_delegated_administrator
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  delegated_account_id,
  name,
  arn,
  email,
  joined_method,
  joined_timestamp,
  status,
  delegation_enabled_date
from
  aws_organizations_delegated_administrator
where
  status = 'ACTIVE';
```