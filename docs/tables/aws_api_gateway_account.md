---
title: "Steampipe Table: aws_api_gateway_account - Query AWS API Gateway Account using SQL"
description: "Allows users to query Account settings in AWS API Gateway. The `aws_api_gateway_account` table in Steampipe provides information about Account settings within AWS API Gateway. This table allows DevOps engineers to query Account-specific details, including throttle settings, CloudWatch role ARN, API key version, and supported features. Users can utilize this table to gather insights on Account settings, such as throttle limits, monitoring configuration, and feature availability. The schema outlines the various attributes of the Account, including the CloudWatch role ARN, throttle settings, API key version, and supported features."
folder: "API Gateway"
---

# Table: aws_api_gateway_account - Query AWS API Gateway Account using SQL

AWS API Gateway Account represents the account-level settings for Amazon API Gateway in a specific region. It contains configuration information such as throttle settings, CloudWatch integration role, API key version, and supported features. These settings apply globally to all APIs within the account and region.

## Table Usage Guide

The `aws_api_gateway_account` table in Steampipe provides you with information about Account settings within AWS API Gateway. This table allows you, as a DevOps engineer, to query Account-specific details, including throttle settings, CloudWatch role ARN, API key version, and supported features. You can utilize this table to gather insights on Account settings, such as throttle limits, monitoring configuration, and feature availability. The schema outlines the various attributes of the Account for you, including the CloudWatch role ARN, throttle settings, API key version, and supported features.

## Examples

### Basic info
Discover the segments that show the configuration of your AWS API Gateway account settings. This query can provide insights into throttle limits, monitoring setup, and feature availability which can be beneficial for optimizing API performance and monitoring.

```sql+postgres
select
  cloudwatch_role_arn,
  api_key_version,
  throttle_burst_limit,
  throttle_rate_limit,
  features,
  region,
  account_id
from
  aws_api_gateway_account;
```

```sql+sqlite
select
  cloudwatch_role_arn,
  api_key_version,
  throttle_burst_limit,
  throttle_rate_limit,
  features,
  region,
  account_id
from
  aws_api_gateway_account;
```

### Check if CloudWatch logging is configured
Determine if CloudWatch logging is properly configured for API Gateway by checking if a CloudWatch role ARN is set.

```sql+postgres
select
  region,
  cloudwatch_role_arn,
  case
    when cloudwatch_role_arn is not null then 'Configured'
    else 'Not Configured'
  end as cloudwatch_status
from
  aws_api_gateway_account;
```

```sql+sqlite
select
  region,
  cloudwatch_role_arn,
  case
    when cloudwatch_role_arn is not null then 'Configured'
    else 'Not Configured'
  end as cloudwatch_status
from
  aws_api_gateway_account;
```

### View throttle settings across regions
Assess the throttle configuration across different regions to ensure consistent API rate limiting policies.

```sql+postgres
select
  region,
  throttle_rate_limit,
  throttle_burst_limit,
  case
    when throttle_rate_limit > 0 then 'Throttling Enabled'
    else 'Throttling Disabled'
  end as throttle_status
from
  aws_api_gateway_account
order by
  region;
```

```sql+sqlite
select
  region,
  throttle_rate_limit,
  throttle_burst_limit,
  case
    when throttle_rate_limit > 0 then 'Throttling Enabled'
    else 'Throttling Disabled'
  end as throttle_status
from
  aws_api_gateway_account
order by
  region;
```

### Check supported features
Identify which features are supported in each region, particularly usage plans functionality.

```sql+postgres
select
  region,
  features,
  case
    when features ? 'UsagePlans' then 'Supported'
    else 'Not Supported'
  end as usage_plans_support
from
  aws_api_gateway_account;
```

```sql+sqlite
select
  region,
  features,
  case
    when json_extract(features, '$[0]') = 'UsagePlans' then 'Supported'
    else 'Not Supported'
  end as usage_plans_support
from
  aws_api_gateway_account;
```
