---
title: "Steampipe Table: aws_identitystore_group - Query AWS Identity Store Groups using SQL"
description: "Allows users to query AWS Identity Store Groups to obtain information about the identity and attributes of groups in AWS."
folder: "Identity Store"
---

# Table: aws_identitystore_group - Query AWS Identity Store Groups using SQL

The AWS Identity Store service provides information about identities in your AWS organization. It enables you to retrieve information about groups, including group name, group ID, and the AWS SSO instance that the group belongs to. This service helps ensure that your applications have access to the identity information they need while adhering to privacy best practices.

## Table Usage Guide

The `aws_identitystore_group` table in Steampipe provides you with information about groups within AWS Identity Store. This table allows you, as a DevOps engineer, to query group-specific details, including group ID, group name, and associated metadata. You can utilize this table to gather insights on groups, such as group names, verification of group identities, and more. The schema outlines the various attributes of the AWS Identity Store group for you, including the group ID, group name, and display name.

## Examples

### Get group by ID
Determine the specific group within AWS Identity Store using a unique identifier. This can be useful for administrators needing to manage or monitor a particular group's settings or activity.

```sql+postgres
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

```sql+sqlite
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

### List groups by name
Determine the areas in which specific user groups are identified within a particular identity store in AWS. This is useful for managing access controls and permissions in a secure environment.

```sql+postgres
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and name = 'test';
```

```sql+sqlite
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and name = 'test';
```