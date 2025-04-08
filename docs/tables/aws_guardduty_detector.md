---
title: "Steampipe Table: aws_guardduty_detector - Query AWS GuardDuty Detector using SQL"
description: "Allows users to query AWS GuardDuty Detector data, including detector details, status, and associated metadata."
folder: "GuardDuty"
---

# Table: aws_guardduty_detector - Query AWS GuardDuty Detector using SQL

The AWS GuardDuty Detector is a part of the Amazon GuardDuty service, a threat detection service that continuously monitors for malicious activity and unauthorized behavior to protect your AWS accounts and workloads. It uses machine learning, anomaly detection, and integrated threat intelligence to identify and prioritize potential threats. GuardDuty analyzes tens of billions of events across multiple AWS data sources, such as AWS CloudTrail, Amazon VPC Flow Logs, and DNS logs.

## Table Usage Guide

The `aws_guardduty_detector` table in Steampipe provides you with information about detectors within AWS GuardDuty. This table allows you, as a security analyst, to query detector-specific details, including detector ID, creation timestamp, status, and associated tags. You can utilize this table to gather insights on detectors, such as their current status, when they were created, and more. The schema outlines the various attributes of the GuardDuty detector for you, including the detector ID, creation timestamp, status, service role, and associated tags.

## Examples

### Basic info
Uncover the details of your AWS GuardDuty detectors, such as their creation date and current status, to gain insights into your security infrastructure and help manage your service roles more effectively.

```sql+postgres
select
  detector_id,
  arn,
  created_at,
  status,
  service_role
from
  aws_guardduty_detector;
```

```sql+sqlite
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
Discover the segments that are actively monitoring for malicious activity by identifying the enabled detectors in your AWS GuardDuty. This can be useful for maintaining security and ensuring that all necessary detectors are functioning properly.

```sql+postgres
select
  detector_id,
  created_at,
  status
from
  aws_guardduty_detector
where
  status = 'ENABLED';
```

```sql+sqlite
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
Uncover the details of each detector's status and the status of their respective data sources. This allows for efficient monitoring and management of security and threat detection systems.

```sql+postgres
select
  detector_id,
  status as detector_status,
  data_sources -> 'CloudTrail' ->> 'Status' as cloud_trail_status,
  data_sources -> 'DNSLogs' ->> 'Status' as dns_logs_status,
  data_sources -> 'FlowLogs' ->> 'Status' as flow_logs_status
from
  aws_guardduty_detector;
```

```sql+sqlite
select
  detector_id,
  status as detector_status,
  json_extract(data_sources, '$.CloudTrail.Status') as cloud_trail_status,
  json_extract(data_sources, '$.DNSLogs.Status') as dns_logs_status,
  json_extract(data_sources, '$.FlowLogs.Status') as flow_logs_status
from
  aws_guardduty_detector;
```

### Get information about the master account relationship
Discover the segments that can provide insights into the relationship status of your master account with AWS GuardDuty. This can be particularly useful in understanding your account's security posture and managing potential threats.

```sql+postgres
select 
  detector_id,
  master_account ->> 'AccountId' as master_account_id,
  master_account ->> 'InvitationId' as invitation_id, 
  master_account ->> 'RelationshipStatus' as relationship_status 
from    
  aws_guardduty_detector
where master_account is not null;
```

```sql+sqlite
select 
  detector_id,
  json_extract(master_account, '$.AccountId') as master_account_id,
  json_extract(master_account, '$.InvitationId') as invitation_id, 
  json_extract(master_account, '$.RelationshipStatus') as relationship_status 
from    
  aws_guardduty_detector
where master_account is not null;
```