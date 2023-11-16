---
title: "Table: aws_organizations_policy_target - Query AWS Organizations Policy Targets using SQL"
description: "Allows users to query AWS Organizations Policy Targets to retrieve detailed information about the application of policies to roots, organizational units (OUs), and accounts."
---

# Table: aws_organizations_policy_target - Query AWS Organizations Policy Targets using SQL

The `aws_organizations_policy_target` table in Steampipe provides information about policy targets within AWS Organizations. This table allows DevOps engineers to query policy target-specific details, including the policy ID, target ID, and the type of target (root, OU, or account). Users can utilize this table to gather insights on policy applications, such as which policies are applied to which roots, OUs, or accounts, and the status of these applications. The schema outlines the various attributes of the policy target, including the ARN, policy ID, target ID, and target type.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_organizations_policy_target` table, you can use the `.inspect aws_organizations_policy_target` command in Steampipe.

**Key columns**:

- `policy_id`: The unique identifier (ID) of the policy that is attached to the target. This is useful for referencing the specific policy in AWS Organizations.
- `target_id`: The unique identifier (ID) of the root, organizational unit, or account. This is useful for identifying the specific target to which the policy is attached.
- `type`: The type of the target (root, OU, or account). This is important for understanding the scope of the policy application.

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
  aws_organizations_policy_target 
where
  type = 'SERVICE_CONTROL_POLICY' 
  and target_id = '123456789098';
```

### List tag policies of a targeted organization that are not managed by AWS

```sql
select
  id,
  name,
  arn,
  type,
  aws_managed 
from
  aws_organizations_policy_target 
where
  not aws_managed 
  and type = 'TAG_POLICY' 
  and target_id = 'ou-jsdhkek';
```

### List backup organization policies of an account

```sql
select
  id,
  name,
  arn,
  type,
  aws_managed
from
  aws_organizations_policy_target
where
  type = 'BACKUP_POLICY'
  and target_id = '123456789098';
```

### Get policy details of the service control policies of a root account

```sql
select
  name,
  id,
  content ->> 'Version' as policy_version,
  content ->> 'Statement' as policy_statement
from
  aws_organizations_policy
where
  type = 'SERVICE_CONTROL_POLICY'
  and target_id = 'r-9ijkl7';
```
