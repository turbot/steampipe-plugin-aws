---
title: "Steampipe Table: aws_costoptimizationhub_recommendation - Query AWS Cost Optimization Recommendations using SQL"
description: "Allows users to query AWS Cost Optimization Hub Recommendations to obtain insights on cost-saving opportunities, resource configuration, and associated metadata."
folder: "Cost Explorer"
---

# Table: aws_costoptimizationhub_recommendation - Query AWS Cost Optimization Recommendations using SQL

The AWS Cost Optimization Hub provides recommendations for reducing costs and optimizing the usage of AWS resources. These recommendations are based on usage patterns, cost analysis, and resource configurations, helping organizations achieve better cost efficiency.

## Table Usage Guide

The `aws_costoptimizationhub_recommendation` table in Steampipe allows you to query detailed cost optimization recommendations. This table helps DevOps engineers, cost analysts, or financial professionals identify potential savings, understand implementation efforts, and track the effectiveness of recommendations.

The schema outlines various attributes of the cost optimization recommendations, including estimated savings, resource details, recommendation types, and implementation efforts. It also provides timestamps for the last refresh and additional metadata such as tags and ARNs.

**Important Notes**
- This table supports optional quals. Queries with optional quals are optimized to use additional filtering provided by the AWS API function to narrow down the results for better query performance.. Optional quals are supported for the following columns:
  - `recommendation_account_id`
  - `action_type`
  - `implementation_effort`
  - `recommendation_id`
  - `resource_region`
  - `resource_arn`
  - `resource_id`
  - `current_resource_type`
  - `recommended_resource_type`
  - `restart_needed`
  - `rollback_possible`

## Examples

### Basic info
Retrieve the basic details of cost optimization recommendations, including the resource and estimated savings.

```sql+postgres
select
  recommendation_id,
  resource_id,
  estimated_monthly_savings,
  estimated_savings_percentage,
  implementation_effort
from
  aws_costoptimizationhub_recommendation;
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  estimated_monthly_savings,
  estimated_savings_percentage,
  implementation_effort
from
  aws_costoptimizationhub_recommendation;
```

### List recommendations with significant savings
Identify recommendations where the estimated savings percentage is greater than 50%. This helps prioritize high-impact cost optimization opportunities.

```sql+postgres
select
  recommendation_id,
  resource_id,
  estimated_monthly_savings,
  estimated_savings_percentage
from
  aws_costoptimizationhub_recommendation
where
  estimated_savings_percentage > 50;
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  estimated_monthly_savings,
  estimated_savings_percentage
from
  aws_costoptimizationhub_recommendation
where
  estimated_savings_percentage > 50;
```

### List recommendations requiring a resource restart
Find recommendations that require a restart to implement. This query helps in planning implementation efforts and minimizing downtime.

```sql+postgres
select
  recommendation_id,
  resource_id,
  implementation_effort,
  restart_needed
from
  aws_costoptimizationhub_recommendation
where
  restart_needed = true;
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  implementation_effort,
  restart_needed
from
  aws_costoptimizationhub_recommendation
where
  restart_needed = 1;
```

### Get recommendations by resource type
Filter recommendations based on specific resource types, such as EC2 or RDS, to analyze opportunities for optimizing particular services.

```sql+postgres
select
  recommendation_id,
  resource_id,
  current_resource_type,
  recommended_resource_type
from
  aws_costoptimizationhub_recommendation
where
  current_resource_type = 'EC2';
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  current_resource_type,
  recommended_resource_type
from
  aws_costoptimizationhub_recommendation
where
  current_resource_type = 'EC2';
```

### List recommendations refreshed in the last 30 days
Track recently updated recommendations to stay up-to-date with the latest cost optimization insights.

```sql+postgres
select
  recommendation_id,
  resource_id,
  last_refresh_timestamp
from
  aws_costoptimizationhub_recommendation
where
  last_refresh_timestamp > now() - interval '30 days';
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  last_refresh_timestamp
from
  aws_costoptimizationhub_recommendation
where
  last_refresh_timestamp > datetime('now','-30 days');
```

### Get the tags associated with a recommendation
Retrieve tags assigned to recommendations to better organize and manage resources.

```sql+postgres
select
  recommendation_id,
  resource_id,
  jsonb_each_text(tags) as tag
from
  aws_costoptimizationhub_recommendation;
```

```sql+sqlite
select
  recommendation_id,
  resource_id,
  json_each(tags_src).key as tag_key,
  json_each(tags_src).value as tag_value
from
  aws_costoptimizationhub_recommendation;
```

