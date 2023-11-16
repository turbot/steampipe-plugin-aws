---
title: "Table: aws_networkfirewall_firewall - Query AWS Network Firewall using SQL"
description: "Allows users to query AWS Network Firewalls for detailed information about each firewall's properties and settings."
---

# Table: aws_networkfirewall_firewall - Query AWS Network Firewall using SQL

The `aws_networkfirewall_firewall` table in Steampipe provides information about each firewall in AWS Network Firewall. This table allows network administrators, security analysts, and DevOps engineers to query specific details about firewalls, including firewall policies, subnet mappings, and associated VPCs. Users can utilize this table to gain insights into firewall configurations, such as firewall policy ARNs, VPC IDs, subnet IDs, and more. The schema outlines the various attributes of the firewall, including the firewall ARN, firewall name, firewall policy ARN, VPC ID, subnet mapping, delete protection status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_networkfirewall_firewall` table, you can use the `.inspect aws_networkfirewall_firewall` command in Steampipe.

### Key columns:

- `firewall_arn`: The Amazon Resource Number (ARN) of the firewall. This is a unique identifier and can be used to join this table with other tables that contain firewall ARNs.
- `firewall_name`: The name of the firewall. This can be useful for joining with other tables that reference firewalls by name.
- `vpc_id`: The ID of the VPC where the firewall is deployed. This column can be used to join this table with other tables that contain VPC IDs, providing context about the network in which the firewall operates.

## Examples

### Basic info

```sql
select
  arn,
  name,
  description,
  vpc_id,
  policy_arn,
  region,
  tags
from
  aws_networkfirewall_firewall;
```

### List firewalls using default encryption

```sql
select
  arn,
  name,
  description,
  encryption_configuration
from
  aws_networkfirewall_firewall
where 
  encryption_configuration ->> 'Type' = `aws_OWNED_KMS_KEY';
```

### List firewalls having deletion protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not delete_protection;
```

### List firewalls having policy change protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not policy_change_protection;
```

### List firewalls having subnet change protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not subnet_change_protection;
```

### Get subnet details for each firewall

```sql
select
  f.arn,
  f.name,
  f.vpc_id,
  s ->> 'SubnetId' as subnet_id,
  cidr_block,
  availability_zone,
  default_for_az
from
  aws_networkfirewall_firewall f,
  jsonb_array_elements(subnet_mappings) s,
  aws_vpc_subnet vs
where
  vs.subnet_id = s ->> 'SubnetId';
```

### Get KMS key details of firewalls encrypted with customer managed keys

```sql
select
  f.arn,
  f.name,
  f.vpc_id,
  k.arn as key_arn,
  key_rotation_enabled
from
  aws_networkfirewall_firewall f,
  aws_kms_key k
where
  k.id = encryption_configuration ->> 'KeyId'
  and not f.encryption_configuration ->> 'Type' = `aws_OWNED_KMS_KEY';
```

