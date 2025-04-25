---
title: "Steampipe Table: aws_ec2_regional_settings - Query AWS EC2 Regional Settings using SQL"
description: "Allows users to query AWS EC2 regional settings, including default EBS encryption and default EBS encryption KMS key."
folder: "EC2"
---

# Table: aws_ec2_regional_settings - Query AWS EC2 Regional Settings using SQL

The AWS EC2 Regional Settings are configurations that apply to an entire region in the Amazon Elastic Compute Cloud (EC2) service. These settings include options such as default VPC, default subnet, and default security group. They allow for the customization and management of resources within a specific AWS region.

## Table Usage Guide

The `aws_ec2_regional_settings` table in Steampipe provides you with information about the regional settings of Amazon Elastic Compute Cloud (EC2). This table allows you, as a cloud administrator, security team member, or developer, to query regional settings, including default EBS encryption and the default EBS encryption KMS key. You can utilize this table to gather insights on regional settings, such as the default EBS encryption status, the default EBS encryption KMS key, and the region name. The schema outlines the various attributes of the regional settings for you, including the region, default EBS encryption, and default EBS encryption KMS key.

## Examples

### Basic settings info
Analyze the settings to understand the default encryption status and key for your AWS EC2 regional settings. This is useful for ensuring your data is secure and encrypted as per your organization's policies.

```sql+postgres
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_regional_settings;
```

```sql+sqlite
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_regional_settings;
```


### Settings info for a particular region
Determine the areas in which default encryption is enabled for a specific region. This query is beneficial for understanding the security configuration of your cloud storage in that particular region.

```sql+postgres
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

```sql+sqlite
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
Identify regions where the default EBS encryption is enabled. This is useful for maintaining data security and compliance by ensuring that encrypted storage is being utilized in those areas.

```sql+postgres
select
  region,
  default_ebs_encryption_enabled,
  default_ebs_encryption_key
from
  aws_ec2_regional_settings
where
  default_ebs_encryption_enabled;
```

```sql+sqlite
select
  region,
  default_ebs_encryption_enabled,
  default_ebs_encryption_key
from
  aws_ec2_regional_settings
where
  default_ebs_encryption_enabled = 1;
```