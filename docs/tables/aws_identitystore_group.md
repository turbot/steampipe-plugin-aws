---
title: "Table: aws_identitystore_group - Query AWS Identity Store Groups using SQL"
description: "Allows users to query AWS Identity Store Groups to obtain information about the identity and attributes of groups in AWS."
---

# Table: aws_identitystore_group - Query AWS Identity Store Groups using SQL

The `aws_identitystore_group` table in Steampipe provides information about groups within AWS Identity Store. This table allows DevOps engineers to query group-specific details, including group ID, group name, and associated metadata. Users can utilize this table to gather insights on groups, such as group names, verification of group identities, and more. The schema outlines the various attributes of the AWS Identity Store group, including the group ID, group name, and display name.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_identitystore_group` table, you can use the `.inspect aws_identitystore_group` command in Steampipe.

### Key columns:

- `group_id`: This is the unique identifier for the group. It can be used to join this table with other tables that require a group ID for querying specific group information.
- `group_name`: This is the name of the group. It is useful for filtering queries based on the group name.
- `display_name`: This is the display name of the group. It can be used for user-friendly querying and reporting.

## Examples

### Get group by ID

```sql
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

### List groups by name

```sql
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and name = 'test';
```
