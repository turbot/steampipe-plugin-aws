# Table: aws_ram_principal_association

AWS Resource Access Manager (RAM) helps you securely share the AWS resources you create in one AWS account with other AWS accounts. If you have multiple AWS accounts, you can create a resource once and use AWS RAM to make that resource usable by those other accounts.

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