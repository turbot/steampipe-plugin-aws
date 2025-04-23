---
title: "Steampipe Table: aws_wellarchitected_workload - Query AWS Well-Architected Tool Workloads using SQL"
description: "Allows users to query AWS Well-Architected Tool Workloads to retrieve and manage workload data, including workload names, ARNs, risk counts, and improvement statuses."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_workload - Query AWS Well-Architected Tool Workloads using SQL

The AWS Well-Architected Tool is a service that helps you review the state of your workloads and compares them to the latest AWS architectural best practices. The tool measures your workloads across five pillars of a well-architected framework: operational excellence, security, reliability, performance efficiency, and cost optimization. It provides a consistent approach for customers and partners to evaluate architectures, and implement designs that can scale over time.

## Table Usage Guide

The `aws_wellarchitected_workload` table in Steampipe provides you with information about workloads within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query workload-specific details, including workload name, ARN, risk count, and improvement status. You can utilize this table to gather insights on workloads, such as identifying workloads with high risk counts, tracking improvement status, and more. The schema outlines the various attributes of the workload for you, including the workload ARN, creation date, risk count, improvement status, and associated tags.

## Examples

## Basic info

```sql+postgres
select
  workload_name,
  workload_id,
  environment,
  industry,
  owner
from
  aws_wellarchitected_workload;
```

```sql+sqlite
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

```sql+postgres
select
  workload_name,
  workload_id,
  environment
from
  aws_wellarchitected_workload
where
  environment = 'PRODUCTION';
```

```sql+sqlite
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

```sql+postgres
select
  workload_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_workload;
```

```sql+sqlite
select
  workload_name,
  workload_id,
  json_extract(risk_counts, '$.HIGH') as high_risk_counts
from
  aws_wellarchitected_workload;
```


## List workloads that do not require a review owner

```sql+postgres
select
  workload_name,
  workload_id,
  is_review_owner_update_acknowledged
from
  aws_wellarchitected_workload
where
  not is_review_owner_update_acknowledged;
```

```sql+sqlite
select
  workload_name,
  workload_id,
  is_review_owner_update_acknowledged
from
  aws_wellarchitected_workload
where
  is_review_owner_update_acknowledged is not 1;
```