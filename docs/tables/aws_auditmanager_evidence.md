---
title: "Steampipe Table: aws_auditmanager_evidence - Query AWS Audit Manager Evidence using SQL"
description: "Allows users to query AWS Audit Manager Evidence, providing detailed information about evidence resources associated with assessments in AWS Audit Manager."
folder: "Audit Manager"
---

# Table: aws_auditmanager_evidence - Query AWS Audit Manager Evidence using SQL

The AWS Audit Manager Evidence is a component of AWS Audit Manager service that automates the collection and organization of evidence for audits. It simplifies the process of gathering necessary documents to demonstrate to auditors that your controls are operating effectively. This resource assists in continuously auditing your AWS usage to simplify risk assessment and compliance with regulations and industry standards.

## Table Usage Guide

The `aws_auditmanager_evidence` table in Steampipe provides you with information about evidence resources within AWS Audit Manager. This table allows you, as a DevOps engineer, to query evidence-specific details, including the source, collection method, and associated metadata. You can utilize this table to gather insights on evidence, such as the evidence state, evidence by type, and the AWS resource from which the evidence was collected. The schema outlines the various attributes of the evidence for you, including the evidence id, assessment id, control set id, evidence folder id, and associated tags.

## Examples

### Basic info
Explore the various pieces of evidence collected in AWS Audit Manager to understand their association with different control sets and IAM identities. This can help in assessing the compliance status of your AWS resources and identifying areas that may need attention.

```sql+postgres
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

```sql+sqlite
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
Analyze the distribution of evidence across different folders in AWS Audit Manager to understand the workload and prioritize accordingly. This can help in efficiently managing and reviewing the collected evidence.

```sql+postgres
select
  evidence_folder_id,
  count(id) as evidence_count
from
  aws_auditmanager_evidence
group by
  evidence_folder_id;
```

```sql+sqlite
select
  evidence_folder_id,
  count(id) as evidence_count
from
  aws_auditmanager_evidence
group by
  evidence_folder_id;
```