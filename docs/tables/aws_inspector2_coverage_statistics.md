---
title: "Table: aws_inspector2_coverage_statistics - Query AWS Inspector2 Coverage Statistics using SQL"
description: "Allows users to query AWS Inspector2 Coverage Statistics to obtain detailed information about the assessment targets and the number of instances they cover."
---

# Table: aws_inspector2_coverage_statistics - Query AWS Inspector2 Coverage Statistics using SQL

The `aws_inspector2_coverage_statistics` table in Steampipe provides information about AWS Inspector2's coverage statistics. This table allows DevOps engineers, security analysts, and other technical professionals to query detailed information about the assessment targets and the number of instances they cover. Users can utilize this table to gather insights on assessment targets, including their ARNs, the number of instances they cover, and other associated metadata. The schema outlines the various attributes of the coverage statistics, including the assessment target ARN, the instance count, and the agent ID.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector2_coverage_statistics` table, you can use the `.inspect aws_inspector2_coverage_statistics` command in Steampipe.

**Key columns**:

- `assessment_target_arn`: The Amazon Resource Name (ARN) of the assessment target. This column is useful for joining with other tables that contain assessment target information.
- `instance_count`: The number of instances that the assessment target covers. This column is useful for understanding the extent of the coverage of each assessment target.
- `agent_id`: The ID of the agent. This column is useful for joining with other tables that contain agent-specific information.

## Examples

### Basic info

```sql
select
  total_counts,
  counts_by_group
from
  aws_inspector2_coverage_statistics;
```

### Get the count of resources within a group

```sql
select
  g ->> 'Count' as count,
  g ->> 'GroupKey' as group_key
from
  aws_inspector2_coverage_statistics,
  jsonb_array_elements(counts_by_group) as g;
```