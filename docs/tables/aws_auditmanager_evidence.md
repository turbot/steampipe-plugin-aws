---
title: "Table: aws_auditmanager_evidence - Query AWS Audit Manager Evidence using SQL"
description: "Allows users to query AWS Audit Manager Evidence, providing detailed information about evidence resources associated with assessments in AWS Audit Manager."
---

# Table: aws_auditmanager_evidence - Query AWS Audit Manager Evidence using SQL

The `aws_auditmanager_evidence` table in Steampipe provides information about evidence resources within AWS Audit Manager. This table allows DevOps engineers to query evidence-specific details, including the source, collection method, and associated metadata. Users can utilize this table to gather insights on evidence, such as the evidence state, evidence by type, and the AWS resource from which the evidence was collected. The schema outlines the various attributes of the evidence, including the evidence id, assessment id, control set id, evidence folder id, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_auditmanager_evidence` table, you can use the `.inspect aws_auditmanager_evidence` command in Steampipe.

**Key columns**:

- `evidence_id`: The unique identifier for the evidence. This can be used to join this table with other tables to get more detailed information about the evidence.
- `assessment_id`: The identifier for the assessment that the evidence is associated with. This can be used to join with the assessment table to get more context about the evidence.
- `control_set_id`: The identifier for the control set that the evidence is associated with. This can be used to join with the control set table to get more information about the controls related to the evidence.

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
