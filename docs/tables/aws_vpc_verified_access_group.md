---
title: "Table: aws_vpc_verified_access_group - Query AWS VPC Verified Access Groups using SQL"
description: "Allows users to query VPC Verified Access Groups within AWS Virtual Private Cloud (VPC). This table provides information about each verified access group within a VPC, including details such as group ID, group name, and the VPC ID it is associated with."
---

# Table: aws_vpc_verified_access_group - Query AWS VPC Verified Access Groups using SQL

The `aws_vpc_verified_access_group` table in Steampipe provides information about each verified access group within a VPC in AWS Virtual Private Cloud (VPC). This table allows network administrators and security personnel to query group-specific details, including the group ID, group name, and the VPC ID it is associated with. Users can utilize this table to gain insights on access groups, such as which VPCs have certain access groups, the names of these groups, and more. The schema outlines the various attributes of the verified access group, including the group ID, group name, and associated VPC ID.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_verified_access_group` table, you can use the `.inspect aws_vpc_verified_access_group` command in Steampipe.

### Key columns:

- `group_id`: This is the unique identifier of the verified access group. It can be used to join with other tables that reference the group ID.
- `group_name`: This is the name of the verified access group. It can be useful for filtering or sorting results based on the group name.
- `vpc_id`: This is the ID of the VPC that the verified access group is associated with. It can be used to join with other tables that contain VPC information.

## Examples

### Basic info

```sql
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

```sql
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

### List active groups

```sql
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

```sql
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
