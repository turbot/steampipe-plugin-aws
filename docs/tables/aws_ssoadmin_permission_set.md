# Table: aws_ssoadmin_permission_set

Contains information about AWS SSO permission sets.

## Examples

### Basic info

```sql
select
  name,
  arn,
  created_date,
  description,
  relay_state,
  session_duration,
  tags
from
  aws_ssoadmin_permission_set;
```
