# Table: aws_ssm_association

An association is a configuration that is assigned to your managed instances. The configuration defines the state that you want to maintain on your instances. For example, an association can specify that antivirus software must be installed and running on your instances, or that certain ports must be closed. The association specifies a schedule for when the configuration is applied once or reapplied at specified times.

## Examples

### Basic info

```sql
select
  association_id,
  association_name,
  association_version,
  last_execution_date,
  document_name,
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
