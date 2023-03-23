# Table: aws_ssoadmin_account_assignment

Contains information about AWS SSO account assignments (an assignment of a SSO permission set to an Identity Store User/Group in an Account).

## Examples

### Assignments for a specific permission set and account

```sql
select
    permission_set_arn,
    target_account_id,
    identity_store_id,
    principal_type,
    principal_id,
from
    aws_ssoadmin_account_assignment
where
    permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef'
    and target_account_id = '012347678910'
```

### Assignments for a specific permission set and account, with user/group information from Identity Store

```sql
with aws_ssoadmin_principal as (
    select
        i.arn   as instance_arn,
        'GROUP' as "type",
        g.id,
        g.title
    from
       aws_ssoadmin_instance i
    left join
       aws_identitystore_group g on i.identity_store_id = g.identity_store_id
    union
    select
       i.arn  as instance_arn,
        'USER' as "type",
        u.id,
        u.title
    from
       aws_ssoadmin_instance i
    left join
       aws_identitystore_user u on i.identity_store_id = u.identity_store_id)
select
    a.target_account_id,
    a.principal_type,
    p.title as principal_title
from
    aws_ssoadmin_account_assignment a
left join
    aws_ssoadmin_principal p on a.principal_type = p.type and a.principal_id = p.id and a.instance_arn = p.instance_arn
where
    a.target_account_id = '012345678901'
    a.permission_set_arn = 'arn:aws:sso:::permissionSet/ssoins-0123456789abcdef/ps-0123456789abcdef'
```
