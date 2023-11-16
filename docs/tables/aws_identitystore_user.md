---
title: "Table: aws_identitystore_user - Query AWS Identity Store User using SQL"
description: "Allows users to query AWS Identity Store User data, providing details such as user ID, username, and ARN. This table is essential for managing and auditing user information within the AWS Identity Store."
---

# Table: aws_identitystore_user - Query AWS Identity Store User using SQL

The `aws_identitystore_user` table in Steampipe provides information about users within the AWS Identity Store. This table allows DevOps engineers to query user-specific details, including user ID, username, and ARN. Users can utilize this table to manage and audit user information, such as user identities, associated roles, and permissions. The schema outlines the various attributes of the user, including the user ID, ARN, username, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_identitystore_user` table, you can use the `.inspect aws_identitystore_user` command in Steampipe.

### Key columns:

- `user_id`: This is the unique identifier for the user. It can be used to join with other tables that require a user ID for specific queries.
- `username`: This is the name associated with the user. This can be useful when joining with tables that use username for identification.
- `arn`: This is the Amazon Resource Name (ARN) for the user. This can be useful when joining with other tables that use ARN for identification and resource linking.

## Examples

### Get user by ID

```sql
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

### List users by name

```sql
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and name = 'test';
```
