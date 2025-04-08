---
title: "Steampipe Table: aws_availability_zone - Query EC2 Availability Zones using SQL"
description: "Allows users to query EC2 Availability Zones in AWS, providing details such as zone ID, name, region, and state."
folder: "EC2"
---

# Table: aws_availability_zone - Query EC2 Availability Zones using SQL

The AWS EC2 Availability Zones are isolated locations within data center regions from which public cloud services originate and operate. They are designed to provide stable, secure, and high availability services by allowing users to run instances in several locations. These zones are an essential component for fault-tolerant and highly available infrastructure design, enabling applications to continue functioning despite a failure within a single location.

## Table Usage Guide

The `aws_availability_zone` table in Steampipe provides you with information about Availability Zones within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query zone-specific details, including zone ID, name, region, and state. You can utilize this table to gather insights on zones, such as zones that are currently available, the regions associated with each zone, and more. The schema outlines the various attributes of the Availability Zone for you, including the zone ID, zone name, region name, and zone state.

## Examples

### Availability zone info
Analyze the settings to understand the distribution and types of availability zones in different regions. This can aid in planning resource deployment for optimal performance and redundancy.

```sql+postgres
select
  name,
  zone_id,
  zone_type,
  group_name,
  region_name
from
  aws_availability_zone;
```

```sql+sqlite
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
Determine the distribution of availability zones across different regions to understand the geographical spread of your AWS resources.

```sql+postgres
select
  region_name,
  count(name) as zone_count
from
  aws_availability_zone
group by
  region_name;
```

```sql+sqlite
select
  region_name,
  count(name) as zone_count
from
  aws_availability_zone
group by
  region_name;
```


### List of AWS availability zones which are not enabled in the account
Identify the AWS availability zones that are not currently enabled within your account. This is useful for understanding which zones you may want to opt into for increased redundancy or global coverage.

```sql+postgres
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

```sql+sqlite
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