---
title: "Table: aws_ssm_managed_instance_compliance - Query AWS SSM Managed Instance Compliance using SQL"
description: "Allows users to query AWS SSM Managed Instance Compliance data, providing details on compliance status, compliance type, and related metadata."
---

# Table: aws_ssm_managed_instance_compliance - Query AWS SSM Managed Instance Compliance using SQL

The `aws_ssm_managed_instance_compliance` table in Steampipe provides information about managed instance compliance within AWS Systems Manager (SSM). This table allows DevOps engineers to query compliance-specific details, including compliance status, compliance type, and associated metadata. Users can utilize this table to gather insights on compliance, such as instances that are non-compliant, compliance with specific standards, and more. The schema outlines the various attributes of the managed instance compliance, including the instance ID, compliance type, compliance status, and compliance severity.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_managed_instance_compliance` table, you can use the `.inspect aws_ssm_managed_instance_compliance` command in Steampipe.

### Key columns:

- `instance_id` - This is the unique identifier of the managed instance. It can be used to join this table with other tables that contain instance-specific information.
- `compliance_type` - This column indicates the type of compliance, such as patch, association, or custom. It is useful for filtering data based on specific compliance types.
- `compliance_status` - This column indicates whether the instance is compliant or not. It can be used to identify non-compliant instances for further investigation or remediation.

## Examples

### Basic info

```sql
select
  id,
  name,
  resource_id,
  status,
  compliance_type,
  severity
from
  aws_ssm_managed_instance_compliance
where
  resource_id = 'i-2a3dc8b11ed9d37a';
```

### List non-compliant associations for a managed instance

```sql
select
  id,
  name,
  resource_id as instance_id,
  status,
  compliance_type,
  severity
from
  aws_ssm_managed_instance_compliance
where
  resource_id = 'i-2a3dc8b11ed9d37a'
  and compliance_type = 'Association'
  and status <> 'COMPLIANT';
```

### List non-compliant patches for a managed instance

```sql
select
  id,
  name,
  resource_id as instance_id,
  status,
  compliance_type,
  severity
from
  aws_ssm_managed_instance_compliance
where
  resource_id = 'i-2a3dc8b11ed9d37a'
  and compliance_type = 'Patch'
  and status <> 'COMPLIANT';
```

### List compliance statuses for all managed instances

```sql
select
  c.resource_id as instance_id,
  id,
  status
from
  aws_ssm_managed_instance i,
  aws_ssm_managed_instance_compliance c
where
  i.instance_id = c.resource_id;
```
