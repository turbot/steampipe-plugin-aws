---
title: "Table: aws_identitystore_group_membership - Query AWS Identity Store Group Memberships using SQL"
description: "Allows users to query AWS Identity Store Group Memberships, providing information about AWS users' membership status within various identity groups."
---

# Table: aws_identitystore_group_membership - Query AWS Identity Store Group Memberships using SQL

The `aws_identitystore_group_membership` table in Steampipe provides information about AWS users' membership status within various identity groups. This table allows DevOps engineers to query group membership-specific details, such as the group name, user's ARN, and membership type. Users can utilize this table to gather insights on group memberships, such as which users belong to which groups, the types of memberships they hold, and more. The schema outlines the various attributes of the group membership, including the group name, user's ARN, and membership type.

**You must specify an Identity Store ID** in a `where` clause (`where identity_store_id='d-1234567890'`), `group_id` can be pass optionally in where clause.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_identitystore_group_membership` table, you can use the `.inspect aws_identitystore_group_membership` command in Steampipe.

**Key columns**:

- `group_id`: This is the unique ID of the group. It can be used to join this table with other tables that contain information about groups.
- `user_id`: This is the unique ID of the user. It can be used to join this table with other tables that contain information about users.
- `membership_type`: This column indicates the type of membership the user has within the group. It can be used to filter or sort the data based on the type of membership.

## Examples

### Basic info

```sql
select
  identity_store_id,
  group_id,
  membership_id,
  member_id
from
  aws_identitystore_group_membership
where 
  identity_store_id = 'd-1234567890';
```

### Get user details of each group membership

```sql
select
  m.membership_id,
  m.group_id,
  m.identity_store_id,
  u.name as user_name 
from
  aws_identitystore_group_membership as m,
  aws_identitystore_user as u 
where
  m.identity_store_id = 'd-1234567890' 
  and u.identity_store_id = m.identity_store_id 
  and u.id = m.member_id;
```

### Get group details of each membership

```sql
select
  m.membership_id,
  m.group_id,
  m.identity_store_id,
  g.name as group_name
from
  aws_identitystore_group_membership as m,
  aws_identitystore_group as g
where
  m.identity_store_id = 'd-1234567890'
  and g.identity_store_id = m.identity_store_id
  and g.id = m.group_id;
```
