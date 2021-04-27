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
