# Table: aws_auditmanager_evidence

Each active assessment in AWS Audit Manager automatically collects evidence from a range of data sources. Every assessment has a defined scope that specifies the AWS services and accounts from which Audit Manager collects data. Each of these defined data sources contains multiple resources, with each resource being a system asset inventory that you own. Essentially, evidence collection in Audit Manager involves the assessment of each in-scope resource.

## Examples

### Basic info

```sql
select
  id,
  arn,
  evidence_folder_id,
  evidence_by_type,
  iam_id,
  control_set_id
from
  aws_auditmanager_evidence;
```

### Get evidence count by evidence folder

```sql
select
  evidence_folder_id,
  count(id) as evidence_count
from
  aws_auditmanager_evidence
group by
  evidence_folder_id;
```
