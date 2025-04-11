---
title: "Steampipe Table: aws_backup_framework - Query AWS Backup Frameworks using SQL"
description: "Allows users to query AWS Backup Frameworks and retrieve comprehensive data about each backup plan, including its unique ARN, version, creation and deletion dates, and more."
folder: "Backup"
---

# Table: aws_backup_framework - Query AWS Backup Frameworks using SQL

The AWS Backup service provides a centralized framework to manage and automate data backup across AWS services. It helps you to meet business and regulatory backup compliance requirements by simplifying the management and reducing the cost of backup operations. AWS Backup offers a cost-effective, fully managed, policy-based backup solution, protecting your data in AWS services.

## Table Usage Guide

The `aws_backup_framework` table in Steampipe provides you with information about each backup framework within AWS Backup service. This table empowers you, as a DevOps engineer, to query backup plan-specific details, including the backup plan's ARN, version, creation date, deletion date, and more. You can utilize this table to gather insights on backup plans, such as their status, associated rules, and other relevant metadata. The schema outlines the various attributes of the backup plan for you, including the backup plan ARN, version, creation and deletion dates, and more.

## Examples

### Basic info
This query is used to gain insights into the deployment status, creation time, and other details of your AWS backup frameworks. The practical application is to understand the configuration and status of your backup systems for effective management and troubleshooting.

```sql+postgres
select
  account_id,
  arn,
  creation_time,
  deployment_status,
  framework_controls,
  framework_description,framework_name,
  framework_status,
  number_of_controls,
  region,
  tags
from
  aws_backup_framework;
```

```sql+sqlite
select
  account_id,
  arn,
  creation_time,
  deployment_status,
  framework_controls,
  framework_description,framework_name,
  framework_status,
  number_of_controls,
  region,
  tags
from
  aws_backup_framework;
```

### List AWS frameworks created within the last 90 days
Determine the AWS frameworks that have been established within the past three months. This is beneficial for understanding recent changes and additions to your AWS environment, allowing you to stay updated on your current configurations and controls.

```sql+postgres
select
  framework_name,
  arn,
  creation_time,
  number_of_controls
from
  aws_backup_framework
where
  creation_time >= (current_date - interval '90' day)
order by
  creation_time;
```

```sql+sqlite
select
  framework_name,
  arn,
  creation_time,
  number_of_controls
from
  aws_backup_framework
where
  creation_time >= date('now','-90 day')
order by
  creation_time;
```

### List frameworks that are using a specific control (`BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK`)
Determine the frameworks which are utilizing a specific control for resource protection in a backup vault. This is useful for identifying potential areas of risk or for compliance monitoring.

```sql+postgres
select
  framework_name
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls
where
  controls ->> 'ControlName' = 'BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK';
```

```sql+sqlite
select
  framework_name
from
  aws_backup_framework
where
  json_extract(framework_controls, '$[*].ControlName') = 'BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK';
```

### List control names and scopes for each framework
Determine the areas in which specific control names and scopes are applied within each framework. This is particularly useful for understanding the scope of control within AWS backup frameworks, aiding in effective resource management and compliance.
This query will return an empty control scope if the control doesn't apply to a specific AWS resource type.
Otherwise, the query will list the control name and the AWS resource type.


```sql+postgres
select
  framework_name,
  controls ->> 'ControlName' as control_name,
  control_scope
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls,
  json_array_elements_text(coalesce(controls -> 'ControlScope' ->> 'ComplianceResourceTypes', '[""]')::json) as control_scope
where
  framework_name = 'framework_name';
```

```sql+sqlite
select
  framework_name,
  json_extract(controls.value, '$.ControlName') as control_name,
  control_scope.value as control_scope
from
  aws_backup_framework,
  json_each(framework_controls) as controls,
  json_each(json(coalesce(json_extract(controls.value, '$.ControlScope.ComplianceResourceTypes'), '[""]'))) as control_scope
where
  framework_name = 'framework_name';
```

### List framework controls that have non-compliant resources
Determine the areas in which framework controls are not compliant with the rules. This can be useful for identifying and rectifying non-compliant resources to ensure adherence to organizational policies and standards.

```sql+postgres
select
  rule_name,
  compliance_result -> 'Compliance' ->> 'ComplianceType' as compliance_type,
  compliance_result -> 'Compliance' -> 'ComplianceContributorCount' ->> 'CappedCount' as count_of_noncompliant_resources
from
  aws_config_rule
inner join
(
  -- The sub-query will create the AWS Config rule name from information stored in the AWS Backup framework table.
  select
    case when framework_information.control_scope = '' then concat(framework_information.control_name, '-', framework_information.framework_uuid)
    else concat(upper(framework_information.control_scope), '-', framework_information.control_name, '-', framework_information.framework_uuid)
    end as rule_name
  from
  (
    select
      framework_name,
      controls ->> 'ControlName' as control_name,
      control_scope,
      right(arn, 36) as framework_uuid
    from
      aws_backup_framework,
      jsonb_array_elements(framework_controls) as controls,
      json_array_elements_text(coalesce(controls -> 'ControlScope' ->> 'ComplianceResourceTypes', '[""]')::json) as control_scope
  ) as framework_information
) as backup_framework
on
  aws_config_rule.name = backup_framework.rule_name,
  jsonb_array_elements(compliance_by_config_rule) as compliance_result
where
  compliance_result -> 'Compliance' ->> 'ComplianceType' = 'NON_COMPLIANT';
```

```sql+sqlite
select
  rule_name,
  json_extract(compliance_result, '$.Compliance.ComplianceType') as compliance_type,
  json_extract(compliance_result, '$.Compliance.ComplianceContributorCount.CappedCount') as count_of_noncompliant_resources
from
  aws_config_rule
join
(
  -- The sub-query will create the AWS Config rule name from information stored in the AWS Backup framework table.
  select
    case when control_scope = '' then control_name || '-' || framework_uuid
    else upper(control_scope) || '-' || control_name || '-' || framework_uuid
    end as rule_name
  from
  (
    select
      framework_name,
      json_extract(controls, '$.ControlName') as control_name,
      control_scope,
      substr(arn, -36) as framework_uuid
    from
      aws_backup_framework,
      json_each(framework_controls) as controls,
      json_each(coalesce(json_extract(controls, '$.
```

### List framework controls that are compliant
Identify the compliant framework controls within your AWS Config rules. This allows you to gain insights into your compliance status and helps in maintaining adherence to regulatory standards.

```sql+postgres
select
  rule_name,
  compliance_result -> 'Compliance' ->> 'ComplianceType' as compliance_type
from
  aws_config_rule
inner join
(
  -- The sub-query will create the AWS Config rule name from information stored in the AWS Backup framework table.
  select
    case when framework_information.control_scope = '' then concat(framework_information.control_name, '-', framework_information.framework_uuid)
    else concat(upper(framework_information.control_scope), '-', framework_information.control_name, '-', framework_information.framework_uuid)
    end as rule_name
  from
  (
    select
      framework_name,
      controls ->> 'ControlName' as control_name,
      control_scope,
      right(arn, 36) as framework_uuid
    from
      aws_backup_framework,
      jsonb_array_elements(framework_controls) as controls,
      json_array_elements_text(coalesce(controls -> 'ControlScope' ->> 'ComplianceResourceTypes', '[""]')::json) as control_scope
  ) as framework_information
) as backup_framework
on
  aws_config_rule.name = backup_framework.rule_name,
  jsonb_array_elements(compliance_by_config_rule) as compliance_result
where
  compliance_result -> 'Compliance' ->> 'ComplianceType' = 'COMPLIANT';
```

```sql+sqlite
select
  rule_name,
  json_extract(compliance_result, '$.Compliance.ComplianceType') as compliance_type
from
  aws_config_rule
inner join
(
  -- The sub-query will create the AWS Config rule name from information stored in the AWS Backup framework table.
  select
    case when framework_information.control_scope = '' then framework_information.control_name || '-' || framework_information.framework_uuid
    else upper(framework_information.control_scope) || '-' || framework_information.control_name || '-' || framework_information.framework_uuid
    end as rule_name
  from
  (
    select
      framework_name,
      json_extract(controls, '$.ControlName') as control_name,
      control_scope,
      substr(arn, -36) as framework_uuid
    from
      aws_backup_framework,
      json_each(framework_controls) as controls,
      json_each(coalesce(json_extract(controls, '$.ControlScope.ComplianceResourceTypes'), '[""]')) as control_scope
```