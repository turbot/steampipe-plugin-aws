# Table: aws_organizations_account

Contains information about AWS accounts that are members of the current organization.

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
