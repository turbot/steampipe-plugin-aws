---
title: "Table: aws_wellarchitected_lens_review_report - Query AWS Well-Architected Tool Lens Review Report using SQL"
description: "Allows users to query Lens Review Reports in the AWS Well-Architected Tool."
---

# Table: aws_wellarchitected_lens_review_report - Query AWS Well-Architected Tool Lens Review Report using SQL

The `aws_wellarchitected_lens_review_report` table in Steampipe provides information about Lens Review Reports within the AWS Well-Architected Tool. This table allows DevOps engineers, architects, and system administrators to query details about the lens review reports, including the lens alias, lens name, lens version, and associated metadata. Users can utilize this table to gather insights on lens review reports, such as identifying the lens version, verifying the lens status, understanding the lens notes, and more. The schema outlines the various attributes of the Lens Review Report, including the lens version, lens status, lens notes, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_lens_review_report` table, you can use the `.inspect aws_wellarchitected_lens_review_report` command in Steampipe.

**Key columns**:

- `lens_alias`: The alias of the lens. This can be used to join with other tables that contain lens alias information.
- `lens_version`: The version of the lens. This is useful for tracking the version history of the lens.
- `lens_status`: The status of the lens. This can be used to monitor and manage the status of the lens.

## Examples

### Basic info

```sql
select
  lens_alias,
  lens_arn,
  workload_id,
  milestone_number,
  base64_string
from
  aws_wellarchitected_lens_review_report;
```

### Get workload details for the review report

```sql
select
  r.workload_name,
  r.workload_id,
  r.base64_string,
  w.environment,
  w.is_review_owner_update_acknowledged
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_workload as w
where
  and r.workload_id = w.workload_id;
```

### Get the review report of custom lenses

```sql
select
  r.lens_alias,
  r.lens_arn,
  r.base64_string,
  l.lens_type
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_lens as l
where
  l.lens_type <> `aws_OFFICIAL';
```