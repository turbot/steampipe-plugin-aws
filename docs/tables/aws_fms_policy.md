# Table: aws_fms_policy

An AWS Firewall Manager Policy is a set of security rules and settings that you can apply to your AWS accounts and resources within those accounts. These policies help you enforce security and compliance standards consistently across your organization, making it easier to manage multiple accounts and resources.

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