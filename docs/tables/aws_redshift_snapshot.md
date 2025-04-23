---
title: "Steampipe Table: aws_redshift_snapshot - Query AWS Redshift Snapshots using SQL"
description: "Allows users to query Redshift Snapshots, providing details about each snapshot's configuration, status, and associated metadata."
folder: "Redshift"
---

# Table: aws_redshift_snapshot - Query AWS Redshift Snapshots using SQL

The AWS Redshift Snapshot is a point-in-time copy of your data in AWS Redshift, a fully managed, petabyte-scale data warehouse service in the cloud. Snapshots are used to back up data and enable fast restore. They are automatically created by Redshift and can also be manually created by users.

## Table Usage Guide

The `aws_redshift_snapshot` table in Steampipe provides you with information about snapshots within AWS Redshift. This table allows you, as a DevOps engineer, to query snapshot-specific details, including the snapshot status, creation time, source cluster, and associated metadata. You can utilize this table to gather insights on snapshots, such as snapshot availability, size, and retention period. The schema outlines the various attributes of the Redshift snapshot for you, including the snapshot identifier, snapshot type, and associated tags.

## Examples

### Basic info
Explore which snapshots in your AWS Redshift database are encrypted. This can help you identify potential security risks and ensure that sensitive data is adequately protected.

```sql+postgres
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  encrypted
from
  aws_redshift_snapshot;
```

```sql+sqlite
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  encrypted
from
  aws_redshift_snapshot;
```


### List manual snapshots
Explore which snapshots have been manually created in your AWS Redshift environment. This can assist in understanding your data backup and recovery practices.

```sql+postgres
select
  snapshot_identifier,
  snapshot_type
from
  aws_redshift_snapshot
where
  snapshot_type = 'manual';
```

```sql+sqlite
select
  snapshot_identifier,
  snapshot_type
from
  aws_redshift_snapshot
where
  snapshot_type = 'manual';
```


### List unencrypted snapshots
Discover the segments that contain unencrypted snapshots in your AWS Redshift database. This is useful for identifying potential security risks and ensuring your data is properly protected.

```sql+postgres
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  number_of_nodes,
  encrypted
from
  aws_redshift_snapshot
where
  not encrypted;
```

```sql+sqlite
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  number_of_nodes,
  encrypted
from
  aws_redshift_snapshot
where
  encrypted = 0;
```


### Get cluster info for each snapshot
Explore the specifics of each snapshot, such as the associated cluster, its size, version, and potential restore options. This is useful for understanding the characteristics of each snapshot and for planning potential restore scenarios.

```sql+postgres
select
  snapshot_identifier,
  cluster_identifier,
  number_of_nodes,
  cluster_version,
  engine_full_version,
  restorable_node_types
from
  aws_redshift_snapshot;
```

```sql+sqlite
select
  snapshot_identifier,
  cluster_identifier,
  number_of_nodes,
  cluster_version,
  engine_full_version,
  restorable_node_types
from
  aws_redshift_snapshot;
```


### List snapshots that are shared with other accounts
Identify instances where snapshots are accessible to other accounts, a crucial step in assessing data sharing and privacy practices within your AWS Redshift environment.

```sql+postgres
select
  snapshot_identifier,
  accounts_with_restore_access
from
  aws_redshift_snapshot
where
  accounts_with_restore_access is not null;
```

```sql+sqlite
select
  snapshot_identifier,
  accounts_with_restore_access
from
  aws_redshift_snapshot
where
  accounts_with_restore_access is not null;
```


### List accounts that are authorized to restore each snapshot
Determine which accounts have permission to restore each snapshot in your AWS Redshift database. This is useful for auditing and managing data recovery permissions across your organization.

```sql+postgres
select
  snapshot_identifier,
  p ->> 'AccountId' as account_id,
  p ->> 'AccountAlias' as account_alias
from
  aws_redshift_snapshot,
  jsonb_array_elements(accounts_with_restore_access) as p;
```

```sql+sqlite
select
  snapshot_identifier,
  json_extract(p.value, '$.AccountId') as account_id,
  json_extract(p.value, '$.AccountAlias') as account_alias
from
  aws_redshift_snapshot,
  json_each(accounts_with_restore_access) as p;
```