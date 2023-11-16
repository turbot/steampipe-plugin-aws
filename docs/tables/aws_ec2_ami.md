---
title: "Table: aws_ec2_ami - Query AWS EC2 AMI using SQL"
description: "Allows users to query AWS EC2 AMIs (Amazon Machine Images) to retrieve detailed information about each AMI available in the AWS account."
---

# Table: aws_ec2_ami - Query AWS EC2 AMI using SQL

The `aws_ec2_ami` table in Steampipe provides information about AMIs (Amazon Machine Images) within Amazon Elastic Compute Cloud (Amazon EC2). This table allows DevOps engineers, system administrators, and other technical professionals to query AMI-specific details, including its attributes, block device mappings, and associated tags. Users can utilize this table to gather insights on AMIs, such as identifying unused or outdated AMIs, verifying AMI permissions, and more. The schema outlines the various attributes of the AMI, including the AMI ID, creation date, owner, and visibility status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_ami` table, you can use the `.inspect aws_ec2_ami` command in Steampipe.

Key columns:

- `image_id`: The unique identifier for the AMI. This column can be used to join this table with other tables that contain AMI IDs, such as the `aws_ec2_instance` table.
- `owner_id`: The AWS account ID of the AMI owner. This column can be used to filter or join the table based on the owner of the AMI.
- `state`: The current state of the AMI (available, pending, or failed). This column can be useful for filtering the table to only include AMIs in a certain state.

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
  aws_ec2_ami;
```

### List public AMIs

```sql
select
  name,
  image_id,
  public
from
  aws_ec2_ami
where
  public;
```

### List failed AMIs

```sql
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

```sql
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
