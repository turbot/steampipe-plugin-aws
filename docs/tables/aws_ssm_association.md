---
title: "Steampipe Table: aws_ssm_association - Query AWS SSM Association using SQL"
description: "Allows users to query AWS SSM Associations for detailed information about the AWS Systems Manager associations, including their status, targets, and parameters."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_association - Query AWS SSM Association using SQL

The AWS SSM Association is a component of AWS Systems Manager that allows you to configure and manage instances at scale. It enables you to perform administrative tasks such as installing patches, updating agents, or applying policies. With SSM Association, you can automate the process of keeping your managed instances in a desired state.

## Table Usage Guide

The `aws_ssm_association` table in Steampipe provides you with information about the AWS Systems Manager (SSM) associations. This table enables you, as a DevOps engineer, to query association-specific details, including the association ID, the instance ID it is associated with, the association version, and the parameters of the association. You can utilize this table to gather insights on associations, such as the status of associations, the targets of associations, and the parameters of associations. The schema outlines the various attributes of the SSM association for you, including the association name, association ID, instance ID, association version, parameters, and more.

## Examples

### Basic info
Explore which AWS Systems Manager (SSM) associations are currently active, including their compliance severity and last execution date. This can help in managing and monitoring the AWS SSM associations effectively.

```sql+postgres
select
  association_id,
  association_name,
  arn,
  association_version,
  last_execution_date,
  document_name,
  compliance_severity,
  region
from
  aws_ssm_association;
```

```sql+sqlite
select
  association_id,
  association_name,
  arn,
  association_version,
  last_execution_date,
  document_name,
  compliance_severity,
  region
from
  aws_ssm_association;
```

### List associations that have a failed status
Identify instances where AWS System Manager associations have failed. This is useful for troubleshooting and rectifying the issues causing the failure.

```sql+postgres
select
  association_id,
  overview ->> 'AssociationStatusAggregatedCount' as association_status_aggregated_count,
  overview ->> 'DetailedStatus' as detailed_status,
  overview ->> 'Status' as status
from
  aws_ssm_association
where
  overview ->> 'Status' = 'Failed';
```

```sql+sqlite
select
  association_id,
  json_extract(overview, '$.AssociationStatusAggregatedCount') as association_status_aggregated_count,
  json_extract(overview, '$.DetailedStatus') as detailed_status,
  json_extract(overview, '$.Status') as status
from
  aws_ssm_association
where
  json_extract(overview, '$.Status') = 'Failed';
```

### List instances targeted by the association
Discover the instances that are targeted by a specific association within the AWS SSM service. This is beneficial for gaining insights into how your resources are being utilized and managed.

```sql+postgres
select
  association.association_id as association_id,
  target ->> 'Key' as target_key,
  target ->> 'Values' as target_value,
  instances
from
  aws_ssm_association as association,
  jsonb_array_elements(targets) as target,
  jsonb_array_elements_text(target -> 'Values') as instances
where
  target ->> 'Key' = 'InstanceIds';
```

```sql+sqlite
select
  association.association_id as association_id,
  json_extract(target.value, '$.Key') as target_key,
  json_extract(target.value, '$.Values') as target_value,
  instances.value as instances
from
  aws_ssm_association as association,
  json_each(association.targets) as target,
  json_each(json_extract(target.value, '$.Values')) as instances
where
  json_extract(target.value, '$.Key') = 'InstanceIds';
```

### List associations with a critical compliance severity level
Discover the associations that have a critical compliance severity level. This can be useful for identifying potential risk areas in your AWS Simple Systems Manager configuration.

```sql+postgres
select
  association_id,
  association_name,
  targets,
  document_name
from
  aws_ssm_association
where
  compliance_severity = 'CRITICAL';
```

```sql+sqlite
select
  association_id,
  association_name,
  targets,
  document_name
from
  aws_ssm_association
where
  compliance_severity = 'CRITICAL';
```