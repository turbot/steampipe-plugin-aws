# Table: aws_ec2_capacity_reservation

On-Demand Capacity Reservations enable you to reserve compute capacity for your Amazon EC2 instances in a specific Availability Zone for any duration. A Capacity Reservation is tied to a specific Availability Zone and, by default automatically utilized by running instances in that Availability Zone.

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
