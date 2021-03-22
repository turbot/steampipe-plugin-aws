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


## List of workload with production environment

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


## List of aws regions associated with the workload

```sql
select
  workload_name,
  workload_id,
  environment,
  aws_regions
from
  aws_wellarchitected_workload;
```


## Industry type of the workloads

```sql
select
  workload_name,
  workload_id,
  industry_type
from
  aws_wellarchitected_workload;
```


## List of workloads with IsReviewOwnerUpdateAcknowledged not enabled

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
