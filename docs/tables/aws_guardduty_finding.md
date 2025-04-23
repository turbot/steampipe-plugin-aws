---
title: "Steampipe Table: aws_guardduty_finding - Query AWS GuardDuty Findings using SQL"
description: "Allows users to query AWS GuardDuty Findings to access detailed information about potential security threats or suspicious activities detected in their AWS environment."
folder: "GuardDuty"
---

# Table: aws_guardduty_finding - Query AWS GuardDuty Findings using SQL

AWS GuardDuty is a threat detection service that continuously monitors for malicious or unauthorized behavior to help protect your AWS accounts and workloads. It identifies unusual or unauthorized activity, like crypto-currency mining or infrastructure deployments in a region that has never been used. GuardDuty analyzes tens of billions of events across multiple AWS data sources, such as AWS CloudTrail, Amazon VPC Flow Logs, and DNS logs.

## Table Usage Guide

The `aws_guardduty_finding` table in Steampipe provides you with information about findings reported by AWS GuardDuty. This table allows you as a security analyst to query finding-specific details, including threat type, severity, and associated resources. You can utilize this table to gather insights on potential security threats, such as unauthorized access attempts, data breaches, or compromised instances. The schema outlines the various attributes of the GuardDuty finding for you, including the finding ID, detector ID, account ID, region, and associated tags.

## Examples

### Basic info
Explore which instances have been identified by AWS GuardDuty. This is useful for assessing the security findings and understanding when they were created.

```sql+postgres
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding;
```

```sql+sqlite
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding;
```

### List findings that are not archived
Discover the segments that consist of unarchived findings in your AWS GuardDuty. This is particularly useful in identifying active threats or issues that are yet to be addressed and archived.

```sql+postgres
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

```sql+sqlite
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding
where
  json_extract(service, '$.Archived') = 'false';
```