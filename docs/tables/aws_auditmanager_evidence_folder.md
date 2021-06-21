# Table: aws_auditmanager_evidence_folder

Evidence folders contain user activity evidence that is collected from CloudTrail logs. Configuration data evidence is collected from snapshots of other AWS services such as Amazon EC2, Amazon S3, or IAM.

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

### Count the number of evidence folders by assessment ID

```sql
select
  assessment_id,
  count(id) as evidence_folder_count
from
  aws_auditmanager_evidence_folder
group by
  assessment_id;
```
