---
title: "Table: aws_auditmanager_evidence_folder - Query AWS Audit Manager Evidence Folders using SQL"
description: "Allows users to query AWS Audit Manager Evidence Folders to get comprehensive details about the evidence folders in the AWS Audit Manager service."
---

# Table: aws_auditmanager_evidence_folder - Query AWS Audit Manager Evidence Folders using SQL

The `aws_auditmanager_evidence_folder` table in Steampipe provides information about evidence folders within AWS Audit Manager. This table allows DevOps engineers to query evidence folder-specific details, including the ID, ARN, name, date created, and associated metadata. Users can utilize this table to gather insights on evidence folders, such as the total count of evidence in the folder, the status of the evidence, verification of evidence source, and more. The schema outlines the various attributes of the evidence folder, including the evidence folder ID, ARN, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_auditmanager_evidence_folder` table, you can use the `.inspect aws_auditmanager_evidence_folder` command in Steampipe.

**Key columns**:

- `id`: The unique identifier for the evidence folder. This column can be used to join this table with other tables that contain evidence folder IDs.
- `arn`: The Amazon Resource Name (ARN) specifying the evidence folder. This column is useful for joining with other tables that use ARNs.
- `assessment_id`: The identifier for the assessment that the evidence folder belongs to. This column can be used to join with the `aws_auditmanager_assessment` table to pull more detailed information about the assessment.

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
