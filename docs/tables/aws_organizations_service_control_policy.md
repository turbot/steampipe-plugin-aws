# Table: aws_organizations_service_control_policy

Policies in AWS Organizations enable you to apply additional types of management to the AWS accounts in your organization.

Note: The `type` column in this table is required to make the API call.

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
  aws_organizations_service_control_policy
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
  aws_organizations_service_control_policy
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
  aws_organizations_service_control_policy
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
  aws_organizations_service_control_policy
where
  type = 'SERVICE_CONTROL_POLICY';
```
