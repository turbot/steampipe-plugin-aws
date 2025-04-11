---
title: "Steampipe Table: aws_identitystore_user - Query AWS Identity Store User using SQL"
description: "Allows users to query AWS Identity Store User data, providing details such as user ID, username, and ARN. This table is essential for managing and auditing user information within the AWS Identity Store."
folder: "Identity Store"
---

# Table: aws_identitystore_user - Query AWS Identity Store User using SQL

The AWS Identity Store User is a resource in AWS Identity Store that allows you to manage user identities. It provides a unified view of users and groups across AWS SSO and AWS Managed Microsoft AD, helping you to simplify identity management. It also enables you to perform identity-based actions in your AWS environment, enhancing the security and governance of your resources.

## Table Usage Guide

The `aws_identitystore_user` table in Steampipe provides you with information about users within the AWS Identity Store. This table allows you, as a DevOps engineer, to query user-specific details, including user ID, username, and ARN. You can utilize this table to manage and audit user information, such as user identities, associated roles, and permissions. The schema outlines the various attributes of the user for you, including the user ID, ARN, username, and status.

## Examples

### Get user by ID
Explore which user is associated with a specific ID in the AWS Identity Store. This is useful to validate user identities and ensure appropriate access controls are in place.

```sql+postgres
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

```sql+sqlite
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

### List users by name
Determine the areas in which specific users are identified within a particular identity store. This is useful for pinpointing the presence and details of specific users within a given identity store, to manage and track user data.

```sql+postgres
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and name = 'test';
```

```sql+sqlite
select
  id,
  name
from
  aws_identitystore_user
where identity_store_id = 'd-1234567890' and name = 'test';
```