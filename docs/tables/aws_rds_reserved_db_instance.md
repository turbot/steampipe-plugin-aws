---
title: "Steampipe Table: aws_rds_reserved_db_instance - Query AWS RDS Reserved DB Instances using SQL"
description: "Allows users to query RDS Reserved DB Instances in AWS, providing details such as reservation status, instance type, duration, and associated costs."
folder: "RDS"
---

# Table: aws_rds_reserved_db_instance - Query AWS RDS Reserved DB Instances using SQL

The AWS RDS Reserved DB Instance is a reservation of resources within the Amazon Relational Database Service. This reservation allows you to receive a significant discount, compared to the on-demand instance pricing, by committing to a one or three year term for using DB Instances. It provides a cost-effective solution to reserve capacity for your Amazon RDS instances in a specific region.

## Table Usage Guide

The `aws_rds_reserved_db_instance` table in Steampipe provides you with information about Reserved DB Instances within Amazon Relational Database Service (RDS). This table enables you, as a database administrator or DevOps engineer, to query detailed information about reserved instances, including the reservation identifier, instance type, duration, and associated costs. You can utilize this table to gather insights on your reserved instances, such as understanding cost allocation, planning capacity, and managing reservations. The schema outlines the various attributes of the Reserved DB Instance for you, including the reservation ARN, offering type, recurring charges, and associated tags.

## Examples

### Basic info
Explore the status and details of your reserved database instances in AWS RDS to better manage your database resources and costs.

```sql+postgres
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance;
```

```sql+sqlite
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
Identify instances where the database instances reserved on AWS RDS do not have the multi-Availability Zone feature enabled. This is useful to pinpoint potential areas of vulnerability in your database architecture, as disabling multi-AZ can impact data redundancy and failover support.

```sql+postgres
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

```sql+sqlite
select
  reserved_db_instance_id,
  arn,
  reserved_db_instances_offering_id,
  state,
  class
from
  aws_rds_reserved_db_instance
where
  multi_az = 0;
```

### List reserved DB instances with offering type `All Upfront`
Identify instances where database reservations have been made with an 'All Upfront' offering type in AWS RDS. This can be useful for cost analysis and budgeting as it highlights where upfront payments have been made for database services.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have reserved database instances, organized by their duration. This can be particularly useful for prioritizing and managing resources effectively within your AWS RDS environment.

```sql+postgres
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

```sql+sqlite
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
Identify the most expensive reserved database instances in your AWS RDS service. This can help prioritize cost management efforts and optimize your cloud resource allocation.

```sql+postgres
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

```sql+sqlite
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
Explore which reserved database instances are not currently active. This can help in managing resources and costs by identifying unused or underutilized instances.

```sql+postgres
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

```sql+sqlite
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
  state != 'active';
```