# Table: aws_vpc_verified_access_instance

A Verified Access instance is an AWS resource that helps you organize your trust providers and Verified Access groups.

## Examples

### Basic info

```sql
select
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time,
  verified_access_trust_providers
from
  aws_vpc_verified_access_instance;
```

### List trusted providers older than 30 days

```sql
select
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_instance
where
  creation_time <= now() - interval '30' day;
```

### Get trusted provider details for each instance

```sql
select
  i.verified_access_instance_id,
  i.creation_time,
  p ->> 'Description' as trust_provider_description,
  p ->> 'TrustProviderType' as trust_provider_type,
  p ->> 'UserTrustProviderType' as user_trust_provider_type,
  p ->> 'DeviceTrustProviderType' as device_trust_provider_type,
  p ->> 'VerifiedAccessTrustProviderId' as verified_access_trust_provider_id,
  t.policy_reference_name as trust_access_policy_reference_name
from
  aws_vpc_verified_access_instance as i,
  aws_vpc_verified_access_trust_provider as t,
  jsonb_array_elements(verified_access_trust_providers) as p
where
  p ->> 'VerifiedAccessTrustProviderId' = t.verified_access_trust_provider_id;
```
