---
title: "Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL"
description: "Allows users to query shared Amazon Machine Images (AMIs) in AWS EC2"
---

# Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL

The `aws_ec2_ami_shared` table in Steampipe provides information about shared Amazon Machine Images (AMIs) within AWS EC2. This table allows system administrators and DevOps engineers to query shared AMI-specific details, including image ID, creation date, state, and associated tags. Users can utilize this table to gather insights on shared AMIs, such as their availability, permissions, and associated metadata. The schema outlines the various attributes of the shared AMI, including the image type, launch permissions, and virtualization type.

**You must specify an Owner ID or Image ID** in a `where` clause (`where owner_id='`), (`where image_id='`).

The `aws_ec2_ami_shared` table can list any image but you must specify `owner_id` or `image_id`.
If you want to list all of the images in your account then you can use the `aws_ec2_ami` table.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_ami_shared` table, you can use the `.inspect aws_ec2_ami_shared` command in Steampipe.

### Key columns:

- `image_id`: The unique identifier for the AMI. This can be used to join this table with other tables that contain AMI information.
- `creation_date`: The date when the AMI was created. This can be useful for tracking the age of the AMI and determining when it may need to be updated or replaced.
- `state`: The current state of the AMI (available, pending, or failed). This can be useful for tracking the availability of the AMI and diagnosing any potential issues.

## Examples

### Basic info

```sql
select
  name,
  image_id,
  state,
  image_location,
  creation_date,
  public,
  root_device_name
from
  aws_ec2_ami_shared
where
  owner_id = '137112412989';
```

### List arm64 AMIs

```sql
select
  name,
  image_id,
  state,
  image_location,
  creation_date,
  public,
  root_device_name
from
  aws_ec2_ami_shared
where
  owner_id = '137112412989'
  and architecture = 'arm64';
```

### List EC2 instances using AMIs owned by a specific AWS account

```sql
select
  i.title,
  i.instance_id,
  i.image_id,
  ami.name,
  ami.description,
  ami.platform_details
from
  aws_ec2_instance as i
 join aws_ec2_ami_shared as ami on i.image_id = ami.image_id
where
  ami.owner_id = '137112412989';
```
