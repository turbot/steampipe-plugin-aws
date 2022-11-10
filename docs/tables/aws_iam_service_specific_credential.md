# Table: aws_iam_service_specific_credential

Service-specific credentials are associated with a specific IAM user and can only be used for the service they were created for. To give IAM roles or federated identities permissions to access all your AWS resources, you should create IAM access keys for AWS authentication and use the SigV4 authentication plugin.

## Examples

### Basic service specific credential info

```sql
select
  service_name,
  service_specific_credential_id,
  create_date,
  user_name
from
  aws_iam_service_specific_credential;
```

### IAM user details for the service specific credential

```sql
select
  s.service_name as service_name,
  s.service_specific_credential_id as service_specific_credential_id,
  u.name as user_name,
  u.user_id as user_id,
  u.password_last_used as password_last_used,
  u.mfa_enabled as mfa_enabled
  iam_group ->> 'GroupId' as group_id,
  iam_group ->> 'CreateDate' as create_date
from
  aws_iam_service_specific_credential as s,
  aws_iam_user as u
where
  s.user_name = u.name;
```