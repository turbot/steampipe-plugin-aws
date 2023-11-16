---
title: "Table: aws_ssoadmin_account_assignment - Query AWS SSO Admin Account Assignment using SQL"
description: "Allows users to query AWS SSO Admin Account Assignments. This table provides information about each AWS SSO admin account assignment within an AWS account."
---

# Table: aws_ssoadmin_account_assignment - Query AWS SSO Admin Account Assignment using SQL

The `aws_ssoadmin_account_assignment` table in Steampipe provides information about each AWS SSO (Single Sign-On) admin account assignment within an AWS account. This table allows DevOps engineers, administrators, and AWS users to query details related to SSO admin account assignments, including the principal type, principal ID, target type, target ID, and permission set. Users can utilize this table to gather insights on SSO admin account assignments, such as the account assignments for a specific principal or target, the permission sets assigned to a target, and more. The schema outlines the various attributes of the SSO admin account assignment, including the instance ARN, principal type, principal ID, target type, target ID, and permission set.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssoadmin_account_assignment` table, you can use the `.inspect aws_ssoadmin_account_assignment` command in Steampipe.

**Key columns**:

- `instance_arn`: The ARN of the SSO instance under which the assignment was made. This can be useful for joining with other tables that also contain SSO instance ARNs.
- `principal_id`: The identifier of the principal (user or group) that the assignment applies to. This can be used to join with other tables that contain principal identifiers.
- `target_id`: The identifier of the target (AWS account or AWS Organizations unit) that the assignment applies to. This can be used to join with other tables that contain target identifiers.

## Examples

### Assignments for a specific permission set and account

```sql
select
  permission_set_arn,
  target_account_id,
  principal_type,
  principal_id
from
  aws_ssoadmin_account_assignment
where
  permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef'
  and target_account_id = '012347678910';
```

### Assignments for a specific permission set and account, with user/group information from Identity Store

```sql
with aws_ssoadmin_principal as
(
  select
    i.arn as instance_arn,
    'GROUP' as "type",
    g.id,
    g.title
  from
    aws_ssoadmin_instance i
    left join
      aws_identitystore_group g
      on i.identity_store_id = g.identity_store_id
    union
    select
      i.arn as instance_arn,
      'USER' as "type",
      u.id,
      u.title
    from
      aws_ssoadmin_instance i
      left join
        aws_identitystore_user u
        on i.identity_store_id = u.identity_store_id
)
select
  a.target_account_id,
  a.principal_type,
  p.title as principal_title
from
  aws_ssoadmin_account_assignment a
  left join
    aws_ssoadmin_principal p
    on a.principal_type = p.type
    and a.principal_id = p.id
    and a.instance_arn = p.instance_arn
where
  a.target_account_id = '012345678901' a.permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef';
```
