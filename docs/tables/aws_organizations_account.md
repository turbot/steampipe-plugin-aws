# Table: aws_organizations_account

Contains information about AWS accounts that are members of the current AWS Organizations organization.

This table can only be queried using credentials from an AWS Organizations management account or a member account that is a delegated administrator for an AWS service.

Note: The `account_id` column in this table is the account ID from which the API calls are being made (often the management account). To get the member account's ID, query the `id` column.

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
