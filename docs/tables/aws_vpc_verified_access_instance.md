---
title: "Steampipe Table: aws_vpc_verified_access_instance - Query AWS VPC Verified Access Instances using SQL"
description: "Allows users to query AWS VPC Verified Access Instances and provides information about the Amazon VPC verified access instances. This table can be used to gather details such as the instance ID, instance state, instance type, and associated tags."
folder: "VPC"
---

# Table: aws_vpc_verified_access_instance - Query AWS VPC Verified Access Instances using SQL

The AWS VPC Verified Access Instances are a part of Amazon's Virtual Private Cloud (VPC) service, allowing users to launch AWS resources into a virtual network that they define. This service provides advanced security features, such as security groups and network access control lists, to enable inbound and outbound filtering at the instance level and subnet level. In addition, you can create a Hardware VPN connection between your corporate datacenter and your VPC and leverage the AWS cloud as an extension of your corporate datacenter.

## Table Usage Guide

The `aws_vpc_verified_access_instance` table in Steampipe provides you with information about the Amazon VPC verified access instances. This table allows you, as a network administrator or security analyst, to query instance-specific details, including instance ID, instance state, instance type, and associated tags. You can utilize this table to gather insights on instances, such as instance state, type, and associated tags. The schema outlines the various attributes of the VPC verified access instance for you, including the instance ID, instance state, instance type, and associated tags.

## Examples

### Basic info
Explore which AWS VPC instances have been verified and gain insights into their creation and last updated times. This is useful for understanding your verified environments and maintaining security compliance.

```sql+postgres
select
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time,
  verified_access_trust_providers
from
  aws_vpc_verified_access_instance;
```

```sql+sqlite
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
Determine the instances in your virtual private cloud that have had verified access for over 30 days. This can be beneficial for auditing purposes, allowing you to identify potential security risks or unused resources.

```sql+postgres
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

```sql+sqlite
select
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_instance
where
  creation_time <= datetime('now', '-30 day');
```

### Get trusted provider details for each instance
Determine the areas in which each instance has a trusted provider by analyzing the provider's description, type, and associated policy. This query is useful for understanding the security measures in place for each instance and helps in managing access control.

```sql+postgres
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

```sql+sqlite
select
  i.verified_access_instance_id,
  i.creation_time,
  json_extract(p.value, '$.Description') as trust_provider_description,
  json_extract(p.value, '$.TrustProviderType') as trust_provider_type,
  json_extract(p.value, '$.UserTrustProviderType') as user_trust_provider_type,
  json_extract(p.value, '$.DeviceTrustProviderType') as device_trust_provider_type,
  json_extract(p.value, '$.VerifiedAccessTrustProviderId') as verified_access_trust_provider_id,
  t.policy_reference_name as trust_access_policy_reference_name
from
  aws_vpc_verified_access_instance as i,
  aws_vpc_verified_access_trust_provider as t,
  json_each(verified_access_trust_providers) as p
where
  json_extract(p.value, '$.VerifiedAccessTrustProviderId') = t.verified_access_trust_provider_id;
```