---
title: "Table: aws_backup_framework - Query AWS Backup Frameworks using SQL"
description: "Allows users to query AWS Backup Frameworks and retrieve comprehensive data about each backup plan, including its unique ARN, version, creation and deletion dates, and more."
---

# Table: aws_backup_framework - Query AWS Backup Frameworks using SQL

The `aws_backup_framework` table in Steampipe provides information about each backup framework within AWS Backup service. This table allows DevOps engineers to query backup plan-specific details, including the backup plan's ARN, version, creation date, deletion date, and more. Users can utilize this table to gather insights on backup plans, such as their status, associated rules, and other relevant metadata. The schema outlines the various attributes of the backup plan, including the backup plan ARN, version, creation and deletion dates, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_framework` table, you can use the `.inspect aws_backup_framework` command in Steampipe.

### Key columns:

- `backup_plan_arn`: The Amazon Resource Name (ARN) for the backup plan. This column is important as it uniquely identifies each backup plan and can be used to join this table with other tables that contain backup plan ARN data.
- `creation_date`: The date and time a backup plan is created. This column is useful for tracking the lifespan of backup plans and identifying those that might need updating or retiring.
- `version_id`: The unique version ID of a backup plan. This column is important for tracking changes and updates to backup plans over time.

## Examples

### Basic info

```sql
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

```sql
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

### List frameworks that are using a specific control (`BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK`)

```sql
select
  framework_name
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls
where
  controls ->> 'ControlName' = 'BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK';
```

### List control names and scopes for each framework

This query will return an empty control scope if the control doesn't apply to a specific AWS resource type.
Otherwise, the query will list the control name and the AWS resource type.

```sql
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

### List framework controls that have non-compliant resources

```sql
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

### List framework controls that are compliant

```sql
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
