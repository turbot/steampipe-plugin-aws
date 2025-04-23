---
title: "Steampipe Table: aws_guardduty_threat_intel_set - Query AWS GuardDuty ThreatIntelSet using SQL"
description: "Allows users to query AWS GuardDuty ThreatIntelSet to fetch information about threat intelligence sets that are associated with a GuardDuty detector."
folder: "GuardDuty"
---

# Table: aws_guardduty_threat_intel_set - Query AWS GuardDuty ThreatIntelSet using SQL

The AWS GuardDuty ThreatIntelSet is a feature within the AWS GuardDuty service. It allows you to manage and use threat intelligence feeds that are tailored to your specific needs. This helps in identifying potential security threats and responding to them swiftly, thereby enhancing the overall security posture of your AWS environment.

## Table Usage Guide

The `aws_guardduty_threat_intel_set` table in Steampipe provides you with information about threat intelligence sets that are associated with a GuardDuty detector in AWS GuardDuty. This table allows you, as a security analyst, to query threat-specific details, including the name, format, location, and status of the threat intelligence set. You can utilize this table to gather insights on threats, such as those that are currently active, the format and location of the threat intelligence set, and more. The schema outlines the various attributes of the threat intelligence set for you, including the threat intelligence set ID, detector ID, name, format, location, status, and associated tags.

## Examples

### Basic info
Discover the segments that are being monitored for potential security threats within your AWS GuardDuty service. This allows for a better understanding of the threat intelligence sets in use and their respective configurations, aiding in effective threat management and response.

```sql+postgres
select
  detector_id,
  threat_intel_set_id,
  name,
  format,
  location
from
  aws_guardduty_threat_intel_set;
```

```sql+sqlite
select
  detector_id,
  threat_intel_set_id,
  name,
  format,
  location
from
  aws_guardduty_threat_intel_set;
```


### List disabled threat intel sets
Identify instances where the threat intelligence sets in your AWS GuardDuty service have been deactivated. This is useful for security audits or when reviewing your threat detection capabilities.

```sql+postgres
select
  threat_intel_set_id,
  status
from
  aws_guardduty_threat_intel_set
where
  status = 'INACTIVE';
```

```sql+sqlite
select
  threat_intel_set_id,
  status
from
  aws_guardduty_threat_intel_set
where
  status = 'INACTIVE';
```