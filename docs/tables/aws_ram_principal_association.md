# Table: aws_ram_principal_association

AWS Resource Access Manager (AWS RAM) helps you securely share the AWS resources that you create in one AWS account with other AWS accounts. If you have multiple AWS accounts, you can create a resource once and use AWS RAM to make that resource usable by those other accounts. It Provides a Resource Access Manager (RAM) principal association. Depending if RAM Sharing with AWS Organizations is enabled, the RAM behavior with different principal types changes.

## Examples

### Basic info

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association;
```

### List permissions attached with each principal associated

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  p ->> 'Arn' as resource_share_permission_arn,
  p ->> 'Status' as resource_share_permission_status
from
  aws_ram_principal_association,
  jsonb_array_elements(resource_share_permission) p;
```

### Get principals that failed association

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association
where
  status = 'FAILED';
```