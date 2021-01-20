# Table: aws_iam_role

An IAM role is an AWS Identity and Access Management (IAM) entity with permissions to make AWS service requests.

## Examples

### List of IAM roles with no inline policy

```sql
select
  name,
  create_date
from
  aws_iam_role
where
  inline_policies is null;
```


### List the policies attached to the roles

```sql
select
  name,
  description,
  split_part(policy, '/', 3) as attached_policy
from
  aws_iam_role
  cross join jsonb_array_elements_text(attached_policy_arns) as policy;
```


### Permission boundary information for each role

```sql
select
  name,
  description,
  permissions_boundary_arn,
  permissions_boundary_type
from
  aws_iam_role;
```