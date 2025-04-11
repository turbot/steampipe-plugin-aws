---
title: "Steampipe Table: aws_ssm_managed_instance_compliance - Query AWS SSM Managed Instance Compliance using SQL"
description: "Allows users to query AWS SSM Managed Instance Compliance data, providing details on compliance status, compliance type, and related metadata."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_managed_instance_compliance - Query AWS SSM Managed Instance Compliance using SQL

The AWS Systems Manager Managed Instance Compliance is a resource that enables you to check the compliance of your managed instances. It uses SQL queries to assess the configuration compliance of your instances according to the policies you've defined. This allows you to ensure your instances are running in accordance with your organization's security and operational best practices.

## Table Usage Guide

The `aws_ssm_managed_instance_compliance` table in Steampipe provides you with information about managed instance compliance within AWS Systems Manager (SSM). This table allows you, as a DevOps engineer, to query compliance-specific details, including compliance status, compliance type, and associated metadata. You can utilize this table to gather insights on compliance, such as instances that are non-compliant, compliance with specific standards, and more. The schema outlines the various attributes of the managed instance compliance for you, including the instance ID, compliance type, compliance status, and compliance severity.

**Important Notes**
- You must specify an Managed Instance ID in a `where` or `join` clause (`where resource_id='`) to query this table.

## Examples

### Basic info
Determine the compliance status and severity level of a specific AWS SSM managed instance. This is useful to identify potential security risks and ensure adherence to compliance standards.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which a managed instance is non-compliant. This query is beneficial in identifying instances where the compliance type is 'Association' and the status is not 'COMPLIANT', providing insights into potential areas for improvement.

```sql+postgres
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

```sql+sqlite
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
  and status != 'COMPLIANT';
```

### List non-compliant patches for a managed instance
Determine the areas in which patches for a managed instance are non-compliant. This can assist in identifying potential vulnerabilities and ensuring system security.

```sql+postgres
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

```sql+sqlite
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
  and status != 'COMPLIANT';
```

### List compliance statuses for all managed instances
Determine the compliance status of all managed instances to ensure adherence to standards and regulations. This is useful in maintaining system integrity and mitigating potential risks.

```sql+postgres
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

```sql+sqlite
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