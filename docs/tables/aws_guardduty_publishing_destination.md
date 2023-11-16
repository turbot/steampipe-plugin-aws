---
title: "Table: aws_guardduty_publishing_destination - Query AWS GuardDuty Publishing Destinations using SQL"
description: "Allows users to query AWS GuardDuty Publishing Destinations to retrieve information about where GuardDuty findings are published."
---

# Table: aws_guardduty_publishing_destination - Query AWS GuardDuty Publishing Destinations using SQL

The `aws_guardduty_publishing_destination` table in Steampipe provides information about publishing destinations in AWS GuardDuty. This table allows security analysts and DevOps engineers to query destination-specific details, including the destination type, status, and associated metadata. Users can utilize this table to gather insights on destinations, such as their current statuses, the types of destinations, and more. The schema outlines the various attributes of the publishing destination, including the destination type, destination ARN, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_publishing_destination` table, you can use the `.inspect aws_guardduty_publishing_destination` command in Steampipe.

**Key columns**:

- `detector_id`: The unique ID of the detector that the publishing destination is associated with. This can be used to join with the `aws_guardduty_detector` table.
- `destination_arn`: The ARN of the destination. This can be used to join with other AWS resource tables that may refer to the destination.
- `destination_type`: The type of the destination (e.g., S3, SNS). This can be useful for filtering or grouping by destination type.

## Examples

### Basic info

```sql
select
  detector_id,
  destination_id,
  arn,
  destination_arn,
  status
from
  aws_guardduty_publishing_destination;
```

### List unverified publishing destinations

```sql
select
  destination_id,
  arn,
  status
from
  aws_guardduty_publishing_destination
where
  status = 'PENDING_VERIFICATION';
```

### List publishing destinations which are not encrypted

```sql
select
  destination_id,
  kms_key_arn,
  status,
  destination_type
from
  aws_guardduty_publishing_destination
where
  kms_key_arn is null;
```

### Count publishing destinations by type

```sql
select
  destination_type,
  count(destination_id)
from
  aws_guardduty_publishing_destination
group by 
  destination_type
order by
  count desc;
```

### Get bucket policies for S3 bucket publishing destinations

```sql
select
  d.destination_id,
  d.destination_arn,
  d.destination_type,
  p ->> 'Sid' as sid,
  p ->> 'Action' as policy_action,
  p ->> 'Effect' as effect,
  p -> 'Principal' ->> 'Service' as principal_service
from
  aws_guardduty_publishing_destination as d,
  aws_s3_bucket as s,
  jsonb_array_elements(s.policy -> 'Statement') as p
where
  d.destination_type = 'S3'
and
  s.arn = d.destination_arn;
```

### Get KMS key policies associated with publishing destinations

```sql
select
  d.destination_id,
  p ->> 'Sid' as sid,
  p ->> 'Action' as policy_action,
  p ->> 'Effect' as effect,
  p ->> 'Principal' as policy_principal,
  p ->> 'Condition' as policy_condition
from
  aws_guardduty_publishing_destination as d,
  aws_kms_key as k,
  jsonb_array_elements(k.policy -> 'Statement') as p
where
  d.kms_key_arn is not null
and
  k.arn = d.kms_key_arn;
```