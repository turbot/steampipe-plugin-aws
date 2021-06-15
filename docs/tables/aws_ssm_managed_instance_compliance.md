# Table: aws_ssm_managed_instance_compliance

A managed instance is any machine configured for AWS Systems Manager. You can configure Amazon Elastic Compute Cloud (Amazon EC2) instances or on-premises machines in a hybrid environment as managed instances. Systems Manager supports various distributions of Linux, including Raspberry Pi devices, macOS, and Microsoft Windows Server.

AWS SSM Managed Instance Compliance provide list of compliance statuses for different resource types for a specified resource ID.

**You must specify an Managed Instance ID** in a `where` clause (`where resource_id='`).

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
  resource_id = 'i-2a49c5a11ed9d37f';
```

### List non-compliant associations of managed instance

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
  resource_id = 'i-2a49c5a11ed9d37f'
  and compliance_type = 'Association'
  and status <> 'COMPLIANT';
```

### List non-compliant patches of managed instance

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
  resource_id = 'i-2a49c5a11ed9d37f'
  and compliance_type = 'Patch'
  and status <> 'COMPLIANT';
```

### List compliance status for all managed instances

```sql
select
  resource_id as instance_id,
  id,
  status
from
  aws_ssm_managed_instance_compliance
where
  resource_id in (
    select
      instance_id
    from
      aws.aws_ssm_managed_instance
  );
```
