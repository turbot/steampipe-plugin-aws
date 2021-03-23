
# Table: aws_redshift_snapshot

Snapshots are point-in-time backups of a cluster. There are two types of snapshots: automated and manual. Amazon Redshift stores these snapshots internally in Amazon S3 by using an encrypted Secure Sockets Layer (SSL) connection.

## Examples

### Redshift snapshot basic info

```sql
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  encrypted
from
  aws_redshift_snapshot;
```


### List of all manual redshift snapshots

```sql
select
  snapshot_identifier,
  snapshot_type
from
  aws_redshift_snapshot
where
  snapshot_type = 'manual';
```


### List of snapshots which are not encrypted

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


### Redshift Cluster info of each redshift snapshot

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


### List snapshots of respective clusters those are not in a VPC

```sql
select
  snapshot_identifier,
  cluster_identifier,
  node_type,
  vpc_id
from
  aws_redshift_snapshot
where
  vpc_id is null;
```


### List snapshots whose AWS customer accounts are not authorized

```sql
select
  snapshot_identifier,
  accounts_with_restore_access
from
  aws_redshift_snapshot
where
  accounts_with_restore_access is null;

```