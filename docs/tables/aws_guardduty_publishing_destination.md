---
title: "Steampipe Table: aws_guardduty_publishing_destination - Query AWS GuardDuty Publishing Destinations using SQL"
description: "Allows users to query AWS GuardDuty Publishing Destinations to retrieve information about where GuardDuty findings are published."
folder: "GuardDuty"
---

# Table: aws_guardduty_publishing_destination - Query AWS GuardDuty Publishing Destinations using SQL

The AWS GuardDuty Publishing Destination is a feature of AWS GuardDuty that allows for the continuous export of findings to an Amazon S3 bucket. This enables long-term archiving of findings and facilitates further analysis and correlation of data. It is an essential tool for maintaining security and monitoring activity within your AWS environment.

## Table Usage Guide

The `aws_guardduty_publishing_destination` table in Steampipe provides you with information about publishing destinations in AWS GuardDuty. This table allows you as a security analyst or DevOps engineer to query destination-specific details, including the destination type, status, and associated metadata. You can utilize this table to gather insights on destinations, such as their current statuses, the types of destinations, and more. The schema outlines the various attributes of the publishing destination for you, including the destination type, destination ARN, and status.

## Examples

### Basic info
Analyze the settings to understand the status and relationship between different elements in AWS GuardDuty's publishing destination. This can help in assessing the effectiveness of your current security configurations and identify any potential issues.

```sql+postgres
select
  detector_id,
  destination_id,
  arn,
  destination_arn,
  status
from
  aws_guardduty_publishing_destination;
```

```sql+sqlite
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
Determine the areas in which publishing destinations are still pending verification within the AWS GuardDuty service. This is useful to ensure all destinations are verified and secure, maintaining the integrity of your data publishing process.

```sql+postgres
select
  destination_id,
  arn,
  status
from
  aws_guardduty_publishing_destination
where
  status = 'PENDING_VERIFICATION';
```

```sql+sqlite
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
Identify instances where publishing destinations are not secured with encryption. This is useful for assessing potential security vulnerabilities within your AWS GuardDuty publishing destinations.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which AWS GuardDuty publishing destinations are most frequently used, providing a clear overview of usage patterns and aiding in resource management.

```sql+postgres
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

```sql+sqlite
select
  destination_type,
  count(destination_id)
from
  aws_guardduty_publishing_destination
group by 
  destination_type
order by
  count(destination_id) desc;
```

### Get bucket policies for S3 bucket publishing destinations
Determine the security policies associated with your S3 bucket publishing destinations. This helps you understand who has access and what actions they can perform, ensuring your data is protected and access is controlled.

```sql+postgres
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

```sql+sqlite
select
  d.destination_id,
  d.destination_arn,
  d.destination_type,
  json_extract(p.value, '$.Sid') as sid,
  json_extract(p.value, '$.Action') as policy_action,
  json_extract(p.value, '$.Effect') as effect,
  json_extract(p.value, '$.Principal.Service') as principal_service
from
  aws_guardduty_publishing_destination as d
join
  aws_s3_bucket as s
on
  s.arn = d.destination_arn
join
  json_each(s.policy, '$.Statement') as p
where
  d.destination_type = 'S3';
```

### Get KMS key policies associated with publishing destinations
This example helps identify the policies associated with the KMS keys used in publishing destinations. It's useful for auditing security settings and ensuring appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
select
  d.destination_id,
  json_extract(p.value, '$.Sid') as sid,
  json_extract(p.value, '$.Action') as policy_action,
  json_extract(p.value, '$.Effect') as effect,
  json_extract(p.value, '$.Principal') as policy_principal,
  json_extract(p.value, '$.Condition') as policy_condition
from
  aws_guardduty_publishing_destination as d,
  aws_kms_key as k,
  json_each(json_extract(k.policy, '$.Statement')) as p
where
  d.kms_key_arn is not null
and
  k.arn = d.kms_key_arn;
```