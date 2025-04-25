---
title: "Steampipe Table: aws_account - Query AWS Accounts using SQL"
description: "Allows users to query AWS Account information, including details about the account's status, owner, and associated resources."
folder: "Account"
---

# Table: aws_account - Query AWS Accounts using SQL

The AWS Account is a container for AWS resources. It is used to sign up for, organize, and manage AWS services, and it provides administrative control access to resources. An AWS Account contains its own data, with its own settings, including billing and payment information.

## Table Usage Guide

The `aws_account` table in Steampipe provides you with information about your AWS Account. This table allows you, as a DevOps engineer, to query account-specific details, including the account status, owner, and associated resources. You can utilize this table to gather insights on your AWS account, such as the account's ARN, creation date, email address, and more. The schema outlines the various attributes of your AWS account, including the account ID, account alias, and whether your account is a root account.

## Examples

### Basic AWS account info
Discover the segments that are associated with your AWS account, including details about the organization and the master account. This can help you manage and understand the relationships within your AWS structure.This query provides a snapshot of basic details about your AWS account, including its alias and associated organization details. It's useful for quickly accessing key information about your account, particularly in larger organizations where multiple accounts may be in use.


```sql+postgres
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

```sql+sqlite
select
  alias.value as alias,
  arn,
  organization_id,
  organization_master_account_email,
  organization_master_account_id
from
  aws_account,
  json_each(account_aliases) as alias;
```

### Organization policy of aws account
This query allows you to delve into the various policies within your AWS account, particularly focusing on the type and status of each policy. It's useful for managing and tracking policy configurations across your organization, ensuring compliance and efficient resource utilization.This query is used to understand the types and status of policies available for an AWS organization. This can be beneficial for auditing purposes, ensuring policy compliance across all accounts within the organization.


```sql+postgres
select
  organization_id,
  policy ->> 'Type' as policy_type,
  policy ->> 'Status' as policy_status
from
  aws_account
  cross join jsonb_array_elements(organization_available_policy_types) as policy;
```

```sql+sqlite
select
  organization_id,
  json_extract(policy.value, '$.Type') as policy_type,
  json_extract(policy.value, '$.Status') as policy_status
from
  aws_account,
  json_each(organization_available_policy_types) as policy;
```