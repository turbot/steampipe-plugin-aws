---
title: "Steampipe Table: aws_vpc_verified_access_group - Query AWS VPC Verified Access Groups using SQL"
description: "Allows users to query VPC Verified Access Groups within AWS Virtual Private Cloud (VPC). This table provides information about each verified access group within a VPC, including details such as group ID, group name, and the VPC ID it is associated with."
folder: "VPC"
---

# Table: aws_vpc_verified_access_group - Query AWS VPC Verified Access Groups using SQL

The AWS VPC Verified Access Groups are used to manage access to your Virtual Private Cloud (VPC) resources. They enable you to specify which principals in your AWS environment have access to the resources in your VPC. This helps in maintaining security and control over your networked AWS resources.

## Table Usage Guide

The `aws_vpc_verified_access_group` table in Steampipe provides you with information about each verified access group within a VPC in AWS Virtual Private Cloud (VPC). This table enables you, as a network administrator or security personnel, to query group-specific details, including the group ID, group name, and the VPC ID it is associated with. You can utilize this table to gain insights on access groups, such as which VPCs have certain access groups, the names of these groups, and more. The schema outlines for you the various attributes of the verified access group, including the group ID, group name, and associated VPC ID.

## Examples

### Basic info
Determine the areas in which your AWS VPC verified access groups were created and last updated. This allows you to monitor and manage your AWS resources effectively, ensuring optimal security and performance.

```sql+postgres
select
  verified_access_group_id,
  arn,
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group;
```

```sql+sqlite
select
  verified_access_group_id,
  arn,
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group;
```

### List groups older than 30 days
Uncover the details of specific access groups in your AWS VPC that were created more than 30 days ago. This can be useful for routine cleanup or auditing purposes, ensuring your environment remains optimized and secure.

```sql+postgres
select
  verified_access_group_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  creation_time <= now() - interval '30' day;
```

```sql+sqlite
select
  verified_access_group_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  creation_time <= datetime('now', '-30 day');
```

### List active groups
Discover the segments that are currently active within your AWS VPC by identifying groups that have not been deleted. This can help in managing and maintaining the security and access control of your virtual private cloud.

```sql+postgres
select
  verified_access_group_id,
  creation_time,
  deletion_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  deletion_time is null;
```

```sql+sqlite
select
  verified_access_group_id,
  creation_time,
  deletion_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  deletion_time is null;
```

### Get trusted provider details for each group
Explore the trusted provider details associated with each group to understand when and how these relationships were established, providing valuable context for managing and optimizing your AWS VPC access security.

```sql+postgres
select
  g.verified_access_group_id,
  g.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  jsonb_pretty(i.verified_access_trust_providers) as verified_access_trust_providers
from
  aws_vpc_verified_access_group as g,
  aws_vpc_verified_access_instance as i
where
  g.verified_access_instance_id = i.verified_access_instance_id;
```

```sql+sqlite
select
  g.verified_access_group_id,
  g.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  i.verified_access_trust_providers as verified_access_trust_providers
from
  aws_vpc_verified_access_group as g
join
  aws_vpc_verified_access_instance as i
on
  g.verified_access_instance_id = i.verified_access_instance_id;
```