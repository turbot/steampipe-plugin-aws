# Table: aws_guardduty_detector

Amazon GuardDuty is a threat detection service that continuously monitors for malicious activity and unauthorized behavior to protect your AWS accounts, workloads, and data stored in Amazon S3.

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