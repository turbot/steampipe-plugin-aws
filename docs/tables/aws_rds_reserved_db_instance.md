# Table: aws_rds_reserved_db_instance

Amazon RDS Reserved Instances give you the option to reserve a DB instance for a one or three year term and in turn receive a significant discount compared to the On-Demand Instance pricing for the DB instance.

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
