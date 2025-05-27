---
title: "Steampipe Table: aws_auditmanager_evidence_folder - Query AWS Audit Manager Evidence Folders using SQL"
description: "Allows users to query AWS Audit Manager Evidence Folders to get comprehensive details about the evidence folders in the AWS Audit Manager service."
folder: "Audit Manager"
---

# Table: aws_auditmanager_evidence_folder - Query AWS Audit Manager Evidence Folders using SQL

The AWS Audit Manager Evidence Folders are used to organize and store evidence collected for assessments. This evidence can be automatically collected by AWS Audit Manager or manually uploaded by users. The evidence folders help in managing compliance audits and providing detailed proof of how the data is being handled within the AWS environment.

## Table Usage Guide

The `aws_auditmanager_evidence_folder` table in Steampipe provides you with information about evidence folders within AWS Audit Manager. This table allows you, as a DevOps engineer, to query evidence folder-specific details, including the ID, ARN, name, date created, and associated metadata. You can utilize this table to gather insights on evidence folders, such as the total count of evidence in the folder, the status of the evidence, verification of evidence source, and more. The schema outlines the various attributes of the evidence folder for you, including the evidence folder ID, ARN, creation date, and associated tags.

## Examples

### Basic info
Explore which evidence folders exist within your AWS Audit Manager to better manage and assess your compliance controls and evidence. This can help you identify areas where you might need to gather additional evidence or focus your auditing efforts.

```sql+postgres
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

```sql+sqlite
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
Explore how many evidence folders are associated with each assessment in your AWS Audit Manager. This is useful for understanding the volume of evidence collected for each audit, aiding in audit management and review processes.

```sql+postgres
select
  assessment_id,
  count(id) as evidence_folder_count
from
  aws_auditmanager_evidence_folder
group by
  assessment_id;
```

```sql+sqlite
select
  assessment_id,
  count(id) as evidence_folder_count
from
  aws_auditmanager_evidence_folder
group by
  assessment_id;
```