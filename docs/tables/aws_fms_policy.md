---
title: "Table: aws_fms_policy - Query AWS Firewall Manager Policies using SQL"
description: "Allows users to query AWS Firewall Manager Policies using SQL. This table provides information about each AWS Firewall Manager (FMS) policy in an AWS account. It can be used to gain insights into policy details such as the policy name, ID, resource type, security service type, and more."
---

# Table: aws_fms_policy - Query AWS Firewall Manager Policies using SQL

The `aws_fms_policy` table in Steampipe provides information about each AWS Firewall Manager (FMS) policy in an AWS account. This table allows DevOps engineers, security professionals, and other users to query policy-specific details, including policy ID, policy name, resource type, security service type, and more. Users can utilize this table to gather insights on policies, such as identifying which resources are protected by which policies, the type of security service provided by each policy, and the status of each policy. The schema outlines the various attributes of the FMS policy, including the policy ID, policy name, resource tags, remediation enabled status, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_fms_policy` table, you can use the `.inspect aws_fms_policy` command in Steampipe.

### Key columns:

- `policy_id`: The ID of the AWS FMS policy. This is a unique identifier for each policy and can be used to join this table with other tables that contain policy information.
- `policy_name`: The name of the AWS FMS policy. This provides a human-readable identifier for each policy and can be useful for identifying policies in queries and reports.
- `resource_type`: The type of AWS resource that the policy applies to. This can be useful for identifying which resources are protected by which policies.

## Examples

### Basic info

```sql
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type
from
  aws_fms_policy;
```

### List policies that has remediation enabled

```sql
select
  policy_name,
  policy_id,
  arn,
  policy_description,
  resource_type,
  remediation_enabled
from
  aws_fms_policy
where
  remediation_enabled;
```

### Count policies by resource type

```sql
select
  policy_name,
  resource_type,
  count(policy_id) as policy_applied
from
  aws_fms_policy
group by
  policy_name,
  resource_type;
```

### List policies that are not active

```sql
select
  policy_name,
  policy_id,
  policy_status
from
  aws_fms_policy
where
  policy_status <> 'ACTIVE';
```