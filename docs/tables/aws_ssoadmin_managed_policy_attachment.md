# Table: aws_ssoadmin_managed_policy_attachment

Contains information about managed IAM policies attached to AWS SSO permission sets.

## Examples

### Basic info

```sql
select
  mpa.managed_policy_arn,
  mpa.name
from
  aws_ssoadmin_managed_policy_attachment as mpa
join
  aws_ssoadmin_permission_set as ps on mpa.permission_set_arn = ps.arn;
```
