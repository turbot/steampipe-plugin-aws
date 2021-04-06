# Table: aws_redshift_snapshot

Snapshots are point-in-time backups of a cluster. There are two types of snapshots: automated and manual. Amazon Redshift stores these snapshots internally in Amazon S3 by using an encrypted Secure Sockets Layer (SSL) connection.

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


### List all manual redshift snapshots

```sql
select
  snapshot_identifier,
  snapshot_type
from
  aws_redshift_snapshot
where
  snapshot_type = 'manual';
```


### List snapshots which are not encrypted

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


### Get cluster info of each redshift snapshot

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


### List accounts that are authorized to restore the snapshots

```sql
select
  snapshot_identifier,
  p ->> 'AccountId' as account_id,
  p ->> 'AccountAlias' as account_alias
from
  aws_redshift_snapshot,
  jsonb_array_elements(accounts_with_restore_access) as p;
```