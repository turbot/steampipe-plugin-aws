---
title: "Steampipe Table: aws_ebs_snapshot - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS EBS snapshots, providing detailed information about each snapshot's configuration, status, and associated metadata."
folder: "EBS"
---

# Table: aws_ebs_snapshot - Query AWS Elastic Block Store (EBS) using SQL

The AWS Elastic Block Store (EBS) provides durable, block-level storage volumes for use with Amazon EC2 instances. These snapshots are point-in-time copies of your data that are used for enabling disaster recovery, migrating data across regions or accounts, improving backup compliance, or creating dev/test environments. EBS snapshots are incremental, meaning that only the blocks on the device that have changed after your most recent snapshot are saved.

## Table Usage Guide

The `aws_ebs_snapshot` table in Steampipe provides you with information about EBS snapshots within AWS Elastic Block Store (EBS). This table allows you, as a DevOps engineer, to query snapshot-specific details, including snapshot ID, description, status, volume size, and associated metadata. You can utilize this table to gather insights on snapshots, such as snapshots with public permissions, snapshots by volume, and more. The schema outlines the various attributes of the EBS snapshot for you, including the snapshot ID, creation time, volume ID, and associated tags.

**Important Notes**
- The `aws_ebs_snapshot` table lists all private snapshots by default.
- You can specify an owner alias, owner ID or snapshot ID** in the `where` clause (`where owner_alias=''`), (`where owner_id=''`) or (`where snapshot_id=''`) to list public or shared snapshots from a specific AWS account.

## Examples

### List of snapshots which are not encrypted
Discover the segments that include unencrypted snapshots in your AWS EBS environment. This is beneficial for enhancing your security measures by identifying potential vulnerabilities.

```sql+postgres
select
  snapshot_id,
  arn,
  encrypted
from
  aws_ebs_snapshot
where
  not encrypted;
```

```sql+sqlite
select
  snapshot_id,
  arn,
  encrypted
from
  aws_ebs_snapshot
where
  encrypted = 0;
```

### List of EBS snapshots which are publicly accessible
Determine the areas in which EBS snapshots are publicly accessible to identify potential security risks. This query is used to uncover instances where EBS snapshots may be exposed to all users, which could lead to unauthorized data access.

```sql+postgres
select
  snapshot_id,
  arn,
  volume_id,
  perm ->> 'UserId' as userid,
  perm ->> 'Group' as group
from
  aws_ebs_snapshot
  cross join jsonb_array_elements(create_volume_permissions) as perm
where
  perm ->> 'Group' = 'all';
```

```sql+sqlite
select
  snapshot_id,
  arn,
  volume_id,
  json_extract(perm, '$.UserId') as userid,
  json_extract(perm, '$.Group') as group
from
  aws_ebs_snapshot,
  json_each(create_volume_permissions) as perm
where
  json_extract(perm, '$.Group') = 'all';
```

### Find the Account IDs with which the snapshots are shared
Determine the accounts that have access to specific snapshots in your AWS EBS setup. This can be useful for auditing purposes, ensuring that only authorized accounts have access to your data.

```sql+postgres
select
  snapshot_id,
  volume_id,
  perm ->> 'UserId' as account_ids
from
  aws_ebs_snapshot
  cross join jsonb_array_elements(create_volume_permissions) as perm;
```

```sql+sqlite
select
  snapshot_id,
  volume_id,
  json_extract(perm.value, '$.UserId') as account_ids
from
  aws_ebs_snapshot
  cross join json_each(create_volume_permissions) as perm;
```

### Find the snapshot count per volume
Assess the elements within each volume to determine the number of snapshots associated with it. This can be useful for understanding the backup frequency and data recovery potential for each volume.

```sql+postgres
select
  volume_id,
  count(snapshot_id) as snapshot_id
from
  aws_ebs_snapshot
group by
  volume_id;
```

```sql+sqlite
select
  volume_id,
  count(snapshot_id) as snapshot_id
from
  aws_ebs_snapshot
group by
  volume_id;
```

### List snapshots owned by a specific AWS account
Determine the areas in which specific AWS accounts own snapshots. This can be useful for managing and tracking resources across different accounts in a cloud environment.

```sql+postgres
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_id = '859788737657';
```

```sql+sqlite
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_id = '859788737657';
```

### Get a specific snapshot by ID
Discover the specific details of a particular snapshot using its unique identifier. This can be useful for auditing purposes, such as confirming the owner or checking if the snapshot is encrypted.

```sql+postgres
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  snapshot_id = 'snap-07bf4f91353ad71ae';
```

```sql+sqlite
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  snapshot_id = 'snap-07bf4f91353ad71ae';
```

### List snapshots owned by Amazon (Note: This will attempt to list ALL public snapshots)
Discover the segments that are owned by Amazon, specifically focusing on public snapshots. This is particularly useful for gaining insights into the distribution and ownership of snapshots within the Amazon ecosystem.

```sql+postgres
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_alias = 'amazon'
```

```sql+sqlite
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_alias = 'amazon'
```