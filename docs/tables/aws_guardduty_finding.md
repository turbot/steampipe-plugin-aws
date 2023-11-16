---
title: "Table: aws_guardduty_finding - Query AWS GuardDuty Findings using SQL"
description: "Allows users to query AWS GuardDuty Findings to access detailed information about potential security threats or suspicious activities detected in their AWS environment."
---

# Table: aws_guardduty_finding - Query AWS GuardDuty Findings using SQL

The `aws_guardduty_finding` table in Steampipe provides information about findings reported by AWS GuardDuty. This table allows security analysts to query finding-specific details, including threat type, severity, and associated resources. Users can utilize this table to gather insights on potential security threats, such as unauthorized access attempts, data breaches, or compromised instances. The schema outlines the various attributes of the GuardDuty finding, including the finding ID, detector ID, account ID, region, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_finding` table, you can use the `.inspect aws_guardduty_finding` command in Steampipe.

**Key columns**:

- `id`: The unique identifier for the finding. This can be used to join this table with other tables that contain information about specific findings.
- `detector_id`: The ID of the GuardDuty detector that generated the finding. This can be used to join with the `aws_guardduty_detector` table which contains information about the detectors.
- `region`: The AWS region in which the finding was detected. This can be used to join with other tables that contain region-specific information.

## Examples

### Basic info

```sql
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding;
```

### List findings that are not archived

```sql
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding
where
  service ->> 'Archived' = 'false';
```
