---
title: "Table: aws_rds_reserved_db_instance - Query AWS RDS Reserved DB Instances using SQL"
description: "Allows users to query RDS Reserved DB Instances in AWS, providing details such as reservation status, instance type, duration, and associated costs."
---

# Table: aws_rds_reserved_db_instance - Query AWS RDS Reserved DB Instances using SQL

The `aws_rds_reserved_db_instance` table in Steampipe provides information about Reserved DB Instances within Amazon Relational Database Service (RDS). This table allows database administrators and DevOps engineers to query detailed information about reserved instances, including the reservation identifier, instance type, duration, and associated costs. Users can utilize this table to gather insights on their reserved instances, such as understanding cost allocation, planning capacity, and managing reservations. The schema outlines the various attributes of the Reserved DB Instance, including the reservation ARN, offering type, recurring charges, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_reserved_db_instance` table, you can use the `.inspect aws_rds_reserved_db_instance` command in Steampipe.

### Key columns:

- `reserved_db_instance_id`: This is the unique identifier of the reserved DB instance. It is crucial for identifying individual reservations and can be used to join with other tables that reference the reservation.
- `db_instance_class`: This represents the compute and memory capacity of the reserved DB instance. It is useful for capacity planning and cost management.
- `reserved_db_instances_arn`: This is the Amazon Resource Name (ARN) of the reserved DB instance. It is a key identifier within AWS and can be used to join with other tables that reference the ARN.

## Examples

### Basic info

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance;
```

### List reserved DB instances with multi-AZ disabled

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance
where
  not multi_az;
```

### List reserved DB instances with offering type `All Upfront`

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance
where
  offering_type = 'All Upfront';
```

### List reserved DB instances order by duration

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance
order by
  duration desc;
```

### List reserved DB instances order by usage price

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class,
  usage_price
from
  aws_rds_reserved_db_instance
order by
  usage_price desc;
```

### List reserved DB instances which are not active

```sql
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class,
  usage_price
from
  aws_rds_reserved_db_instance
where
  state <> 'active';
```
