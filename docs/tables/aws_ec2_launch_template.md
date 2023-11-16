---
title: "Table: aws_ec2_launch_template - Query AWS EC2 Launch Templates using SQL"
description: "Allows users to query AWS EC2 Launch Templates to retrieve detailed information, including the associated AMI, instance type, key pair, security groups, and user data."
---

# Table: aws_ec2_launch_template - Query AWS EC2 Launch Templates using SQL

The `aws_ec2_launch_template` table in Steampipe provides information about EC2 Launch Templates within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query template-specific details, including instance type, key pair, security groups, and user data. Users can utilize this table to gather insights on templates, such as associated AMIs, security configurations, instance configurations, and more. The schema outlines the various attributes of the EC2 Launch Template, including the template ID, creation date, default version, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_launch_template` table, you can use the `.inspect aws_ec2_launch_template` command in Steampipe.

**Key columns**:

- `launch_template_id`: This is the unique ID of the launch template. It can be used to join with other tables when you need to correlate data based on the launch template.

- `default_version_number`: This is the default version number of the launch template. It can be used to join with other tables when you need to correlate data based on the default version of the launch template.

- `created_by`: This is the AWS account ID of the user who created the launch template. It can be used to join with other tables when you need to correlate data based on the creator of the launch template.

## Examples

### Basic info

```sql
select
  launch_template_name,
  launch_template_id,
  created_time,
  created_by,
  default_version_number,
  latest_version_number
from
  aws_ec2_launch_template;
```

### List launch templates created by a user

```sql
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by
from
  aws_ec2_launch_template
where
  created_by like '%turbot';
```

### List launch templates created in the last 30 days

```sql
select
  launch_template_name,
  launch_template_id,
  create_time
from
  aws_ec2_launch_template
where
  create_time >= now() - interval '30' day;
```