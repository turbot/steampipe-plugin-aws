---
title: "Steampipe Table: aws_trusted_advisor_check_result - Query AWS Trusted Advisor Check Results using SQL"
description: "Allows users to query AWS Trusted Advisor Check Results, providing detailed information about individual flagged resources identified by Trusted Advisor checks."
folder: "Support"
---

# Table: aws_trusted_advisor_check_result - Query AWS Trusted Advisor Check Results using SQL

AWS Trusted Advisor is a real-time guidance tool that helps you provision your resources following AWS best practices. It inspects your AWS environment and makes recommendations for saving money, improving system performance and reliability, and closing security gaps. The `aws_trusted_advisor_check_result` table provides information about individual flagged resources identified by Trusted Advisor checks.

## Table Usage Guide

The `aws_trusted_advisor_check_result` table in Steampipe provides you with information about individual flagged resources from AWS Trusted Advisor checks. This table allows you, as a DevOps engineer, cloud administrator, or security analyst, to query resource-specific details identified by Trusted Advisor, including resource IDs, statuses, regions, and metadata. You can utilize this table to gather insights on flagged resources, such as identifying critical issues, analyzing resource distribution across regions, and tracking resource-level recommendations. The schema outlines the various attributes of flagged resources, including the check ID, resource information, suppression status, and associated metadata.

**Important Notes**
- You must specify both `language` and `check_id` in the WHERE clause to query this table.
- Each row represents an individual flagged resource rather than a check summary.
- You must have a Business, Enterprise On-Ramp, or Enterprise Support plan to use AWS Trusted Advisor.
- Amazon Web Services Support API currently supports the following languages for Trusted Advisor:
  - Chinese, Simplified - zh
  - Chinese, Traditional - zh_TW
  - English - en
  - French - fr
  - German - de
  - Indonesian - id
  - Italian - it
  - Japanese - ja
  - Korean - ko
  - Portuguese, Brazilian - pt_BR
  - Spanish - es

## Examples

### Basic info
Retrieve fundamental information about flagged resources for a specific Trusted Advisor check, including resource IDs, statuses, and regions.

```sql+postgres
select
  check_id,
  flagged_resource_id,
  flagged_resource_status,
  flagged_resource_region,
  flagged_resource_is_suppressed
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5';
```

```sql+sqlite
select
  check_id,
  flagged_resource_id,
  flagged_resource_status,
  flagged_resource_region,
  flagged_resource_is_suppressed
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5';
```

### Get error status resources only
Identify resources that are flagged with an error status, indicating critical issues that need immediate attention.

```sql+postgres
select
  check_id,
  flagged_resource_id,
  flagged_resource_status,
  flagged_resource_region,
  flagged_resource_metadata
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and flagged_resource_status = 'error';
```

```sql+sqlite
select
  check_id,
  flagged_resource_id,
  flagged_resource_status,
  flagged_resource_region,
  flagged_resource_metadata
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and flagged_resource_status = 'error';
```

### Count flagged resources by status
Get a summary count of flagged resources grouped by their status to understand the distribution of issues.

```sql+postgres
select
  flagged_resource_status,
  count(*) as resource_count
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  flagged_resource_status
order by
  resource_count desc;
```

```sql+sqlite
select
  flagged_resource_status,
  count(*) as resource_count
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  flagged_resource_status
order by
  resource_count desc;
```

### List flagged resources by region
Analyze the regional distribution of flagged resources to identify regions with the most issues.

```sql+postgres
select
  flagged_resource_region,
  flagged_resource_status,
  count(*) as resource_count
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  flagged_resource_region,
  flagged_resource_status
order by
  flagged_resource_region,
  resource_count desc;
```

```sql+sqlite
select
  flagged_resource_region,
  flagged_resource_status,
  count(*) as resource_count
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  flagged_resource_region,
  flagged_resource_status
order by
  flagged_resource_region,
  resource_count desc;
```

### Get suppressed resources
Identify resources that have been suppressed by users, indicating issues that have been acknowledged but not resolved.

```sql+postgres
select
  check_id,
  flagged_resource_id,
  flagged_resource_region,
  flagged_resource_status,
  flagged_resource_metadata
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and flagged_resource_is_suppressed = true;
```

```sql+sqlite
select
  check_id,
  flagged_resource_id,
  flagged_resource_region,
  flagged_resource_status,
  flagged_resource_metadata
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and flagged_resource_is_suppressed = 1;
```

### View check summary information
Get an overview of the overall check status and resource summary statistics along with individual flagged resources.

```sql+postgres
select
  check_id,
  status as check_status,
  timestamp as check_timestamp,
  resources_summary,
  category_specific_summary,
  count(*) as total_flagged_resources
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  check_id,
  status,
  timestamp,
  resources_summary,
  category_specific_summary;
```

```sql+sqlite
select
  check_id,
  status as check_status,
  timestamp as check_timestamp,
  resources_summary,
  category_specific_summary,
  count(*) as total_flagged_resources
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
group by
  check_id,
  status,
  timestamp,
  resources_summary,
  category_specific_summary;
```

### Analyze cost optimization opportunities
For cost optimization checks, examine the category-specific summary along with flagged resources to understand potential savings.

```sql+postgres
select
  check_id,
  flagged_resource_id,
  flagged_resource_metadata,
  category_specific_summary
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and category_specific_summary is not null;
```

```sql+sqlite
select
  check_id,
  flagged_resource_id,
  flagged_resource_metadata,
  category_specific_summary
from
  aws_trusted_advisor_check_result
where
  language = 'en'
  and check_id = 'L4dfs2Q4C5'
  and category_specific_summary is not null;
```
