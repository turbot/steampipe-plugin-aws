---
title: "Steampipe Table: aws_ec2_launch_template - Query AWS EC2 Launch Templates using SQL"
description: "Allows users to query AWS EC2 Launch Templates to retrieve detailed information, including the associated AMI, instance type, key pair, security groups, and user data."
folder: "EC2"
---

# Table: aws_ec2_launch_template - Query AWS EC2 Launch Templates using SQL

The AWS EC2 Launch Template is a resource within the Amazon Elastic Compute Cloud (EC2) service. It allows you to save launch parameters within Amazon EC2 so you can quickly launch instances with those settings. This ensures consistency across instances and reduces the manual effort required to configure individual instances.

## Table Usage Guide

The `aws_ec2_launch_template` table in Steampipe provides you with information about EC2 Launch Templates within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query template-specific details, including instance type, key pair, security groups, and user data. You can utilize this table to gather insights on templates, such as associated AMIs, security configurations, instance configurations, and more. The schema outlines the various attributes of the EC2 Launch Template for you, including the template ID, creation date, default version, and associated tags.

## Examples

### Basic info
Explore which AWS EC2 launch templates have been created, by whom, and when. This can help in understanding the evolution of your infrastructure, including the original and most recent versions of each template.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by,
  default_version_number,
  latest_version_number
from
  aws_ec2_launch_template;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by,
  default_version_number,
  latest_version_number
from
  aws_ec2_launch_template;
```

### List launch templates created by a user
Discover the segments that include launch templates created by a specific user in AWS EC2. This is beneficial for understanding and managing user-specific resources within the cloud infrastructure.

```sql+postgres
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

```sql+sqlite
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
Identify recently created launch templates within the past month. This is useful for monitoring new additions and ensuring proper configuration and usage.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  create_time
from
  aws_ec2_launch_template
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  create_time
from
  aws_ec2_launch_template
where
  create_time >= datetime('now', '-30 days');
```