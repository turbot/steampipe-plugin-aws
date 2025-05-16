---
title: "Steampipe Table: aws_organizations_delegated_services_for_account - Query AWS Organizations Delegated Services for an Account using SQL"
description: "Allows users to query AWS Organizations delegated services for a specific account, providing details on which services have been granted delegated administrative privileges."
folder: "Organizations"
---

# Table: aws_organizations_delegated_services_for_account - Query AWS Organizations Delegated Services for an Account using SQL

The AWS Organizations Delegated Services for an Account is a resource that allows you to query information about the services that have been delegated administrative privileges for a given AWS account within an AWS Organization service. This allows you to query details about these services, including the service principal and the date the delegation was enabled. This is useful for managing and auditing delegated service administration within your organization.

## Table Usage Guide

The `aws_organizations_delegated_administrator` table in Steampipe lets you query the services granted delegated administration access to a specific account. You specify the delegated account ID, and the table returns a list of services and their delegation status. The schema outlines various attributes of each delegated service.

**Important Notes:**
- This table supports the optional list key column `delegated_account_id`.
- The `delegated_account_id` is the ID of the account which has been granted delegated administration, *not* the management account making the API calls.

## Examples

### Basic info
Retrieve basic information about all delegated services.

```sql+postgres
select
  delegated_account_id,
  service_principal,
  delegation_enabled_date
from
  aws_organizations_delegated_services_for_account
```

```sql+sqlite
select
  delegated_account_id,
  service_principal,
  delegation_enabled_date
from
  aws_organizations_delegated_services_for_account
```

### Basic info for a specific account
Retrieve basic information about all delegated services for a specific account.

```sql+postgres
select
  delegated_account_id,
  service_principal,
  delegation_enabled_date
from
  aws_organizations_delegated_services_for_account
where
  delegated_account_id = '123456789012';
```

```sql+sqlite
select
  delegated_account_id,
  service_principal,
  delegation_enabled_date
from
  aws_organizations_delegated_services_for_account
where
  delegated_account_id = '123456789012';
```