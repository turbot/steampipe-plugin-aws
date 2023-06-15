# Table: aws_organizations_policy_target

AWS organization policy target refers to the entity to which an organization policy is applied. It can be an AWS account or an organizational unit (OU) within the AWS organization. The policy is enforced on the specified target, defining the rules and restrictions that govern the resources and actions within that target. By setting up organization policies and assigning them to appropriate targets, administrators can effectively manage and control the AWS resources and actions within their organization.

**Note:** The `type` and `target_id` columns in this table are required to make the API call.

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
and
  target_id = '123456789098';
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
and
  type = 'TAG_POLICY'
and
  target_id = 'ou-jsdhkek';
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
and
  target_id = '123456789098';
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
and
  target_id = 'r-9ijkl7';
```