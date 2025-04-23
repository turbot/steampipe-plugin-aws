---
title: "Steampipe Table: aws_ec2_reserved_instance - Query AWS EC2 Reserved Instances using SQL"
description: "Allows users to query AWS EC2 Reserved Instances to gather comprehensive insights on the reserved instances, such as their configurations, state, and associated tags."
folder: "EC2"
---

# Table: aws_ec2_reserved_instance - Query AWS EC2 Reserved Instances using SQL

The AWS EC2 Reserved Instances are a type of Amazon EC2 instance that allows you to reserve compute capacity for your AWS account in a specific Availability Zone, providing a significant discount compared to On-Demand pricing. These instances are recommended for applications with steady state usage, offering up to 75% savings compared to on-demand instances. AWS EC2 Reserved Instances can be purchased with a one-time payment and used throughout the term you select.

## Table Usage Guide

The `aws_ec2_reserved_instance` table in Steampipe provides you with information about Reserved Instances within Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query reserved instance-specific details, including instance type, offering class, and state. You can utilize this table to gather insights on reserved instances, such as their configurations, reserved instance state, and associated tags. The schema outlines the various attributes of the reserved instance for you, including its ARN, instance type, offering class, and associated tags.

## Examples

### Basic Info
Determine the areas in which you can gain insights into your Amazon EC2 reserved instances, such as understanding the instance type, state, and costs associated with the reservation. This is useful for managing your AWS resources and optimizing your cloud cost.

```sql+postgres
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

```sql+sqlite
select
  reserved_instance_id,
  arn,
  instance_type,
  instance_state,
  currency_code,
  CAST(fixed_price AS text),
  offering_class, scope,
  CAST(usage_price AS text)
from
  aws_ec2_reserved_instance;
```

### Count reserved instances by instance type
Determine the number of reserved instances per type to better manage your AWS EC2 resources and optimize your cloud infrastructure.

```sql+postgres
select
  instance_type,
  count(instance_count) as count
from
  aws_ec2_reserved_instance
group by
  instance_type;
```

```sql+sqlite
select
  instance_type,
  count(instance_count) as count
from
  aws_ec2_reserved_instance
group by
  instance_type;
```

### List reserved instances provisioned with undesired(for example t2.large and m3.medium is desired) instance type(s)
Determine the areas in which the provisioned reserved instances are not of the desired types such as t2.large and m3.medium. This can help in optimizing resources and better cost management.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that consist of standard offering class type within reserved instances in AWS EC2, which can assist in better management of resource allocation and cost optimization.

```sql+postgres
select
  reserved_instance_id,
  instance_type,
  offering_class
from
  aws_ec2_reserved_instance
where
  offering_class = 'standard';
```

```sql+sqlite
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
Determine the areas in which active reserved instances are being utilized within your AWS EC2 service. This can help in managing resources and optimizing costs by identifying instances that are currently in active use.

```sql+postgres
select
  reserved_instance_id,
  instance_type,
  instance_state
from
  aws_ec2_reserved_instance
where
  instance_state = 'active';
```

```sql+sqlite
select
  reserved_instance_id,
  instance_type,
  instance_state
from
  aws_ec2_reserved_instance
where
  instance_state = 'active';
```