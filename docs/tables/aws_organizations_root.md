---
title: "Steampipe Table: aws_organizations_root - Query AWS Organizations Root using SQL"
description: "Allows users to query AWS Organizations Root to retrieve detailed information on AWS Organizations Root account. This table can be utilized to gain insights on organizations root account."
folder: "Organizations"
---

# Table: aws_organizations_root - Query AWS Organizations Root using SQL

AWS Organizations uses a hierarchical structure to manage accounts. At the top of this hierarchy is the "root." The root is the starting point for organizing your AWS accounts. The root acts as the parent container for all the accounts in your organization. It can also contain organizational units (OUs), which are sub-containers that can themselves contain accounts or further nested OUs.

## Table Usage Guide

The `aws_organizations_root` table in Steampipe provides you the information about AWS Organizations Root Account.

## Examples

### Basic info
It's particularly useful in contexts where managing or auditing AWS Organizations.

```sql+postgres
select
  name,
  id,
  arn
from
  aws_organizations_root;
```

```sql+sqlite
select
  name,
  id,
  arn
from
  aws_organizations_root;
```

### Get the policy details attached to organization root account
The types of policies that are currently enabled for the root and therefore can be attached to the root or to its OUs or accounts.

```sql+postgres
select
  id,
  name,
  p ->> 'Status' as policy_status,
  p ->> 'Type' as policy_type
from
  aws_organizations_root,
  jsonb_array_elements(policy_types) as p;
```

```sql+sqlite
select
  id,
  name,
  json_extract(json_each.value, '$.Status') AS policy_status,
  json_extract(json_each.value, '$.Type') AS policy_type
from
  aws_organizations_root,
  json_each(policy_types) as p;
```