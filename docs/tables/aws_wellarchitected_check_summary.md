---
title: "Table: aws_wellarchitected_check_summary - Query AWS Well-Architected Tool Check Summary using SQL"
description: "Allows users to query AWS Well-Architected Tool Check Summary for detailed information about the checks for all workloads. This table provides insights into the state of your workloads, highlighting potential risks and areas for improvement."
---

# Table: aws_wellarchitected_check_summary - Query AWS Well-Architected Tool Check Summary using SQL

The `aws_wellarchitected_check_summary` table in Steampipe provides information about the check summaries of all workloads within AWS Well-Architected Tool. This table allows DevOps engineers to query check-specific details, including the workload ID, lens alias, pillar ID, and risk counts. Users can utilize this table to gather insights on checks, such as the number of high-risk items, medium-risk items, and checks that are not applicable. The schema outlines the various attributes of the check summary, including the workload ID, lens alias, pillar ID, and risk counts.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_check_summary` table, you can use the `.inspect aws_wellarchitected_check_summary` command in Steampipe.

**Key columns**:

- `workload_id`: This is the ID of the workload. It can be used to join with the `aws_wellarchitected_workload` table to get more details about the workload.
- `lens_alias`: This is the alias of the lens. It can be used to join with the `aws_wellarchitected_lens` table to get more information about the lens.
- `pillar_id`: This is the ID of the pillar. It can be used to join with the `aws_wellarchitected_pillar` table to get more information about the pillar.

## Examples

### Basic info

```sql
select
  id,
  name,
  description,
  jsonb_pretty(account_summary) as account_summary,
  choice_id,
  lens_arn,
  pillar_id,
  question_id,
  status,
  region,
  workload_id
from
  aws_wellarchitected_check_summary;
```

### Get summarized trusted advisor check report for a workload

```sql
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  status,
  choice_id,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

### List trusted advisor checks with errors

```sql
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  status = 'ERROR';
```

### Get account summary for trusted advisor checks

```sql
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary;
```

### Get account summary for trusted advisor checks for well-architected lens in a particular workload

```sql
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```