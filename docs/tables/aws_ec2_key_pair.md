---
title: "Steampipe Table: aws_ec2_key_pair - Query AWS EC2 Key Pairs using SQL"
description: "Allows users to query AWS EC2 Key Pairs, providing information about key pairs which are used to securely log into EC2 instances."
folder: "EC2"
---

# Table: aws_ec2_key_pair - Query AWS EC2 Key Pairs using SQL

The AWS EC2 Key Pair is a security feature utilized within Amazon's Elastic Compute Cloud (EC2). It provides a simple, secure way to log into your instances using SSH. The key pair is composed of a public key that AWS stores, and a private key file that you store, enabling an encrypted connection to your instance.

## Table Usage Guide

The `aws_ec2_key_pair` table in Steampipe provides you with information about Key Pairs within AWS EC2 (Elastic Compute Cloud). This table allows you, as a DevOps engineer, security team member, or system administrator, to query key pair-specific details, including key fingerprints, key material, and associated tags. You can utilize this table to gather insights on key pairs, such as verifying key fingerprints, checking the existence of specific key pairs, and more. The schema outlines the various attributes of the EC2 key pair for you, including the key pair name, key pair ID, key type, public key, and associated tags.

## Examples

### Basic keypair info
Analyze the settings to understand the distribution of your EC2 key pairs across various regions. This can help in managing your AWS resources efficiently and ensuring balanced utilization.

```sql+postgres
select
  key_name,
  key_pair_id,
  region
from
  aws_ec2_key_pair;
```

```sql+sqlite
select
  key_name,
  key_pair_id,
  region
from
  aws_ec2_key_pair;
```

### List of keypairs without owner tag key
Identify instances where AWS EC2 key pairs are not tagged with an owner. This is useful for maintaining efficient tag management and ensuring accountability for key pair usage.

```sql+postgres
select
  key_name,
  tags
from
  aws_ec2_key_pair
where
  not tags :: JSONB ? 'owner';
```

```sql+sqlite
select
  key_name,
  tags
from
  aws_ec2_key_pair
where
  json_extract(tags, '$.owner') IS NULL;
```