---
title: "Steampipe Table: aws_savings_plan - Query AWS Savings Plans using SQL"
description: "Allows users to query AWS Savings Plans to retrieve information about purchased savings plans, including commitment amounts, payment options, and coverage details."
folder: "Savings Plans"
---

# Table: aws_savings_plan - Query AWS Savings Plans using SQL

AWS Savings Plans offer a flexible pricing model that provides significant savings on AWS usage. They provide lower prices on Amazon EC2 instances usage, AWS Lambda, and AWS Fargate, regardless of instance family, size, OS, tenancy, or region. Savings Plans come in three types: Compute Savings Plans, EC2 Instance Savings Plans, and SageMaker Savings Plans.

## Table Usage Guide

The `aws_savings_plan` table in Steampipe provides you with information about AWS Savings Plans within your AWS account. This table allows you, as a DevOps engineer, cloud architect, or financial analyst, to query savings plan details, including commitment amounts, payment options, duration terms, and current states. You can utilize this table to gather insights on cost optimization, track savings plan utilization, and manage financial commitments. The schema outlines the various attributes of the savings plan, including the plan ID, type, payment details, and coverage information.

**Important Notes**
- This table supports optional quals. Queries with optional quals are optimized to reduce query time and cost. Optional quals are supported for the following columns:
  - `commitment` - Filter by commitment amount (using `=` operator)
  - `ec2_instance_family` - Filter by EC2 instance family (using `=` operator)
  - `end_time` - Filter by end time (using `<=` operator)
  - `payment_option` - Filter by payment option (using `=` operator)
  - `region` - Filter by AWS region (using `=` operator)
  - `savings_plan_type` - Filter by savings plan type (using `=` operator)
  - `start_time` - Filter by start time (using `>=` operator)
  - `state` - Filter by savings plan state
  - `term_duration_in_seconds` - Filter by term duration (using `=` operator)

## Examples

### Basic info
Explore your AWS Savings Plans to understand their current state, commitment amounts, and terms. This helps in tracking your cost optimization commitments and their effectiveness.

```sql+postgres
select
  savings_plan_id,
  arn,
  savings_plan_type,
  state,
  commitment,
  currency,
  payment_option,
  start_time,
  end_time
from
  aws_savings_plan
order by
  start_time desc;
```

```sql+sqlite
select
  savings_plan_id,
  arn,
  savings_plan_type,
  state,
  commitment,
  currency,
  payment_option,
  start_time,
  end_time
from
  aws_savings_plan
order by
  start_time desc;
```

### List active savings plans
Identify all currently active savings plans to understand your ongoing commitments and their coverage periods.

```sql+postgres
select
  savings_plan_id,
  savings_plan_type,
  commitment,
  currency,
  start_time,
  end_time,
  extract(days from (end_time - start_time)) as duration_days
from
  aws_savings_plan
where
  state = 'active'
order by
  commitment desc;
```

```sql+sqlite
select
  savings_plan_id,
  savings_plan_type,
  commitment,
  currency,
  start_time,
  end_time,
  julianday(end_time) - julianday(start_time) as duration_days
from
  aws_savings_plan
where
  state = 'active'
order by
  commitment desc;
```

### Get savings plans by payment option
Analyze your savings plans based on different payment options to understand your financial commitment structure.

```sql+postgres
select
  payment_option,
  count(*) as plan_count,
  sum(commitment::numeric) as total_commitment,
  sum(upfront_payment_amount::numeric) as total_upfront_payment,
  sum(recurring_payment_amount::numeric) as total_recurring_payment
from
  aws_savings_plan
group by
  payment_option
order by
  total_commitment desc;
```

```sql+sqlite
select
  payment_option,
  count(*) as plan_count,
  sum(cast(commitment as decimal)) as total_commitment,
  sum(cast(upfront_payment_amount as decimal)) as total_upfront_payment,
  sum(cast(recurring_payment_amount as decimal)) as total_recurring_payment
from
  aws_savings_plan
group by
  payment_option
order by
  total_commitment desc;
```

### Find savings plans nearing expiration
Identify savings plans that are approaching their end date to help with renewal planning and continued cost optimization.

```sql+postgres
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  end_time,
  extract(days from (end_time - now())) as days_until_expiration
from
  aws_savings_plan
where
  state = 'active'
  and end_time <= now() + interval '90 days'
order by
  end_time asc;
```

```sql+sqlite
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  end_time,
  julianday(end_time) - julianday('now') as days_until_expiration
from
  aws_savings_plan
where
  state = 'active'
  and julianday(end_time) <= julianday('now', '+90 days')
order by
  end_time asc;
```

### Get EC2 instance savings plans by family
Analyze EC2 Instance Savings Plans grouped by instance family to understand your compute savings strategy.

```sql+postgres
select
  ec2_instance_family,
  count(*) as plan_count,
  sum(commitment::numeric) as total_commitment,
  avg(commitment::numeric) as avg_commitment,
  string_agg(distinct state, ', ') as states
from
  aws_savings_plan
where
  savings_plan_type = 'EC2Instance'
  and ec2_instance_family is not null
group by
  ec2_instance_family
order by
  total_commitment desc;
```

```sql+sqlite
select
  ec2_instance_family,
  count(*) as plan_count,
  sum(cast(commitment as decimal)) as total_commitment,
  avg(cast(commitment as decimal)) as avg_commitment,
  group_concat(distinct state, ', ') as states
from
  aws_savings_plan
where
  savings_plan_type = 'EC2Instance'
  and ec2_instance_family is not null
group by
  ec2_instance_family
order by
  total_commitment desc;
```

### List savings plans with tags
Explore savings plans that have been tagged for better resource management and cost allocation.

```sql+postgres
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  tags
from
  aws_savings_plan
where
  tags is not null
  and jsonb_array_length(tags) > 0
order by
  savings_plan_id;
```

```sql+sqlite
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  tags
from
  aws_savings_plan
where
  tags is not null
  and json_array_length(tags) > 0
order by
  savings_plan_id;
```

### Get savings plans by region
Analyze the distribution of your savings plans across different AWS regions to understand regional cost optimization coverage.

```sql+postgres
select
  region,
  count(*) as plan_count,
  sum(commitment::numeric) as total_commitment,
  string_agg(distinct savings_plan_type, ', ') as plan_types
from
  aws_savings_plan
group by
  region
order by
  total_commitment desc;
```

```sql+sqlite
select
  region,
  count(*) as plan_count,
  sum(cast(commitment as decimal)) as total_commitment,
  group_concat(distinct savings_plan_type, ', ') as plan_types
from
  aws_savings_plan
group by
  region
order by
  total_commitment desc;
```

### Get returnable savings plans
Identify savings plans that can still be returned, which is useful for adjusting your savings commitments if needed.

```sql+postgres
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  start_time,
  returnable_until,
  extract(days from (returnable_until - now())) as days_until_return_deadline
from
  aws_savings_plan
where
  returnable_until is not null
  and returnable_until > now()
order by
  returnable_until asc;
```

```sql+sqlite
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  start_time,
  returnable_until,
  julianday(returnable_until) - julianday('now') as days_until_return_deadline
from
  aws_savings_plan
where
  returnable_until is not null
  and returnable_until > datetime('now')
order by
  returnable_until asc;
```

### Filter savings plans by time range
Query savings plans that started within a specific time period and are set to end before a certain date. This example demonstrates how to use timestamp filtering with the optional quals.

```sql+postgres
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  start_time,
  end_time,
  extract(days from (end_time - start_time)) as duration_days
from
  aws_savings_plan
where
  start_time >= '2023-01-01T00:00:00Z'
  and end_time <= '2025-12-31T23:59:59Z'
order by
  start_time desc;
```

```sql+sqlite
select
  savings_plan_id,
  savings_plan_type,
  state,
  commitment,
  currency,
  start_time,
  end_time,
  julianday(end_time) - julianday(start_time) as duration_days
from
  aws_savings_plan
where
  start_time >= '2023-01-01T00:00:00Z'
  and end_time <= '2025-12-31T23:59:59Z'
order by
  start_time desc;
```
