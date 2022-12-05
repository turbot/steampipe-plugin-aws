# Table: aws_vpc_verified_access_trust_provider

A trust provider is a service that sends information about users and devices, called trust data, to AWS Verified Access. Trust data may include attributes based on user identity such as an email address or membership in the "sales" organization, or device management information such as security patches or antivirus software version.

## Examples

### Basic info

```sql
select
  verified_access_trust_provider_id,
  creation_time,
  device_trust_provider_type,
  last_updated_time,
  policy_reference_name,
  trust_provider_type
from
  aws_vpc_verified_access_trust_provider;
```

### List trusted providers that are of type user

```sql
select
  verified_access_trust_provider_id,
  creation_time,
  device_trust_provider_type,
  last_updated_time,
  policy_reference_name,
  trust_provider_type
from
  aws_vpc_verified_access_trust_provider
where
  trust_provider_type = 'user';
```

### List trusted providers older than 90 days

```sql
select
  verified_access_trust_provider_id,
  creation_time,
  last_updated_time,
  policy_reference_name,
  trust_provider_type
from
  aws_vpc_verified_access_trust_provider
where
  creation_time >= now() - interval '90' days;
```
