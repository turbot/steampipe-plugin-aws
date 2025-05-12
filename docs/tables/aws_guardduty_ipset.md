---
title: "Steampipe Table: aws_guardduty_ipset - Query AWS GuardDuty IPSet using SQL"
description: "Allows users to query AWS GuardDuty IPSet to retrieve information about the IPSet, such as the detector ID, IPSet ID, name, format, location, and status."
folder: "GuardDuty"
---

# Table: aws_guardduty_ipset - Query AWS GuardDuty IPSet using SQL

The AWS GuardDuty IPSet is a component of Amazon GuardDuty that allows you to manage and use lists of trusted IP addresses. It can help you to more effectively detect and respond to potential security threats by defining IP conditions and filtering findings. This contributes to the overall security and integrity of your AWS environment by providing an additional layer of protection against unauthorized or malicious activity.

## Table Usage Guide

The `aws_guardduty_ipset` table in Steampipe provides you with information about IPSet within AWS GuardDuty. This table allows you, as a security analyst, to query IPSet-specific details, including the detector ID, IPSet ID, name, format, location, and status. You can utilize this table to gather insights on IPSet, such as the list of IP addresses used by GuardDuty to simulate trusted IP addresses when generating test findings. The schema outlines the various attributes of the IPSet for you, including the detector ID, IPSet ID, name, format, location, and status.

## Examples

### Basic info
Determine the areas in which potential security threats can be identified within the AWS GuardDuty service. This query is useful for gaining insights into the specific locations and formats of these threats, helping to enhance your overall security posture.

```sql+postgres
select
  detector_id,
  ipset_id,
  name,
  format,
  location
from
  aws_guardduty_ipset;
```

```sql+sqlite
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
Identify instances where certain IP sets within the AWS GuardDuty service are inactive. This is useful for maintaining network security by ensuring all necessary IP sets are active and functioning as expected.

```sql+postgres
select
  ipset_id,
  name,
  status
from
  aws_guardduty_ipset
where
  status = 'INACTIVE';
```

```sql+sqlite
select
  ipset_id,
  name,
  status
from
  aws_guardduty_ipset
where
  status = 'INACTIVE';
```