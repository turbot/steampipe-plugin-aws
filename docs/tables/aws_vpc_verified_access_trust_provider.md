---
title: "Steampipe Table: aws_vpc_verified_access_trust_provider - Query AWS VPC Verified Access Trust Providers using SQL"
description: "Allows users to query AWS VPC Verified Access Trust Providers, providing information about the trust providers for VPC endpoints in AWS. This table can be used to gain insights into the trust relationships between VPC endpoints and the services they access."
folder: "VPC"
---

# Table: aws_vpc_verified_access_trust_provider - Query AWS VPC Verified Access Trust Providers using SQL

The AWS VPC Verified Access Trust Provider is an AWS service that helps manage and verify access to your Virtual Private Cloud (VPC). This service allows you to control and secure network access to your AWS resources within the VPC. It provides a layer of security that helps you control who can access your resources within a VPC from the internet.

## Table Usage Guide

The `aws_vpc_verified_access_trust_provider` table in Steampipe provides you with information about the trust providers for VPC endpoints within AWS Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query trust provider-specific details, including the provider type, owner, and associated metadata. You can utilize this table to gather insights on trust relationships, such as the services that VPC endpoints can access, verification of trust providers, and more. The schema outlines the various attributes of the trust provider for you, including the provider type, owner, and associated tags.

## Examples

### Basic info
Explore the creation and update timeline of verified access trust providers in your AWS VPC. This can help in maintaining security by identifying the type of trust providers and understanding their policy references.

```sql+postgres
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

```sql+sqlite
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

### List trusted providers that are of the user type
Explore which trusted providers are specifically categorized as 'user' type within your AWS VPC. This can help in managing access controls and understanding the security posture of your virtual private cloud environment.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are trusted providers and have been active for over 90 days. This can be useful for assessing the longevity and reliability of these providers in your AWS VPC environment.

```sql+postgres
select
  verified_access_trust_provider_id,
  creation_time,
  last_updated_time,
  policy_reference_name,
  trust_provider_type
from
  aws_vpc_verified_access_trust_provider
where
  creation_time >= now() - interval '90' day;
```

```sql+sqlite
select
  verified_access_trust_provider_id,
  creation_time,
  last_updated_time,
  policy_reference_name,
  trust_provider_type
from
  aws_vpc_verified_access_trust_provider
where
  creation_time >= datetime('now', '-90 day');
```