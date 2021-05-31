# Table: aws_auditmanager_evidence_folder

The evidence folder that contains the evidence. User activity evidence is collected from AWS CloudTrail logs. Configuration data evidence is collected from snapshots of other AWS services such as Amazon EC2, Amazon S3, or IAM.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  assessment_id,
  control_set_id,
  control_id,
  total_evidence
from
  aws_auditmanager_evidence_folder;
```

### Get evidence folder count by assessment id

```sql
select
  assessment_id,
  count(id) as evidence_folder_count
from
  aws_auditmanager_evidence_folder
group by
  assessment_id;
```
