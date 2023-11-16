---
title: "Table: aws_ec2_regional_settings - Query AWS EC2 Regional Settings using SQL"
description: "Allows users to query AWS EC2 regional settings, including default EBS encryption and default EBS encryption KMS key."
---

# Table: aws_ec2_regional_settings - Query AWS EC2 Regional Settings using SQL

The `aws_ec2_regional_settings` table in Steampipe provides information about the regional settings of Amazon Elastic Compute Cloud (EC2). This table allows cloud administrators, security teams, and developers to query regional settings, including default EBS encryption and the default EBS encryption KMS key. Users can utilize this table to gather insights on regional settings, such as the default EBS encryption status, the default EBS encryption KMS key, and the region name. The schema outlines the various attributes of the regional settings, including the region, default EBS encryption, and default EBS encryption KMS key.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_regional_settings` table, you can use the `.inspect aws_ec2_regional_settings` command in Steampipe.

Key columns:

- `region`: This column provides the name of the region. It is useful for joining with other tables that contain regional data.
- `default_ebs_encryption`: This column indicates whether EBS encryption is enabled by default. It can be used to join with other tables that track encryption settings.
- `default_ebs_encryption_kms_key_id`: This column provides the ID of the KMS key used for EBS encryption. It is useful for joining with other tables that contain KMS key data.

## Examples

### Basic settings info

```sql
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_regional_settings;
```


### Settings info for a particular region

```sql
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_regional_settings
where
  region = 'ap-south-1';
```


### List the regions along with the key where default EBS encryption is enabled

```sql
select
  region,
  default_ebs_encryption_enabled,
  default_ebs_encryption_key
from
  aws_ec2_regional_settings
where
  default_ebs_encryption_enabled;
```
