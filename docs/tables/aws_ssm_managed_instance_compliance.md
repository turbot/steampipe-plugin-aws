# Table: aws_ssm_managed_instance_compliance

A managed instance is any machine configured for AWS Systems Manager. You can configure Amazon Elastic Compute Cloud (Amazon EC2) instances or on-premises machines in a hybrid environment as managed instances. Systems Manager supports various distributions of Linux, including Raspberry Pi devices, macOS, and Microsoft Windows Server.

AWS SSM Managed Instance Compliance provide list of compliance statuses for different resource types for a specified resource ID.

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
  aws_ssm_managed_instance_compliance;
```

### List non compliant associations of managed instances

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
  compliance_type = 'Association'
  and status <> 'COMPLIANT';
```

### List non compliant patches of managed instances

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
  compliance_type = 'Patch'
  and status <> 'COMPLIANT';
```