---
title: "Steampipe Table: aws_inspector2_organization_configuration - Query AWS Inspector2 Regional Organization Configuration using SQL"
description: "Allows users to query AWS Inspector2 regional organization configuration, including auto-enable settings for scan types and account limit status per region."
folder: "Inspector2"
---

# Table: aws_inspector2_organization_configuration - Query AWS Inspector2 Regional Organization Configuration using SQL

The AWS Inspector2 Regional Organization Configuration contains settings that determine which scan types are automatically enabled for new members of your organization within a specific region and whether the organization has reached its account limit. These configurations help manage the security posture across your AWS organization on a regional basis.

## Table Usage Guide

The `aws_inspector2_organization_configuration` table in Steampipe provides you with information about the regional organization configuration of Amazon Inspector2. This table allows you, as a cloud administrator, security team member, or compliance officer, to query regional organization settings, including which scan types are automatically enabled for new members and whether the organization has reached its account limit. You can utilize this table to gather insights on regional organization configuration, such as EC2, ECR, Lambda, and Lambda Code scan auto-enablement status, account limit status, and region information. The schema outlines the various attributes of the regional organization configuration for you, including the region, scan type auto-enablement settings, and account limit status.

**Important Notes**
- To query this table, the account must be registered as the delegated administrator. For more details, see: https://docs.aws.amazon.com/inspector/latest/user/admin-member-relationship.html 

## Examples

### Basic info
Analyze the regional organization configuration to understand which scan types are automatically enabled for new members and whether the organization has reached its account limit. This is useful for ensuring your security posture is properly configured across regions.

```sql+postgres
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  title
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  title
from
  aws_inspector2_organization_configuration;
```

### List regions with scan type auto-enablement settings
Identify which scan types are automatically enabled for new members of your organization across all regions. This helps in understanding the default security posture for new accounts.

```sql+postgres
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable
from
  aws_inspector2_organization_configuration;
```

### List regions with organization account limit status
Determine whether your organization has reached the maximum AWS account limit for Amazon Inspector across all regions. This is important for capacity planning and understanding organizational constraints.

```sql+postgres
select
  region,
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
  region,
  max_account_limit_reached,
  case
    when max_account_limit_reached then 'Organization has reached maximum account limit'
    else 'Organization can add more accounts'
  end as limit_status
from
  aws_inspector2_organization_configuration;
```

### List regions with comprehensive scan type coverage
Identify regions that have all scan types (EC2, ECR, Lambda, and Lambda Code) automatically enabled for new members. This indicates a comprehensive security posture across regions.

```sql+postgres
select
  region,
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
  region,
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

### List regions with EC2 scan auto-enablement status
Focus specifically on EC2 scan configuration to understand if EC2 scans are automatically enabled for new members across regions.

```sql+postgres
select
  region,
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
  region,
  ec2_auto_enable,
  case
    when ec2_auto_enable then 'EC2 scans are auto-enabled'
    else 'EC2 scans are not auto-enabled'
  end as ec2_status
from
  aws_inspector2_organization_configuration;
```

### List regions with ECR scan auto-enablement status
Focus specifically on ECR scan configuration to understand if ECR scans are automatically enabled for new members across regions.

```sql+postgres
select
  region,
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
  region,
  ecr_auto_enable,
  case
    when ecr_auto_enable then 'ECR scans are auto-enabled'
    else 'ECR scans are not auto-enabled'
  end as ecr_status
from
  aws_inspector2_organization_configuration;
```

### List regions with Lambda scan auto-enablement status
Focus specifically on Lambda scan configuration to understand if Lambda scans are automatically enabled for new members across regions.

```sql+postgres
select
  region,
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
  region,
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

### List complete regional organization configurations
Retrieve the complete organization configuration including all auto-enable settings and account limit status for comprehensive analysis across all regions.

```sql+postgres
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  title
from
  aws_inspector2_organization_configuration;
```

```sql+sqlite
select
  region,
  ec2_auto_enable,
  ecr_auto_enable,
  lambda_auto_enable,
  lambda_code_auto_enable,
  max_account_limit_reached,
  title
from
  aws_inspector2_organization_configuration;
```
