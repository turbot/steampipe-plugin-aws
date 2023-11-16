---
title: "Table: aws_guardduty_threat_intel_set - Query AWS GuardDuty ThreatIntelSet using SQL"
description: "Allows users to query AWS GuardDuty ThreatIntelSet to fetch information about threat intelligence sets that are associated with a GuardDuty detector."
---

# Table: aws_guardduty_threat_intel_set - Query AWS GuardDuty ThreatIntelSet using SQL

The `aws_guardduty_threat_intel_set` table in Steampipe provides information about threat intelligence sets that are associated with a GuardDuty detector in AWS GuardDuty. This table allows security analysts to query threat-specific details, including the name, format, location, and status of the threat intelligence set. Users can utilize this table to gather insights on threats, such as those that are currently active, the format and location of the threat intelligence set, and more. The schema outlines the various attributes of the threat intelligence set, including the threat intelligence set ID, detector ID, name, format, location, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_threat_intel_set` table, you can use the `.inspect aws_guardduty_threat_intel_set` command in Steampipe.

**Key columns**:

- `threat_intel_set_id`: This is the unique identifier for the threat intelligence set. It can be used to join this table with other tables that contain information about specific threat intelligence sets.
- `detector_id`: This is the unique identifier for the GuardDuty detector associated with the threat intelligence set. It can be used to join this table with other tables that contain information about specific detectors.
- `name`: This is the name of the threat intelligence set. It can be used to join this table with other tables that contain information about specific threat intelligence sets, allowing for easy identification of the sets.

## Examples

### Basic info

```sql
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

```sql
select
  threat_intel_set_id,
  status
from
  aws_guardduty_threat_intel_set
where
  status = 'INACTIVE';
```
