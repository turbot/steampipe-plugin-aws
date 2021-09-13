# Table: aws_organizations_account

Contains information about AWS accounts that are members of the current AWS Organizations organization.

Note: This table can only be queried using IAM access key credentials for the AWS Organizations management account or a member account that is a delegated administrator for an AWS service.

## Examples

### Basic info

```sql
select
  id,
  arn,
  email,
  joined_method,
  joined_timestamp,
  name,
  status,
  tags
from
  aws_organizations_account;
```

### List suspended accounts

```sql
select
  id,
  name,
  arn,
  email,
  joined_method,
  joined_timestamp,
  status
from
  aws_organizations_account
where
  status = 'SUSPENDED';
```