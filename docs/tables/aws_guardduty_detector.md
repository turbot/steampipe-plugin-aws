---
title: "Table: aws_guardduty_detector - Query AWS GuardDuty Detector using SQL"
description: "Allows users to query AWS GuardDuty Detector data, including detector details, status, and associated metadata."
---

# Table: aws_guardduty_detector - Query AWS GuardDuty Detector using SQL

The `aws_guardduty_detector` table in Steampipe provides information about detectors within AWS GuardDuty. This table allows security analysts to query detector-specific details, including detector ID, creation timestamp, status, and associated tags. Users can utilize this table to gather insights on detectors, such as their current status, when they were created, and more. The schema outlines the various attributes of the GuardDuty detector, including the detector ID, creation timestamp, status, service role, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_detector` table, you can use the `.inspect aws_guardduty_detector` command in Steampipe.

### Key columns:

- `detector_id`: The unique ID of the detector. This can be used to join this table with other tables that contain GuardDuty detector information.
- `status`: The status of the detector (e.g., ENABLED or DISABLED). This is important for understanding the operational state of the detector.
- `created_at`: The timestamp of when the detector was created. This can be useful for auditing and tracking the lifecycle of the detector.

## Examples

### Basic info

```sql
select
  detector_id,
  arn,
  created_at,
  status,
  service_role
from
  aws_guardduty_detector;
```

### List enabled detectors

```sql
select
  detector_id,
  created_at,
  status
from
  aws_guardduty_detector
where
  status = 'ENABLED';
```

### Get data source status info for each detector

```sql
select
  detector_id,
  status as detector_status,
  data_sources -> 'CloudTrail' ->> 'Status' as cloud_trail_status,
  data_sources -> 'DNSLogs' ->> 'Status' as dns_logs_status,
  data_sources -> 'FlowLogs' ->> 'Status' as flow_logs_status
from
  aws_guardduty_detector;
```

### Get information about the master account relationship

```sql
select 
  detector_id,
  master_account ->> 'AccountId' as master_account_id,
  master_account ->> 'InvitationId' as invitation_id, 
  master_account ->> 'RelationshipStatus' as relationship_status 
from    
  aws_guardduty_detector
where master_account is not null;
```