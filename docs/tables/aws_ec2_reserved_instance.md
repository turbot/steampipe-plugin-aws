---
title: "Table: aws_ec2_reserved_instance - Query AWS EC2 Reserved Instances using SQL"
description: "Allows users to query AWS EC2 Reserved Instances to gather comprehensive insights on the reserved instances, such as their configurations, state, and associated tags."
---

# Table: aws_ec2_reserved_instance - Query AWS EC2 Reserved Instances using SQL

The `aws_ec2_reserved_instance` table in Steampipe provides information about Reserved Instances within Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query reserved instance-specific details, including instance type, offering class, and state. Users can utilize this table to gather insights on reserved instances, such as their configurations, reserved instance state, and associated tags. The schema outlines the various attributes of the reserved instance, including its ARN, instance type, offering class, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_reserved_instance` table, you can use the `.inspect aws_ec2_reserved_instance` command in Steampipe.

**Key columns**:

- `reservation_id`: The ID of the Reserved Instance. This can be used to join this table with other tables that contain information about the reserved instance.
- `instance_type`: The instance type on which the Reserved Instance can be used. This can be used to join this table with other tables that contain information about the instance type.
- `state`: The state of the Reserved Instance. This can be used to join this table with other tables that contain information about the state of the reserved instance.

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
