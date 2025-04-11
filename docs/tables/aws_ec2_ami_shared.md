---
title: "Steampipe Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL"
description: "Allows users to query shared Amazon Machine Images (AMIs) in AWS EC2"
folder: "EC2"
---

# Table: aws_ec2_ami_shared - Query AWS EC2 AMI using SQL

The AWS EC2 AMI (Amazon Machine Image) provides the information necessary to launch an instance, which is a virtual server in the cloud. You can specify an AMI when you launch instances, and you can launch as many instances from the AMI as you need. You can also share your own custom AMI with other AWS accounts, enabling them to launch instances with identical configurations.

## Table Usage Guide

The `aws_ec2_ami_shared` table in Steampipe provides you with information about shared Amazon Machine Images (AMIs) within AWS EC2. This table enables you, as a system administrator or DevOps engineer, to query shared AMI-specific details, including image ID, creation date, state, and associated tags. You can utilize this table to gather insights on shared AMIs, such as their availability, permissions, and associated metadata. The schema outlines the various attributes of the shared AMI, including the image type, launch permissions, and virtualization type.

**Important Notes**
- You must specify an Owner ID or Image ID in the `where` clause (`where owner_id='`), (`where image_id='`).
- The `aws_ec2_ami_shared` table can list any image but you must specify `owner_id` or `image_id`.
- To optimize query timing and API calls, use the optional query parameters `owner_ids` or `image_ids` to perform batch operations.
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

### Retrieve details of multiple shared AMIs in a single query
Fetches metadata of multiple shared AMIs, including their state, visibility, and creation details, to streamline AMI management.

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
  image_ids = '["ami-08df646e18b182346", "ami-04c5f154a6c2fec00"]';
```

### Batch API operation, ensuring AMIs are from trusted sources
Any AWS customer can publish an Amazon Machine Image (AMI) for other AWS customers to launch instances from. AWS only vets a handful of images in the AWS Marketplace, there is no guarantee that other publicly shared AMIs are free of vulnerabilities or malicious code. While it's common for vendors to share their software as an AMI, it's also possible someone in your organization has launched an instance from a compromised image.

```sql+postgres
with instances as (
  select
    instance_id,
    instance_type,
    account_id,
    tags ->> 'Name' as instance_name,
    _ctx ->> 'connection_name' as account_name,
    instance_state,
    region,
    image_id
  from
    aws_ec2_instance
),
all_image_ids as (
  select
    json_agg(image_id)::jsonb as image_ids  -- Cast to jsonb
  from
    instances
),
shared_ami as (
  select
    s.*
  from
    aws_ec2_ami_shared as s,
    all_image_ids
  where s.image_ids = all_image_ids.image_ids
)
select distinct
  shared_ami.image_id as image_id,
  shared_ami.owner_id as image_owner_id,
  shared_ami.image_owner_alias as image_owner_name,
  instances.instance_name,
  instances.account_name,
  instances.region,
  shared_ami.name as image_name
from
  instances
left join shared_ami on shared_ami.image_id=instances.image_id
where shared_ami.image_owner_alias != 'amazon'
and shared_ami.image_owner_alias != 'self';
```

```sql+sqlite
with instances as (
  select
    instance_id,
    instance_type,
    account_id,
    json_extract(tags, '$.Name') as instance_name,
    json_extract(_ctx, '$.connection_name') as account_name,
    instance_state,
    region,
    image_id
  from
    aws_ec2_instance
),
all_image_ids as (
  select
    json_group_array(image_id) as image_ids
  from
    instances
),
shared_ami as (
  select
    s.*,
    ai.image_ids
  from
    aws_ec2_ami_shared as s,
    all_image_ids as ai
  where json_array_contains(ai.image_ids, s.image_id)
)
select distinct
  shared_ami.image_id as image_id,
  shared_ami.owner_id as image_owner_id,
  shared_ami.image_owner_alias as image_owner_name,
  instances.instance_name,
  instances.account_name,
  instances.region,
  shared_ami.name as image_name
from
  instances
left join shared_ami on shared_ami.image_id = instances.image_id
where shared_ami.image_owner_alias != 'amazon'
  and shared_ami.image_owner_alias != 'self';
```
