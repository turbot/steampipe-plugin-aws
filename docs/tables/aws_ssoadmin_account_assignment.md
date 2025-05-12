---
title: "Steampipe Table: aws_ssoadmin_account_assignment - Query AWS SSO Admin Account Assignment using SQL"
description: "Allows users to query AWS SSO Admin Account Assignments. This table provides information about each AWS SSO admin account assignment within an AWS account."
folder: "SSO"
---

# Table: aws_ssoadmin_account_assignment - Query AWS SSO Admin Account Assignment using SQL

The AWS SSO Admin Account Assignment is a resource within AWS Single Sign-On (SSO) service that allows you to manage the assignment of access permissions to users. It enables the administrators to assign user access to AWS accounts, SSO instances, and permission sets using SQL queries. This resource plays a crucial role in managing and controlling access to AWS resources and services, enhancing the security and governance of your AWS environment.

## Table Usage Guide

The `aws_ssoadmin_account_assignment` table in Steampipe provides you with information about each AWS SSO (Single Sign-On) admin account assignment within your AWS account. This table allows you, as a DevOps engineer, administrator, or AWS user, to query details related to SSO admin account assignments, including the principal type, principal ID, target type, target ID, and permission set. You can utilize this table to gather insights on SSO admin account assignments, such as the account assignments for a specific principal or target, the permission sets assigned to a target, and more. The schema outlines the various attributes of the SSO admin account assignment for you, including the instance ARN, principal type, principal ID, target type, target ID, and permission set.

## Examples

### Assignments for a specific permission set and account
Determine the areas in which specific permissions are assigned within a particular account. This can provide insights into how access control is structured within your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Explore which user or group has been assigned a specific permission set in an AWS account. This is useful for understanding access controls and managing permissions within your organization.

```sql+postgres
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
  a.target_account_id = '012345678901' and a.permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef';
```

```sql+sqlite
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
    and a.instance_arn = p.instance_arn
where
  a.target_account_id = '012345678901' and a.permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef';
```