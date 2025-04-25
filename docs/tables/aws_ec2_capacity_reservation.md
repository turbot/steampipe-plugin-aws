---
title: "Steampipe Table: aws_ec2_capacity_reservation - Query AWS EC2 Capacity Reservations using SQL"
description: "Allows users to query AWS EC2 Capacity Reservations to provide information about the reservations within AWS Elastic Compute Cloud (EC2)."
folder: "EC2"
---

# Table: aws_ec2_capacity_reservation - Query AWS EC2 Capacity Reservations using SQL

An AWS EC2 Capacity Reservation ensures that you have reserved capacity for your Amazon EC2 instances in a specific Availability Zone for any duration. This capacity reservation helps to reduce the risks of insufficient capacity for launching instances into an Availability Zone, providing predictable instance launch times. It's a useful tool for capacity planning and managing costs, particularly for applications with predictable peaks in demand.

## Table Usage Guide

The `aws_ec2_capacity_reservation` table in Steampipe provides you with information about Capacity Reservations within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query reservation-specific details, including reservation ID, reservation ARN, state, instance type, and associated metadata. You can utilize this table to gather insights on reservations, such as reservations per availability zone, reservations per instance type, and more. The schema outlines for you the various attributes of the EC2 Capacity Reservation, including the reservation ID, creation date, instance count, and associated tags.

## Examples

### Basic info
Identify instances where Amazon EC2 capacity reservations are made, gaining insights into the type of instances reserved and their current state. This can help in efficiently managing resources and understanding reservation patterns.

```sql+postgres
select
  capacity_reservation_id,
  capacity_reservation_arn,
  instance_type,
  state
from
  aws_ec2_capacity_reservation;
```

```sql+sqlite
select
  capacity_reservation_id,
  capacity_reservation_arn,
  instance_type,
  state
from
  aws_ec2_capacity_reservation;
```

### List EC2 expired capacity reservations
Identify instances where Amazon EC2 capacity reservations have expired. This information can be useful in managing resources and potentially freeing up unused capacity.

```sql+postgres
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

```sql+sqlite
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
Determine the status and type of a specific EC2 capacity reservation in AWS, which can be useful for managing and optimizing resource allocation.

```sql+postgres
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

```sql+sqlite
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