# Table: aws_ec2_reserved_instance

Amazon EC2 Reserved Instances (RI) provide a significant discount (up to 72%) compared to On-Demand pricing and provide a capacity reservation when used in a specific Availability Zone.

## Examples

### Basic Info

```sql
select
  reserved_instance_id,
  arn,
  instance_type,
  instance_state,
  currency_code,
  CAST(fixed_price AS varchar),
  offering_class, scope,
  CAST(usage_price AS varchar)
from
  aws_ec2_reserved_instance;
```

### Count reserved instances by instance type

```sql
select
  instance_type,
  count(instance_count) as count
from
  aws_ec2_reserved_instance
group by
  instance_type;
```

### List reserved instances provisioned with undesired(for example t2.large and m3.medium is desired) instance type(s)

```sql
select
  instance_type,
  count(*) as count
from
  aws_ec2_reserved_instance
where
  instance_type not in ('t2.large', 'm3.medium')
group by
  instance_type;
```

### List standard offering class type reserved instances

```sql
select
  reserved_instance_id,
  instance_type,
  offering_class
from
  aws_ec2_reserved_instance
where
  offering_class = 'standard';
```

### List active reserved instances

```sql
select
  reserved_instance_id,
  instance_type,
  instance_state
from
  aws_ec2_reserved_instance
where
  instance_state = 'active';
```
