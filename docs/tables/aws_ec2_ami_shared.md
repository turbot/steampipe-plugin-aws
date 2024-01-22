---
title: "Steampipe Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL"
description: "Allows users to query shared Amazon Machine Images (AMIs) in AWS EC2"
---

# Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL

The AWS EC2 AMI (Amazon Machine Image) provides the information necessary to launch an instance, which is a virtual server in the cloud. You can specify an AMI when you launch instances, and you can launch as many instances from the AMI as you need. You can also share your own custom AMI with other AWS accounts, enabling them to launch instances with identical configurations.

## Table Usage Guide

The `aws_ec2_ami_shared` table in Steampipe provides you with information about shared Amazon Machine Images (AMIs) within AWS EC2. This table enables you, as a system administrator or DevOps engineer, to query shared AMI-specific details, including image ID, creation date, state, and associated tags. You can utilize this table to gather insights on shared AMIs, such as their availability, permissions, and associated metadata. The schema outlines the various attributes of the shared AMI, including the image type, launch permissions, and virtualization type.

**Important Notes**
- You must specify an Owner ID or Image ID in the `where` clause (`where owner_id='`), (`where image_id='`).
- The `aws_ec2_ami_shared` table can list any image but you must specify `owner_id` or `image_id`.
- If you want to list all of the images in your account then you can use the `aws_ec2_ami` table.

## Examples

### Basic info
Explore which AWS EC2 shared AMI resources are owned by a specific user to understand their configurations. This can be useful in auditing access and managing resources across your organization.

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
  aws_ec2_ami_shared
where
  owner_id = '137112412989';
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
  aws_ec2_ami_shared
where
  owner_id = '137112412989';
```

### List arm64 AMIs
Explore which Amazon Machine Images (AMIs) with 'arm64' architecture are shared by a specific owner. This can be useful in identifying suitable AMIs for deployment on 'arm64' architecture instances.

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
  aws_ec2_ami_shared
where
  owner_id = '137112412989'
  and architecture = 'arm64';
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
  aws_ec2_ami_shared
where
  owner_id = '137112412989'
  and architecture = 'arm64';
```

### List EC2 instances using AMIs owned by a specific AWS account
Explore which EC2 instances are using AMIs owned by a particular AWS account. This is useful to maintain account security and manage resources efficiently.

```sql+postgres
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

```sql+sqlite
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