---
title: "Table: aws_availability_zone - Query EC2 Availability Zones using SQL"
description: "Allows users to query EC2 Availability Zones in AWS, providing details such as zone ID, name, region, and state."
---

# Table: aws_availability_zone - Query EC2 Availability Zones using SQL

The `aws_availability_zone` table in Steampipe provides information about Availability Zones within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query zone-specific details, including zone ID, name, region, and state. Users can utilize this table to gather insights on zones, such as zones that are currently available, the regions associated with each zone, and more. The schema outlines the various attributes of the Availability Zone, including the zone ID, zone name, region name, and zone state.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_availability_zone` table, you can use the `.inspect aws_availability_zone` command in Steampipe.

Key columns:

- `zone_id`: The ID of the availability zone. This can be used to join this table with other tables that contain availability zone information.
- `zone_name`: The name of the availability zone. This is useful for identifying the specific zone within a region.
- `region_name`: The name of the region that the availability zone is located in. This can be used to join this table with other tables that contain region-specific information.

## Examples

### Availability zone info

```sql
select
  name,
  zone_id,
  zone_type,
  group_name,
  region_name
from
  aws_availability_zone;
```


### Count of availability zone per region

```sql
select
  region_name,
  count(name) as zone_count
from
  aws_availability_zone
group by
  region_name;
```


### List of AWS availability zones which are not enabled in the account

```sql
select
  name,
  zone_id,
  region_name,
  opt_in_status
from
  aws_availability_zone
where
  opt_in_status = 'not-opted-in';
```
