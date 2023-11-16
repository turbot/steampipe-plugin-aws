---
title: "Table: aws_ssm_association - Query AWS SSM Association using SQL"
description: "Allows users to query AWS SSM Associations for detailed information about the AWS Systems Manager associations, including their status, targets, and parameters."
---

# Table: aws_ssm_association - Query AWS SSM Association using SQL

The `aws_ssm_association` table in Steampipe provides information about the AWS Systems Manager (SSM) associations. This table allows DevOps engineers to query association-specific details, including the association ID, the instance ID it is associated with, the association version, and the parameters of the association. Users can utilize this table to gather insights on associations, such as the status of associations, the targets of associations, and the parameters of associations. The schema outlines the various attributes of the SSM association, including the association name, association ID, instance ID, association version, parameters, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_association` table, you can use the `.inspect aws_ssm_association` command in Steampipe.

**Key columns**:

- `association_id`: The unique identifier for the association. This column can be used to join this table with other tables that require the association ID.
- `instance_id`: The ID of the instance that the association is associated with. This column can be used to join this table with other tables that require the instance ID.
- `name`: The name of the association. This column can be used to join this table with other tables that require the association name.

## Examples

### Basic info

```sql
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

```sql
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

### List instances targeted by the association

```sql
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

### List associations with a critical compliance severity level

```sql
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
