---
title: "Steampipe Table: aws_organizations_organizational_unit - Query AWS Organizations Organizational Units using SQL"
description: "Allows users to query AWS Organizations Organizational Units and provides information about each OU."
folder: "Organizations"
---

# Table: aws_organizations_organizational_unit

A container for accounts within a root. An OU also can contain other OUs, enabling you to create a hierarchy that resembles an upside-down tree, with a root at the top and branches of OUs that reach down, ending in accounts that are the leaves of the tree. When you attach a policy to one of the nodes in the hierarchy, it flows down and affects all the branches (OUs) and leaves (accounts) beneath it. An OU can have exactly one parent, and currently each account can be a member of exactly one OU.

## Table Usage Guide

The `aws_organizations_organizational_unit` table in Steampipe provides you with information about  the hierarchical structure, the table includes a `path` column. This column is crucial for understanding the relationship between different OUs in the hierarchy. Due to compatibility issues with the `ltree` type, which is typically used for representing tree-like structures in PostgreSQL, the standard hyphen (-) in the path values has been replaced with an underscore (\_). This modification ensures proper functionality of the `ltree` operations and queries.

By default, querying the table without any specific filters will return all OUs from the root of the hierarchy. Users have the option to query the table using a specific `parent_id`. This allows for the retrieval of all direct child OUs under the specified parent.

## Examples

### Basic info
This query helps AWS administrators and cloud architects to efficiently manage, audit, and report on the structure and composition of their AWS Organizations.

```sql+postgres
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

```sql+sqlite
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
By filtering OUs based on their path, the query efficiently retrieves information about a specific subset of your organization's structure, which is particularly useful for large organizations with complex hierarchies.

```sql+postgres
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path <@ 'r_wxnb.ou_wxnb_m8l8t123';
```

```sql+sqlite
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path like 'r_wxnb.ou_wxnb_m8l8t123%'
```

### Select all organizational units at a certain level in the hierarchy
Retrieving a list of organizational units (OUs) from a structured hierarchy, specifically those that exist at a particular level. In the context of a database or a management system like AWS Organizations, this involves using a query to filter and display only the OUs that are positioned at the same depth or stage in the hierarchical structure.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  (length(path) - length(replace(path, '.', ''))) = 2;
```

### Get all ancestors of a given organizational unit
Ancestors are the units in the hierarchy that precede the given OU. An ancestor can be a direct parent (the immediate higher-level unit), or it can be any higher-level unit up to the root of the hierarchy.

```sql+postgres
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  'r_wxnb.ou_wxnb_m8l123aq.ou_wxnb_5gri123b' @> path;
```

```sql+sqlite
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path like 'r_wxnb.ou_wxnb_m8l123aq.ou_wxnb_5gri123b%';
```

### Retrieve all siblings of a specific organizational unit
The query is useful for retrieving information about sibling organizational units corresponding to a specified organizational unit.

```sql+postgres
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

```sql+sqlite
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
This query is designed to retrieve organizational units that have a specific hierarchical path pattern within an AWS (Amazon Web Services) organization's structure.

```sql+postgres
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path ~ 'r_wxnb.*.ou_wxnb_m81234aq.*';
```

```sql+sqlite
select
  name,
  id,
  parent_id,
  path
from
  aws_organizations_organizational_unit
where
  path like 'r_wxnb%ou_wxnb_m81234aq%';
```
