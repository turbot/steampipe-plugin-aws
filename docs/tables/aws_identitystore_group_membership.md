# Table: aws_identitystore_group_membership

AWS Identity Store Group Membership represents the membership of a user in a group. It is used to manage and track the association between users and groups within the identity store.

By managing group memberships, you can control access and permissions for users within your organization and assign them to specific groups to define their level of access to resources and services.

**You must specify an Identity Store ID** in a `where` clause (`where identity_store_id='d-1234567890'`), `group_id` can be pass optionally in where clause.

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
