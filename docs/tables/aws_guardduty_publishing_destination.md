# Table: aws_guardduty_publishing_destination

AWS Guardduty Publishing Destinations provide a resource to export the guard duty findings. This requires an existing GuardDuty Detector.

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