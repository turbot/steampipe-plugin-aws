---
title: "Table: aws_ec2_key_pair - Query AWS EC2 Key Pairs using SQL"
description: "Allows users to query AWS EC2 Key Pairs, providing information about key pairs which are used to securely log into EC2 instances."
---

# Table: aws_ec2_key_pair - Query AWS EC2 Key Pairs using SQL

The `aws_ec2_key_pair` table in Steampipe provides information about Key Pairs within AWS EC2 (Elastic Compute Cloud). This table allows DevOps engineers, security teams, and system administrators to query key pair-specific details, including key fingerprints, key material, and associated tags. Users can utilize this table to gather insights on key pairs, such as verifying key fingerprints, checking the existence of specific key pairs, and more. The schema outlines the various attributes of the EC2 key pair, including the key pair name, key pair ID, key type, public key, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_key_pair` table, you can use the `.inspect aws_ec2_key_pair` command in Steampipe.

### Key columns:

- `key_pair_id`: The ID of the key pair. This column can be used to join this table with other tables that contain key pair ID information.
- `key_name`: The name of the key pair. This is an important column as it is commonly used to identify key pairs in AWS.
- `key_type`: The type of key pair, either `rsa` or `ecdsa`. This column is useful to understand the encryption method used by the key pair.

## Examples

### Basic keypair info

```sql
select
  key_name,
  key_pair_id,
  region
from
  aws_ec2_key_pair;
```


### List of keypairs without owner tag key

```sql
select
  key_name,
  tags
from
  aws_ec2_key_pair
where
  not tags :: JSONB ? 'owner';
```