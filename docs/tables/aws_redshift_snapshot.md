---
title: "Table: aws_redshift_snapshot - Query AWS Redshift Snapshots using SQL"
description: "Allows users to query Redshift Snapshots, providing details about each snapshot's configuration, status, and associated metadata."
---

# Table: aws_redshift_snapshot - Query AWS Redshift Snapshots using SQL

The `aws_redshift_snapshot` table in Steampipe provides information about snapshots within AWS Redshift. This table allows DevOps engineers to query snapshot-specific details, including the snapshot status, creation time, source cluster, and associated metadata. Users can utilize this table to gather insights on snapshots, such as snapshot availability, size, and retention period. The schema outlines the various attributes of the Redshift snapshot, including the snapshot identifier, snapshot type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_snapshot` table, you can use the `.inspect aws_redshift_snapshot` command in Steampipe.

**Key columns**:

- `snapshot_identifier`: The unique identifier for the snapshot. This column is important as it is the primary identifier for each snapshot.
- `cluster_identifier`: The identifier of the cluster for which the snapshot was taken. This column is useful for joining this table with other tables that contain information about Redshift clusters.
- `snapshot_type`: The type of the snapshot. This column is useful for filtering snapshots based on whether they are automated or manual.

## Examples

### Basic info

```sql
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  encrypted
from
  aws_redshift_snapshot;
```


### List manual snapshots

```sql
select
  snapshot_identifier,
  snapshot_type
from
  aws_redshift_snapshot
where
  snapshot_type = 'manual';
```


### List unencrypted snapshots

```sql
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


### Get cluster info for each snapshot

```sql
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

```sql
select
  snapshot_identifier,
  accounts_with_restore_access
from
  aws_redshift_snapshot
where
  accounts_with_restore_access is not null;
```


### List accounts that are authorized to restore each snapshot

```sql
select
  snapshot_identifier,
  p ->> 'AccountId' as account_id,
  p ->> 'AccountAlias' as account_alias
from
  aws_redshift_snapshot,
  jsonb_array_elements(accounts_with_restore_access) as p;
```
