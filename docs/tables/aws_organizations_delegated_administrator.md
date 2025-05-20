---
title: "Steampipe Table: aws_organizations_delegated_administrator - Query AWS Organizations Delegated Administrators using SQL"
description: "Allows users to query AWS Organizations Delegated Administrators and provides information about each account delegated administrative privileges within an AWS Organization."
folder: "Organizations"
---

# Table: aws_organizations_delegated_administrator - Query AWS Organizations Delegated Administrators using SQL

The AWS Organizations Delegated Administrator is a resource within AWS Organizations service that allows you to query details about your delegated administrator accounts, including their status, the date delegation was enabled, and associated metadata. This is useful for managing and auditing delegated administration within your organization.

## Table Usage Guide

The `aws_organizations_delegated_administrator` table in Steampipe lets you query details of accounts designated as delegated administrators in AWS Organizations. This table allows you, as a DevOps engineer, to identify accounts with specific statuses, view the date delegation was enabled, and gather other relevant account information. The schema outlines various attributes of each delegated administrator account.

**Important Notes:**
* This table returns details about *delegated administrator* accounts, **not** the management account executing the API call.
* The `account_id` column shows the ID of the account that made the API request (typically the management account). To retrieve the ID of the delegated administrator (member account), refer to the `id` column.

## Examples

### Basic info
Retrieve basic information about all delegated administrators in your organization.

```sql+postgres
select
  id,
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
  id,
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
  id,
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
  id,
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