---
title: "Table: aws_wellarchitected_workload - Query AWS Well-Architected Tool Workloads using SQL"
description: "Allows users to query AWS Well-Architected Tool Workloads to retrieve and manage workload data, including workload names, ARNs, risk counts, and improvement statuses."
---

# Table: aws_wellarchitected_workload - Query AWS Well-Architected Tool Workloads using SQL

The `aws_wellarchitected_workload` table in Steampipe provides information about workloads within AWS Well-Architected Tool. This table allows DevOps engineers to query workload-specific details, including workload name, ARN, risk count, and improvement status. Users can utilize this table to gather insights on workloads, such as identifying workloads with high risk counts, tracking improvement status, and more. The schema outlines the various attributes of the workload, including the workload ARN, creation date, risk count, improvement status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_workload` table, you can use the `.inspect aws_wellarchitected_workload` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the workload. This is the unique identifier for the workload and can be used to join this table with others that contain workload ARN data.
- `workload_name`: The name of the workload. This can be used to filter or sort data by workload name.
- `risk_count`: The count of risks identified for the workload. This column can be used to identify workloads with high risk counts.

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
