---
title: "Table: aws_vpc_verified_access_trust_provider - Query AWS VPC Verified Access Trust Providers using SQL"
description: "Allows users to query AWS VPC Verified Access Trust Providers, providing information about the trust providers for VPC endpoints in AWS. This table can be used to gain insights into the trust relationships between VPC endpoints and the services they access."
---

# Table: aws_vpc_verified_access_trust_provider - Query AWS VPC Verified Access Trust Providers using SQL

The `aws_vpc_verified_access_trust_provider` table in Steampipe provides information about the trust providers for VPC endpoints within AWS Virtual Private Cloud (VPC). This table allows DevOps engineers to query trust provider-specific details, including the provider type, owner, and associated metadata. Users can utilize this table to gather insights on trust relationships, such as the services that VPC endpoints can access, verification of trust providers, and more. The schema outlines the various attributes of the trust provider, including the provider type, owner, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_verified_access_trust_provider` table, you can use the `.inspect aws_vpc_verified_access_trust_provider` command in Steampipe.

**Key columns**:

- `provider_type`: This column provides the type of the trust provider. It can be useful for understanding the nature of the trust relationship.
- `owner`: This column provides the owner of the trust provider. It can be useful for identifying the entity responsible for the trust provider.
- `tags`: This column provides any associated tags with the trust provider. It can be useful for categorizing and managing trust providers.

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

### List trusted providers that are of the user type

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
  creation_time >= now() - interval '90' day;
```
