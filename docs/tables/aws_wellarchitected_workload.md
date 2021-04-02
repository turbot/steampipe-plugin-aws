# Table: aws_wellarchitected_workload

AWS Well-Architected helps cloud architects build secure, high-performing, resilient, and efficient infrastructure for their applications and workloads.

## Examples

## Basic info

```sql
select
  workload_name,
  workload_id,
  environment,
  industry,
  owner
from
  aws_wellarchitected_workload;
```


## List production workloads

```sql
select
  workload_name,
  workload_id,
  environment
from
  aws_wellarchitected_workload
where
  environment = 'PRODUCTION';
```


## Get high risk issue counts for each workload

```sql
select
  workload_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_workload;
```


## List workloads that do not require a review owner

```sql
select
  workload_name,
  workload_id,
  is_review_owner_update_acknowledged
from
  aws_wellarchitected_workload
where
  not is_review_owner_update_acknowledged;
```
