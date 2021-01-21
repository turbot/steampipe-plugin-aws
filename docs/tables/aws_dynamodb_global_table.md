# Table: aws_dynamodb_global_table

Global tables eliminate the difficult work of replicating data between regions and resolving update conflicts, enabling you to focus on application's business logic. A global table consists of multiple replica tables (one per region that you choose) that DynamoDB treats as a single unit.

## Examples

### List of regions where global table replicas are present

```sql
select
  global_table_name,
  rg -> 'RegionName' as region_name
from
  aws_dynamodb_global_table
  cross join jsonb_array_elements(replication_group) as rg;
```


### DynamoDB global table replica info

```sql
select
  global_table_name,
  global_table_status,
  rg -> 'GlobalSecondaryIndexes' as global_secondary_indexes,
  rg -> 'RegionName' as region_name,
  rg -> 'ReplicaInaccessibleDateTime' as replica_inaccessible_date_time,
  rg -> 'ReplicaStatus' as replica_status,
  rg -> 'ReplicaStatusDescription' as replica_status_description,
  rg -> 'ReplicaStatusPercentProgress' as replica_status_percent_progress
from
  aws_dynamodb_global_table
  cross join jsonb_array_elements(replication_group) as rg;
```
