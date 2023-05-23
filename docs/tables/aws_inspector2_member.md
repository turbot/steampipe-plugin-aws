# Table: aws_inspector2_finding

AWS Inspector Member refers to a user account or IAM role that has been granted permission to access and interact with the AWS Inspector service within an AWS organization.

As part of the AWS Organizations feature, you can create an organization that includes multiple AWS accounts. Within this organization, you can designate specific accounts as "member accounts." These member accounts can then be granted access to various AWS services, including AWS Inspector.

## Examples

### Basic info

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member;
```

### Retrieve a list of members whose status remains unchanged in the past 30 days

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status,
  updated_at
from
  aws_inspector2_member
where
  updated_at >= now() - interval '30' day;
```

### List invited members

```sql
select
  member_account_id,
  delegated_admin_account_id,
  relationship_status
from
  aws_inspector2_member
where
  relationship_status = 'INVITED';
```