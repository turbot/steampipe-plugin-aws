---
title: "Steampipe Table: aws_opensearch_reserved_instance - Query AWS OpenSearch Reserved Instances using SQL"
description: "Allows users to query AWS OpenSearch Reserved Instances and provides information about each Reserved Instance within an AWS account."
folder: "OpenSearch"
---

# Table: aws_opensearch_reserved_instance - Query AWS OpenSearch Reserved Instances using SQL

An AWS OpenSearch Reserved Instance provides a significant discount compared to On-Demand instance pricing in exchange for a commitment to use the instance for a one or three year term. This table allows you to query details about Reserved Instances that have been purchased for OpenSearch domains.

## Table Usage Guide

The `aws_opensearch_reserved_instance` table in Steampipe provides you with information about Reserved Instances within AWS OpenSearch. This table allows you, as a DevOps engineer, financial analyst, or other technical professional, to query Reserved Instance-specific details, including instance type, start time, duration, pricing information, and current state. You can utilize this table to gather insights on Reserved Instances, such as tracking reservation usage, analyzing cost savings, monitoring expiration dates, and auditing payment options. The schema outlines the various attributes of the Reserved Instance for you, including the reservation ID, instance type, duration, fixed price, usage price, and associated metadata.

## Examples

### Basic Reserved Instance info
Explore which OpenSearch Reserved Instances are active in your account, including their type, state, and pricing details.

```sql+postgres
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  state,
  instance_count,
  payment_option,
  currency_code
from
  aws_opensearch_reserved_instance;
```

```sql+sqlite
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  state,
  instance_count,
  payment_option,
  currency_code
from
  aws_opensearch_reserved_instance;
```

### List all active Reserved Instances
Discover which OpenSearch Reserved Instances are currently active and in use.

```sql+postgres
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  instance_count,
  start_time,
  duration
from
  aws_opensearch_reserved_instance
where
  state = 'active';
```

```sql+sqlite
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  instance_count,
  start_time,
  duration
from
  aws_opensearch_reserved_instance
where
  state = 'active';
```

### Calculate cost savings and analyze Reserved Instance pricing
Analyze the financial aspects of your Reserved Instances by examining upfront costs, hourly rates, and payment options.

```sql+postgres
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  payment_option,
  fixed_price,
  usage_price,
  duration,
  (fixed_price + (usage_price * duration / 3600)) as total_cost
from
  aws_opensearch_reserved_instance
order by
  total_cost desc;
```

```sql+sqlite
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  payment_option,
  fixed_price,
  usage_price,
  duration,
  (fixed_price + (usage_price * duration / 3600.0)) as total_cost
from
  aws_opensearch_reserved_instance
order by
  total_cost desc;
```

### Reserved Instances expiring soon
Identify Reserved Instances that will expire in the next 90 days so you can plan for renewals or capacity adjustments.

```sql+postgres
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  instance_count,
  start_time,
  start_time + (duration || ' seconds')::interval as end_time,
  extract(days from (start_time + (duration || ' seconds')::interval - now())) as days_until_expiration
from
  aws_opensearch_reserved_instance
where
  start_time + (duration || ' seconds')::interval <= now() + interval '90 days'
  and state = 'active'
order by
  days_until_expiration;
```

```sql+sqlite
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  instance_count,
  start_time,
  datetime(start_time, '+' || duration || ' seconds') as end_time,
  julianday(datetime(start_time, '+' || duration || ' seconds')) - julianday('now') as days_until_expiration
from
  aws_opensearch_reserved_instance
where
  datetime(start_time, '+' || duration || ' seconds') <= datetime('now', '+90 days')
  and state = 'active'
order by
  days_until_expiration;
```

### Group Reserved Instances by instance type
Analyze your Reserved Instance portfolio by grouping reservations by instance type to understand capacity distribution.

```sql+postgres
select
  instance_type,
  count(*) as reservation_count,
  sum(instance_count) as total_instances,
  sum(fixed_price) as total_upfront_cost,
  avg(usage_price) as avg_hourly_rate
from
  aws_opensearch_reserved_instance
where
  state = 'active'
group by
  instance_type
order by
  total_instances desc;
```

```sql+sqlite
select
  instance_type,
  count(*) as reservation_count,
  sum(instance_count) as total_instances,
  sum(fixed_price) as total_upfront_cost,
  avg(usage_price) as avg_hourly_rate
from
  aws_opensearch_reserved_instance
where
  state = 'active'
group by
  instance_type
order by
  total_instances desc;
```

### Analyze recurring charges
Examine the recurring charges associated with your Reserved Instances to understand ongoing costs.

```sql+postgres
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  payment_option,
  jsonb_pretty(recurring_charges) as recurring_charges_detail
from
  aws_opensearch_reserved_instance
where
  recurring_charges is not null
  and jsonb_array_length(recurring_charges) > 0;
```

```sql+sqlite
select
  reserved_instance_id,
  reservation_name,
  instance_type,
  payment_option,
  recurring_charges as recurring_charges_detail
from
  aws_opensearch_reserved_instance
where
  recurring_charges is not null
  and json_array_length(recurring_charges) > 0;
```

