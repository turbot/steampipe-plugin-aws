---
title: "Table: aws_vpc_verified_access_instance - Query AWS VPC Verified Access Instances using SQL"
description: "Allows users to query AWS VPC Verified Access Instances and provides information about the Amazon VPC verified access instances. This table can be used to gather details such as the instance ID, instance state, instance type, and associated tags."
---

# Table: aws_vpc_verified_access_instance - Query AWS VPC Verified Access Instances using SQL

The `aws_vpc_verified_access_instance` table in Steampipe provides information about the Amazon VPC verified access instances. This table allows network administrators and security analysts to query instance-specific details, including instance ID, instance state, instance type, and associated tags. Users can utilize this table to gather insights on instances, such as instance state, type, and associated tags. The schema outlines the various attributes of the VPC verified access instance, including the instance ID, instance state, instance type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_verified_access_instance` table, you can use the `.inspect aws_vpc_verified_access_instance` command in Steampipe.

### Key columns:

- `instance_id`: The unique identifier for the instance. This can be used to join with other tables that contain instance-specific information.
- `instance_state`: The current state of the instance. This is useful for identifying instances that are inactive or in an error state.    
- `instance_type`: The type of the instance. This is useful for identifying the capacity and performance characteristics of the instance.

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

### List VPC access verified instances older than 30 days

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
