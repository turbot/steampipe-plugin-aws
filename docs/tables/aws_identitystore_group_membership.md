---
title: "Steampipe Table: aws_identitystore_group_membership - Query AWS Identity Store Group Memberships using SQL"
description: "Allows users to query AWS Identity Store Group Memberships, providing information about AWS users' membership status within various identity groups."
folder: "Identity Store"
---

# Table: aws_identitystore_group_membership - Query AWS Identity Store Group Memberships using SQL

The AWS Identity Store Group Membership is a part of AWS Identity Store service, which is used to manage identities in AWS. It provides a unified view of users and groups, allowing you to manage access permissions across AWS organizations. This service helps in ensuring secure access to AWS resources and data by managing group memberships effectively.

## Table Usage Guide

The `aws_identitystore_group_membership` table in Steampipe provides you with information about your AWS users' membership status within various identity groups. You can use this table to query group membership-specific details, such as the group name, user's ARN, and membership type. This table allows you to gather insights on group memberships, such as which users belong to which groups, the types of memberships they hold, and more. The schema outlines the various attributes of the group membership, including the group name, user's ARN, and membership type.

**Important Notes**
- You must specify an Identity Store ID in a `where` clause (`where identity_store_id='d-1234567890'`). You can optionally pass `group_id` in the where clause.

## Examples

### Basic info
Analyze the settings to understand the membership details within a specific AWS Identity Store. This can be useful to manage and review user access and permissions within your organization.

```sql+postgres
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

```sql+sqlite
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
This query is useful for identifying the users associated with each group membership in a specific identity store. It is particularly beneficial for managing user access rights and understanding the composition of each group within your AWS Identity Store.

```sql+postgres
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

```sql+sqlite
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
Explore which memberships are associated with specific groups within a particular identity store, providing a comprehensive overview of group details for each membership. This can be useful for assessing the organization and management of group memberships in a large-scale identity management context.

```sql+postgres
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

```sql+sqlite
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