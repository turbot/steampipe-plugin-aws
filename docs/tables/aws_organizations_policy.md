---
title: "Table: aws_organizations_policy - Query AWS Organizations Policy using SQL"
description: "Allows users to query AWS Organizations Policy to retrieve detailed information on policies within AWS Organizations. This table can be utilized to gain insights on policy-specific details, such as policy type, content, and associated metadata."
---

# Table: aws_organizations_policy - Query AWS Organizations Policy using SQL

The `aws_organizations_policy` table in Steampipe provides information about policies within AWS Organizations. This table allows DevOps engineers to query policy-specific details, including policy type, content, and associated metadata. Users can utilize this table to gather insights on policies, such as policy names, policy types, and the contents of the policies. The schema outlines the various attributes of the policy, including the policy ARN, policy type, policy content, policy name, and associated tags.

**Note**: The `type` column in this table is required to make the API call.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_organizations_policy` table, you can use the `.inspect aws_organizations_policy` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the policy. This can be used to join this table with other tables that contain policy ARN as a column.
- `name`: The friendly name of the policy. This can be used to join this table with other tables that contain policy name as a column.
- `type`: The type of the policy. This can be useful in filtering policies based on their types.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```

### List tag policies that are not managed by AWS

```sql
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  not aws_managed
  and type = 'TAG_POLICY';
```

### List backup policies

```sql
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy
where
  type = 'BACKUP_POLICY';
```

### Get policy details of the service control policies

```sql
select
  name,
  id,
  content ->> 'Version' as policy_version,
  content ->> 'Statement' as policy_statement
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```
