# Table: aws_backup_framework

AWS Backup Framework is a collection of controls that you can use to evaluate your backup practices.
The AWS Backup Framework will then evaluate whether your backup practices comply with your policies and highlights which
resources are not yet in compliance.

## Examples

### Basic Info

```sql
select
  framework_name,
  arn,
  framework_description,
  deployment_status,
  creation_time,
  number_of_controls
from
  aws_backup_framework;
```

### List AWS Backup Frameworks withn the last 90 days

```sql
select
  framework_name,
  arn,
  framework_description,
  deployment_status,
  creation_time,
  number_of_controls
from
  aws_backup_framework
where
  creation_time >= (current_date - interval '90' day)
order by
  creation_time;
```

### List controls that are applied to AWS Backup Frameworks

```sql
select
  framework_name,
  controls ->> 'ControlName' as control_name
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls;
```

### List AWS Backup Frameworks that are using a specific control

For this example, we are trying to get all the AWS Backup Frameworks that have the control name `BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK`.

```sql
select
  framework_name
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls
where
  controls ->> 'ControlName' = 'BACKUP_RESOURCES_PROTECTED_BY_BACKUP_VAULT_LOCK';
```

### List control name and control scope for backup AWS Backup Frameworks

This query will return an empty control scope if the control doesn't apply to a specific AWS resource type.
Otherwise if query will list the control name and the AWS resource type.

```sql
select
  framework_name,
  controls ->> 'ControlName' as control_name,
  jsonb_array_elements_text(
    case when (controls -> 'ControlScope' ->> 'ComplianceResourceTypes') is null then to_jsonb('[""]'::json)
    else (controls -> 'ControlScope' -> 'ComplianceResourceTypes')
    end
  ) as control_scope
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls
```

### List control name and control scope for a specific AWS Backup Framework

This query will return an empty control scope if the control doesn't apply to a specific AWS resource type.
Otherwise if query will list the control name and the AWS resource type.

```sql
select
  framework_name,
  controls ->> 'ControlName' as control_name,
  jsonb_array_elements_text(
    case when (controls -> 'ControlScope' ->> 'ComplianceResourceTypes') is null then to_jsonb('[""]'::json)
    else (controls -> 'ControlScope' -> 'ComplianceResourceTypes')
    end
  ) as control_scope
from
  aws_backup_framework,
  jsonb_array_elements(framework_controls) as controls
where
  framework_name = 'framework_name'
```

### Query to create the AWS config rule name from the AWS Backup Framework controls

The below query will demonstrate how to create the AWS Config rule name from the AWS Backup Framework control.
If the rule exists in AWS Config, then the end user can query the AWS Config table to see if resources in the account are compliant to the rule and what remediation action should take place.

```sql
select
  case when framework_information.control_scope = '' then concat(framework_information.control_name, '-', framework_information.framework_uuid)
  else concat(upper(framework_information.control_scope), '-', framework_information.control_name, '-', framework_information.framework_uuid)
  end as rule_name
from
(
  select
    framework_name,
    controls ->> 'ControlName' as control_name,
    jsonb_array_elements_text(
      case when (controls -> 'ControlScope' ->> 'ComplianceResourceTypes') is NULL THEN to_jsonb('[""]'::json)
      else (controls -> 'ControlScope' -> 'ComplianceResourceTypes')
      end
    ) as control_scope,
    right(arn, 36) as framework_uuid
  from
    aws_backup_framework,
    jsonb_array_elements(framework_controls) as controls
) as framework_information;
```

### Querying rules in AWS Config from AWS Backup Framework rules

This query will link the rules that were created in AWS Backup Framework to the rules defined in AWS Config.
The query will report the if a resource is compliant the defined rules in AWS Config.

```sql
select
  rule_state,
  compliance_by_config_rule,
  input_parameters
from
  aws_config_rule
inner join
(
  select
    case when framework_information.control_scope = '' then concat(framework_information.control_name, '-', framework_information.framework_uuid)
    else concat(upper(framework_information.control_scope), '-', framework_information.control_name, '-', framework_information.framework_uuid)
    end as rule_name
  from
  (
    select
      framework_name,
      controls ->> 'ControlName' as control_name,
      jsonb_array_elements_text(
        case when (controls -> 'ControlScope' ->> 'ComplianceResourceTypes') is NULL THEN to_jsonb('[""]'::json)
        else (controls -> 'ControlScope' -> 'ComplianceResourceTypes')
        end
      ) as control_scope,
      right(arn, 36) as framework_uuid
    from
      aws_backup_framework,
      jsonb_array_elements(framework_controls) as controls
  ) as framework_information
) as backup_framework
on
  aws_config_rule.name = backup_framework.rule_name;
```
