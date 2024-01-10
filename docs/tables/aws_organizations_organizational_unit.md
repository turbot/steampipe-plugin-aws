# Table: aws_organizations_organizational_unit

A container for accounts within a root. An OU also can contain other OUs, enabling you to create a hierarchy that resembles an upside-down tree, with a root at the top and branches of OUs that reach down, ending in accounts that are the leaves of the tree. When you attach a policy to one of the nodes in the hierarchy, it flows down and affects all the branches (OUs) and leaves (accounts) beneath it. An OU can have exactly one parent, and currently each account can be a member of exactly one OU.

To represent the hierarchical structure, the table includes a `path` column. This column is crucial for understanding the relationship between different OUs in the hierarchy. Due to compatibility issues with the ltree type, which is typically used for representing tree-like structures in PostgreSQL, the standard hyphen (-) in the path values has been replaced with an underscore (\_). This modification ensures proper functionality of the ltree operations and queries.

By default, querying the table without any specific filters will return all OUs from the root of the hierarchy. Users have the option to query the table using a specific `parent_id`. This allows for the retrieval of all direct child OUs under the specified parent.

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

### Find a specific organizational unit and all its descendants

```sql
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path <@ 'r_wxnb.ou_wxnb_m8l8tpaq';
```

### Select all organizational units at a certain level in the hierarchy

```sql
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  nlevel(path) = 3;
```

### Get all ancestors of a given organizational unit

```sql
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  'r_wxnb.ou_wxnb_m8l8tpaq.ou_wxnb_5grilgkb' @> path;
```

### Retrieve all siblings of a specific organizational unit

```sql
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  parent_id =
  (
    select
      parent_id
    from
      aws_organizations_organizational_unit
    where
      name = 'Punisher'
  );
```

### Select organizational units with a path that matches a specific pattern

```sql
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path ~ 'r_wxnb.*.ou_wxnb_m8l8tpaq.*';
```
