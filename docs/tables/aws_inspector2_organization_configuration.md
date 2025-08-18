---
title: "Steampipe Table: aws_inspector2_organization_configuration - Query AWS Inspector2 Organization Configuration using SQL"
description: "Allows users to query AWS Inspector2 Organization Configuration and retrieve detailed information about organization-level settings for Amazon Inspector."
folder: "Inspector2"
---

# Table: aws_inspector2_organization_configuration - Query AWS Inspector2 Organization Configuration using SQL

AWS Inspector2 Organization Configuration provides organization-level settings for Amazon Inspector, including which scan types are automatically enabled for new members and whether the organization has reached the maximum account limit.

## Table Usage Guide

The `aws_inspector2_organization_configuration` table in Steampipe provides you with information about organization-level settings for AWS Inspector2. This table allows you, as a DevOps engineer, to query configuration details including auto-enable settings for different scan types and account limit status. You can utilize this table to gather insights on how Inspector2 is configured at the organization level, such as which scan types are automatically enabled for new members and whether the organization has reached its account limit.

## Examples

### Basic info
Explore the organization configuration for AWS Inspector2 to understand which scan types are automatically enabled for new members and whether the organization has reached its account limit.

```sql+postgres
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached
from
  aws_inspector2_organization_configuration;
```

### Check which scan types are auto-enabled
Identify which scan types are automatically enabled for new members of your organization. This helps in understanding the default security posture for new accounts.

```sql+postgres
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable
from
  aws_inspector2_organization_configuration;
```

### Check if organization has reached account limit
Determine whether your organization has reached the maximum AWS account limit for Amazon Inspector. This is important for capacity planning and understanding organizational constraints.

```sql+postgres
select
  max_account_limit_reached,
  case
    when max_account_limit_reached then 'Organization has reached maximum account limit'
    else 'Organization can add more accounts'
  end as limit_status
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  max_account_limit_reached,
  case
    when max_account_limit_reached then 'Organization has reached maximum account limit'
    else 'Organization can add more accounts'
  end as limit_status
from
  aws_inspector2_organization_configuration;
```

### Find organizations with all scan types enabled
Identify organizations that have all scan types (EC2, ECR, Lambda, and Lambda Code) automatically enabled for new members. This indicates a comprehensive security posture.

```sql+postgres
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when ec2_auto_enable
      and ecr_auto_enable
      and lambda_auto_enable
      and lambda_code_auto_enable
    then 'All scan types enabled'
    else 'Some scan types disabled'
  end as scan_coverage
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when ec2_auto_enable
      and ecr_auto_enable
      and lambda_auto_enable
      and lambda_code_auto_enable
    then 'All scan types enabled'
    else 'Some scan types disabled'
  end as scan_coverage
from
  aws_inspector2_organization_configuration;
```

### Check EC2 scan configuration
Focus specifically on EC2 scan configuration to understand if EC2 scans are automatically enabled for new members.

```sql+postgres
select
  ec2_auto_enable,
  case
    when ec2_auto_enable then 'EC2 scans are auto-enabled'
    else 'EC2 scans are not auto-enabled'
  end as ec2_status
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  case
    when ec2_auto_enable then 'EC2 scans are auto-enabled'
    else 'EC2 scans are not auto-enabled'
  end as ec2_status
from
  aws_inspector2_organization_configuration;
```

### Check ECR scan configuration
Focus specifically on ECR scan configuration to understand if ECR scans are automatically enabled for new members.

```sql+postgres
select
  ecr_auto_enable,
  case
    when ecr_auto_enable then 'ECR scans are auto-enabled'
    else 'ECR scans are not auto-enabled'
  end as ecr_status
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ecr_auto_enable,
  case
    when ecr_auto_enable then 'ECR scans are auto-enabled'
    else 'ECR scans are not auto-enabled'
  end as ecr_status
from
  aws_inspector2_organization_configuration;
```

### Check Lambda scan configuration
Focus specifically on Lambda scan configuration to understand if Lambda scans are automatically enabled for new members.

```sql+postgres
select
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when lambda_auto_enable and lambda_code_auto_enable
    then 'Both Lambda and Lambda Code scans are enabled'
    when lambda_auto_enable
    then 'Only Lambda scans are enabled'
    when lambda_code_auto_enable
    then 'Only Lambda Code scans are enabled'
    else 'No Lambda scans are enabled'
  end as lambda_scan_status
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when lambda_auto_enable and lambda_code_auto_enable
    then 'Both Lambda and Lambda Code scans are enabled'
    when lambda_auto_enable
    then 'Only Lambda scans are enabled'
    when lambda_code_auto_enable
    then 'Only Lambda Code scans are enabled'
    else 'No Lambda scans are enabled'
  end as lambda_scan_status
from
  aws_inspector2_organization_configuration;
```

### Get complete organization configuration
Retrieve the complete organization configuration including all auto-enable settings and account limit status for comprehensive analysis.

```sql+postgres
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  region,
  title
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  region,
  title
from
  aws_inspector2_organization_configuration;
```

### Find organizations with disabled scan types
Identify which scan types are not automatically enabled for new members. This helps in understanding potential security gaps in the organization's default configuration.

```sql+postgres
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when not ec2_auto_enable then 'EC2 scans disabled'
    when not ecr_auto_enable then 'ECR scans disabled'
    when not lambda_auto_enable then 'Lambda scans disabled'
    when not lambda_code_auto_enable then 'Lambda Code scans disabled'
    else 'All scan types enabled'
  end as disabled_scan_types
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  case
    when not ec2_auto_enable then 'EC2 scans disabled'
    when not ecr_auto_enable then 'ECR scans disabled'
    when not lambda_auto_enable then 'Lambda scans disabled'
    when not lambda_code_auto_enable then 'Lambda Code scans disabled'
    else 'All scan types enabled'
  end as disabled_scan_types
from
  aws_inspector2_organization_configuration;
```

### Count enabled scan types per organization
Count how many scan types are automatically enabled for new members to understand the organization's security coverage level.

```sql+postgres
select
  (ec2_auto_enable::int + ecr_auto_enable::int + lambda_auto_enable::int + lambda_code_auto_enable::int) as enabled_scan_count,
  case
    when (ec2_auto_enable::int + ecr_auto_enable::int + lambda_auto_enable::int + lambda_code_auto_enable::int) = 4 then 'Full coverage'
    when (ec2_auto_enable::int + ecr_auto_enable::int + lambda_auto_enable::int + lambda_code_auto_enable::int) >= 2 then 'Partial coverage'
    else 'Minimal coverage'
  end as coverage_level
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  (ec2_auto_enable + ecr_auto_enable + lambda_auto_enable + lambda_code_auto_enable) as enabled_scan_count,
  case
    when (ec2_auto_enable + ecr_auto_enable + lambda_auto_enable + lambda_code_auto_enable) = 4 then 'Full coverage'
    when (ec2_auto_enable + ecr_auto_enable + lambda_auto_enable + lambda_code_auto_enable) >= 2 then 'Partial coverage'
    else 'Minimal coverage'
  end as coverage_level
from
  aws_inspector2_organization_configuration;
```