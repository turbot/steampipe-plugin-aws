---
title: "Steampipe Table: aws_ec2_ami - Query AWS EC2 AMI using SQL"
description: "Allows users to query AWS EC2 AMIs (Amazon Machine Images) to retrieve detailed information about each AMI available in the AWS account."
folder: "EC2"
---

# Table: aws_ec2_ami - Query AWS EC2 AMI using SQL

The AWS EC2 AMI (Amazon Machine Image) provides the information necessary to launch an instance, which is a virtual server in the cloud. You specify an AMI when you launch an instance, and you can launch as many instances from the AMI as you need. AMIs are designed to provide a stable, secure, and high performance execution environment for applications running on Amazon EC2.

## Table Usage Guide

The `aws_ec2_ami` table in Steampipe provides you with information about AMIs (Amazon Machine Images) within Amazon Elastic Compute Cloud (Amazon EC2). This table allows you, as a DevOps engineer, system administrator, or other technical professional, to query AMI-specific details, including its attributes, block device mappings, and associated tags. You can utilize this table to gather insights on AMIs, such as identifying unused or outdated AMIs, verifying AMI permissions, and more. The schema outlines the various attributes of the AMI for you, including the AMI ID, creation date, owner, and visibility status.

**Important Notes**
- The `aws_ec2_ami` table only lists images in your account. To list other images shared with you, please use the `aws_ec2_ami_shared` table.

## Examples

### Basic info
Explore the different Amazon Machine Images (AMIs) in your AWS EC2 environment to understand their status, location, creation date, visibility, and root device. This is useful for auditing your resources, ensuring security compliance, and managing your infrastructure.

```sql+postgres
select
  name,
  image_id,
  state,
  image_location,
  creation_date,
  public,
  root_device_name
from
  aws_ec2_ami;
```

```sql+sqlite
select
  name,
  image_id,
  state,
  image_location,
  creation_date,
  public,
  root_device_name
from
  aws_ec2_ami;
```

### List public AMIs
Discover the segments that contain public Amazon Machine Images (AMIs) to help manage and maintain your AWS resources more effectively.

```sql+postgres
select
  name,
  image_id,
  public
from
  aws_ec2_ami
where
  public;
```

```sql+sqlite
select
  name,
  image_id,
  public
from
  aws_ec2_ami
where
  public = 1;
```

### List failed AMIs
Determine the areas in which Amazon Machine Images (AMIs) have failed. This can be useful for troubleshooting and identifying potential issues within your AWS EC2 instances.

```sql+postgres
select
  name,
  image_id,
  public,
  state
from
  aws_ec2_ami
where
  state = 'failed';
```

```sql+sqlite
select
  name,
  image_id,
  public,
  state
from
  aws_ec2_ami
where
  state = 'failed';
```

### Get volume info for each AMI
Explore the characteristics of each Amazon Machine Image (AMI), such as volume size and type, encryption status, and deletion policy. This information is vital for managing storage resources efficiently and ensuring data security within your AWS EC2 environment.

```sql+postgres
select
  name,
  image_id,
  mapping -> 'Ebs' ->> 'VolumeSize' as volume_size,
  mapping -> 'Ebs' ->> 'VolumeType' as volume_type,
  mapping -> 'Ebs' ->> 'Encrypted' as encryption_status,
  mapping -> 'Ebs' ->> 'KmsKeyId' as kms_key,
  mapping -> 'Ebs' ->> 'DeleteOnTermination' as delete_on_termination
from
  aws_ec2_ami
  cross join jsonb_array_elements(block_device_mappings) as mapping;
```

```sql+sqlite
select
  name,
  image_id,
  json_extract(mapping.value, '$.Ebs.VolumeSize') as volume_size,
  json_extract(mapping.value, '$.Ebs.VolumeType') as volume_type,
  json_extract(mapping.value, '$.Ebs.Encrypted') as encryption_status,
  json_extract(mapping.value, '$.Ebs.KmsKeyId') as kms_key,
  json_extract(mapping.value, '$.Ebs.DeleteOnTermination') as delete_on_termination
from
  aws_ec2_ami,
  json_each(block_device_mappings) as mapping;
```