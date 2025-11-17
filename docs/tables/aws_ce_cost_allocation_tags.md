---
title: "Steampipe Table: aws_ce_cost_allocation_tags - Query AWS Cost Allocation Tags using SQL"
description: "Allows users to query Cost Allocation Tags from AWS Cost Explorer, providing information about tag keys, types, and active/inactive status. Useful for cost governance and tag-based cost allocation."
folder: "Cost Explorer"
---

# Table: aws_ce_cost_allocation_tags - Query AWS Cost Allocation Tags using SQL

AWS Cost Allocation Tags help you organize and track your AWS costs by adding user-defined or AWS-generated tags to your resources. Cost Allocation Tags can be activated or deactivated to control which tags are included in your cost allocation reports and Cost Explorer analysis.

## Table Usage Guide

The `aws_ce_cost_allocation_tags` table in Steampipe provides you with information about Cost Allocation Tags in your AWS account. This table allows you to query tag-specific details, including tag keys, tag types (user-defined or AWS-generated), and active/inactive status. You can utilize this table to gather insights on which tags are actively being used for cost allocation, last used dates, and tag metadata. The schema outlines the various attributes of a cost allocation tag including the tag key, type, status, and update dates.

## Examples

### List all cost allocation tags

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_updated_date
from
  aws_ce_cost_allocation_tags;
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_updated_date
from
  aws_ce_cost_allocation_tags;
```

### List active cost allocation tags

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_used_date,
  last_updated_date
from
  aws_ce_cost_allocation_tags
where
  status = 'Active';
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_used_date,
  last_updated_date
from
  aws_ce_cost_allocation_tags
where
  status = 'Active';
```

### List inactive cost allocation tags

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_used_date,
  last_updated_date
from
  aws_ce_cost_allocation_tags
where
  status = 'Inactive';
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_used_date,
  last_updated_date
from
  aws_ce_cost_allocation_tags
where
  status = 'Inactive';
```

### List user-defined cost allocation tags

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  tag_type = 'UserDefined'
order by
  tag_key;
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  tag_type = 'UserDefined'
order by
  tag_key;
```

### List AWS-generated cost allocation tags

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  tag_type = 'AWSGenerated'
order by
  tag_key;
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  tag_type = 'AWSGenerated'
order by
  tag_key;
```

### Find unused tags (not used in the last 90 days)

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  last_used_date < current_date - interval '90 days'
  or last_used_date is null
order by
  last_used_date;
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  last_used_date
from
  aws_ce_cost_allocation_tags
where
  last_used_date < date('now', '-90 days')
  or last_used_date is null
order by
  last_used_date;
```

### List tags by status and type

```sql+postgres
select
  tag_key,
  tag_type,
  status,
  count(*) as tag_count
from
  aws_ce_cost_allocation_tags
group by
  tag_type,
  status
order by
  tag_type,
  status;
```

```sql+sqlite
select
  tag_key,
  tag_type,
  status,
  count(*) as tag_count
from
  aws_ce_cost_allocation_tags
group by
  tag_type,
  status
order by
  tag_type,
  status;
```

