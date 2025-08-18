---
title: "Steampipe Table: aws_ssm_service_setting - Query AWS SSM Service Settings using SQL"
description: "Allows users to query AWS SSM Service Settings, providing information about account-level settings for AWS Systems Manager services."
folder: "SSM"
---

# Table: aws_ssm_service_setting - Query AWS SSM Service Settings using SQL

AWS SSM Service Settings are regional settings for AWS Systems Manager services that apply to each region. These settings define how users interact with or use a service or a feature of a service within a specific region. Service settings are managed by AWS service teams and can be customized by users with appropriate permissions.

## Table Usage Guide

The `aws_ssm_service_setting` table in Steampipe provides you with information about SSM service settings within your AWS account. This table allows you, as a DevOps engineer or system administrator, to query service setting details including setting values, status, modification history, and more. You can utilize this table to gather insights on service settings, such as default parameter tiers, automation configurations, and managed instance settings.

**Important Notes**
- In order to query service settings, the `setting_id` column must be specified in the WHERE clause. The table requires a specific setting ID to retrieve the corresponding service setting information.
- The following setting IDs are supported, for reference see: [AWS SSM Service Setting](https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_GetServiceSetting.html#API_GetServiceSetting_RequestSyntax):
  - `/ssm/managed-instance/default-ec2-instance-management-role`
  - `/ssm/automation/customer-script-log-destination`
  - `/ssm/automation/customer-script-log-group-name`
  - `/ssm/documents/console/public-sharing-permission`
  - `/ssm/managed-instance/activation-tier`
  - `/ssm/opsinsights/opscenter`
  - `/ssm/parameter-store/default-parameter-tier`
  - `/ssm/parameter-store/high-throughput-enabled`

## Examples

### Get a specific service setting across regions
Retrieve details for a specific service setting using the setting ID. This is useful for auditing individual service configurations.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier';
```

### Check parameter store settings
Review Parameter Store service settings to understand default parameter tiers and high-throughput configurations.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and region = 'us-east-1';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and region = 'us-east-1';
```

### Find customized settings
Identify service settings that have been customized from their default values. This helps track which settings have been modified by users.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and status = 'Customized';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and status = 'Customized';
```

### Check managed instance settings for a region
Review managed instance service settings to understand EC2 instance management role and activation tier configurations.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/managed-instance/default-ec2-instance-management-role'
  and region = 'us-west-2';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/managed-instance/default-ec2-instance-management-role'
  and region = 'us-west-2';
```

### Find a recently modified setting across regions
Identify service settings that have been modified recently. This is useful for tracking configuration changes.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and last_modified_date >= now() - interval '30 days';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and last_modified_date >= datetime('now', '-30 days');
```

### Get service setting with ARN for a region
Retrieve details for a specific service setting including the ARN for reference.

```sql+postgres
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  arn,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and region = 'us-east-1';
```

```sql+sqlite
select
  setting_id,
  setting_value,
  status,
  last_modified_date,
  last_modified_user,
  arn,
  region
from
  aws_ssm_service_setting
where
  setting_id = '/ssm/parameter-store/default-parameter-tier'
  and region = 'us-east-1';
```
