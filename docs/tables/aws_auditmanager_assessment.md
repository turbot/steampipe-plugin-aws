---
title: "Table: aws_auditmanager_assessment - Query AWS Audit Manager Assessments using SQL"
description: "Allows users to query AWS Audit Manager Assessments to retrieve detailed information about each assessment."
---

# Table: aws_auditmanager_assessment - Query AWS Audit Manager Assessments using SQL

The `aws_auditmanager_assessment` table in Steampipe provides information about assessments within AWS Audit Manager. This table allows DevOps engineers to query assessment-specific details, including the assessment status, scope, roles, and associated metadata. Users can utilize this table to gather insights on assessments, such as assessment status, scope of the assessments, roles associated with the assessments, and more. The schema outlines the various attributes of the AWS Audit Manager assessment, including the assessment ID, name, description, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_auditmanager_assessment` table, you can use the `.inspect aws_auditmanager_assessment` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the assessment. This can be used to join this table with other tables that contain AWS resource ARNs.
- `id`: The unique identifier for the assessment. This can be used to join this table with other tables that contain AWS Audit Manager assessment IDs.
- `status`: The status of the assessment. This can be useful for filtering assessments based on their current status.

## Examples

### Basic info

```sql
select
  name,
  arn,
  status,
  compliance_type
from
  aws_auditmanager_assessment;
```


### List assessments with public audit bucket

```sql
select
  a.name,
  a.arn,
  a.assessment_report_destination,
  a.assessment_report_destination_type,
  b.bucket_policy_is_public as is_public_bucket
from
  aws_auditmanager_assessment as a
join aws_s3_bucket as b on a.assessment_report_destination = 's3://' || b.Name and b.bucket_policy_is_public;
```


### List inactive assessments

```sql
select
  name,
  arn,
  status
from
  aws_auditmanager_assessment
where
  status <> 'ACTIVE';
```
