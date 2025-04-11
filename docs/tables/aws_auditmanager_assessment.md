---
title: "Steampipe Table: aws_auditmanager_assessment - Query AWS Audit Manager Assessments using SQL"
description: "Allows users to query AWS Audit Manager Assessments to retrieve detailed information about each assessment."
folder: "Audit Manager"
---

# Table: aws_auditmanager_assessment - Query AWS Audit Manager Assessments using SQL

The AWS Audit Manager Assessment is a feature of AWS Audit Manager that helps you continuously audit your AWS usage to simplify your risk management and compliance. It automates evidence collection to enable you to scale your audit capability as your AWS usage grows. This tool facilitates assessment of the effectiveness of your controls and helps you maintain continuous compliance by managing audits throughout their lifecycle.

## Table Usage Guide

The `aws_auditmanager_assessment` table in Steampipe provides you with information about assessments within AWS Audit Manager. This table allows you, as a DevOps engineer, to query assessment-specific details, including the assessment status, scope, roles, and associated metadata. You can utilize this table to gather insights on assessments, such as assessment status, scope of the assessments, roles associated with the assessments, and more. The schema outlines the various attributes of the AWS Audit Manager assessment for you, including the assessment ID, name, description, status, and associated tags.

## Examples

### Basic info
Explore which AWS Audit Manager assessments are currently active and what their compliance types are. This can be useful for keeping track of your organization's compliance status and ensuring all assessments are functioning as expected.

```sql+postgres
select
  name,
  arn,
  status,
  compliance_type
from
  aws_auditmanager_assessment;
```

```sql+sqlite
select
  name,
  arn,
  status,
  compliance_type
from
  aws_auditmanager_assessment;
```


### List assessments with public audit bucket
This query is useful for identifying assessments that are associated with a public audit bucket. This can help in enhancing the security measures by pinpointing potential areas of vulnerability, as public audit buckets can be accessed by anyone.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which assessments are not currently active, enabling you to focus resources on those that require attention or action.

```sql+postgres
select
  name,
  arn,
  status
from
  aws_auditmanager_assessment
where
  status <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  arn,
  status
from
  aws_auditmanager_assessment
where
  status != 'ACTIVE';
```