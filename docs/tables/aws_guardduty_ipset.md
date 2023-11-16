---
title: "Table: aws_guardduty_ipset - Query AWS GuardDuty IPSet using SQL"
description: "Allows users to query AWS GuardDuty IPSet to retrieve information about the IPSet, such as the detector ID, IPSet ID, name, format, location, and status."
---

# Table: aws_guardduty_ipset - Query AWS GuardDuty IPSet using SQL

The `aws_guardduty_ipset` table in Steampipe provides information about IPSet within AWS GuardDuty. This table allows security analysts to query IPSet-specific details, including the detector ID, IPSet ID, name, format, location, and status. Users can utilize this table to gather insights on IPSet, such as the list of IP addresses used by GuardDuty to simulate trusted IP addresses when generating test findings. The schema outlines the various attributes of the IPSet, including the detector ID, IPSet ID, name, format, location, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_ipset` table, you can use the `.inspect aws_guardduty_ipset` command in Steampipe.

**Key columns**:

- `detector_id`: The unique ID of the detector that the IPSet is associated with. This is useful for joining this table with other tables that contain information about the detector.
- `ip_set_id`: The unique ID of the IPSet. This is useful for joining this table with other tables that contain information about the IPSet.
- `name`: The user-friendly name for the IPSet. This can be useful for human-readable queries and reports.

## Examples

### Basic info

```sql
select
  detector_id,
  ipset_id,
  name,
  format,
  location
from
  aws_guardduty_ipset;
```


### List IPSets which are not active

```sql
select
  ipset_id,
  name,
  status
from
  aws_guardduty_ipset
where
  status = 'INACTIVE';
```
