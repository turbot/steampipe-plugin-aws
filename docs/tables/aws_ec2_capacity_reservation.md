---
title: "Table: aws_ec2_capacity_reservation - Query AWS EC2 Capacity Reservations using SQL"
description: "Allows users to query AWS EC2 Capacity Reservations to provide information about the reservations within AWS Elastic Compute Cloud (EC2)."
---

# Table: aws_ec2_capacity_reservation - Query AWS EC2 Capacity Reservations using SQL

The `aws_ec2_capacity_reservation` table in Steampipe provides information about Capacity Reservations within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query reservation-specific details, including reservation ID, reservation ARN, state, instance type, and associated metadata. Users can utilize this table to gather insights on reservations, such as reservations per availability zone, reservations per instance type, and more. The schema outlines the various attributes of the EC2 Capacity Reservation, including the reservation ID, creation date, instance count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_capacity_reservation` table, you can use the `.inspect aws_ec2_capacity_reservation` command in Steampipe.

Key columns:

- `reservation_id`: This is the unique identifier for the reservation. It can be used to join this table with other tables that contain reservation-specific information.
- `availability_zone`: This column provides information about the availability zone associated with the reservation. It can be used to join this table with other tables that contain availability zone-specific information.
- `instance_type`: This column provides information about the instance type associated with the reservation. It can be used to join this table with other tables that contain instance type-specific information.

## Examples

### Basic info

```sql
select
  capacity_reservation_id,
  capacity_reservation_arn,
  instance_type,
  state
from
  aws_ec2_capacity_reservation;
```

### List EC2 expired capacity reservations

```sql
select
  capacity_reservation_id,
  capacity_reservation_arn,
  instance_type,
  state
from
  aws_ec2_capacity_reservation
where
  state = 'expired';
```

### Get EC2 capacity reservation by ID

```sql
select
  capacity_reservation_id,
  capacity_reservation_arn,
  instance_type,
  state
from
  aws_ec2_capacity_reservation
where
  capacity_reservation_id = 'cr-0b30935e9fc2da81e';
```
