# Table: aws_organizations_organizational_unit

A container for accounts within a root. An OU also can contain other OUs, enabling you to create a hierarchy that resembles an upside-down tree, with a root at the top and branches of OUs that reach down, ending in accounts that are the leaves of the tree. When you attach a policy to one of the nodes in the hierarchy, it flows down and affects all the branches (OUs) and leaves (accounts) beneath it. An OU can have exactly one parent, and currently each account can be a member of exactly one OU.

You **_must_** specify a single `parent_id` or `id` in a where or join clause in order to use this table.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  parent_id,
  title,
  akas
from
  aws_organizations_organizational_unit;
```