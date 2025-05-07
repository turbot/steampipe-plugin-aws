---
title: "Steampipe Table: aws_macie2_finding - Query AWS Macie2 Findings using SQL"
description: "Allows users to query AWS Macie2 Findings to retrieve detailed information about security and data privacy findings, including their type, severity, status, and associated metadata."
folder: "Macie2"
---

# Table: aws_macie2_finding - Query AWS Macie2 Findings using SQL

AWS Macie2 is a data security and data privacy service that uses machine learning and pattern matching to discover and protect your sensitive data in AWS. It automatically discovers, classifies, and protects sensitive data stored in Amazon S3. The `aws_macie2_finding` table provides information about security and data privacy findings that Macie2 has generated, including details about sensitive data discovery, policy violations, and other security-related events.

## Table Usage Guide

The `aws_macie2_finding` table in Steampipe provides you with information about security and data privacy findings within AWS Macie2. This table allows you, as a security analyst, compliance officer, or DevOps engineer, to query finding-specific details, including the finding type, severity, status, and associated metadata. You can utilize this table to gather insights on findings, such as identifying sensitive data exposure, monitoring policy violations, and tracking remediation efforts. The schema outlines the various attributes of the Macie2 finding for you, including the finding ID, ARN, type, severity, status, and associated tags.

## Examples

### List all findings

Retrieve basic information about your AWS Macie2 findings, including their type, severity, and status. This can be useful for getting an overview of security and data privacy issues in your AWS account.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  created_at
from
  aws_macie2_finding;
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  created_at
from
  aws_macie2_finding;
```

### Get details of a specific finding

Query detailed information about a specific finding to understand its impact and required actions.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  description,
  category,
  count,
  resources_affected
from
  aws_macie2_finding
where
  id = '12345678901234567890123456789012';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  description,
  category,
  count,
  resources_affected
from
  aws_macie2_finding
where
  id = '12345678901234567890123456789012';
```

### List high severity findings

Identify findings that require immediate attention based on their severity level.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  severity = 'HIGH';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  severity = 'HIGH';
```

### List findings by type

Filter findings by their type to focus on specific categories of security or data privacy issues.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  type = 'SensitiveData:S3Object/Multiple';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  type = 'SensitiveData:S3Object/Multiple';
```

### List active findings

Monitor currently active findings that require attention or remediation.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  created_at,
  description
from
  aws_macie2_finding
where
  status = 'ACTIVE';
```

### List findings with remediation information

Review findings that include specific remediation steps or recommendations.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  remediation
from
  aws_macie2_finding
where
  remediation is not null;
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  remediation
from
  aws_macie2_finding
where
  remediation is not null;
```

### List findings with classification details

Examine findings that include detailed classification information about the discovered data.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  classification_details
from
  aws_macie2_finding
where
  classification_details is not null;
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  classification_details
from
  aws_macie2_finding
where
  classification_details is not null;
```

### List findings by date range

Monitor findings generated within a specific time period.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  created_at
from
  aws_macie2_finding
where
  created_at >= now() - interval '7 days';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  created_at
from
  aws_macie2_finding
where
  created_at >= datetime('now', '-7 days');
```

### List findings with affected resources

Identify findings that include information about affected AWS resources.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  resources_affected
from
  aws_macie2_finding
where
  resources_affected is not null;
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  resources_affected
from
  aws_macie2_finding
where
  resources_affected is not null;
```

### List findings with sample data

Review findings that include sample data to better understand the nature of the discovered sensitive information.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  sample
from
  aws_macie2_finding
where
  sample is not null;
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  sample
from
  aws_macie2_finding
where
  sample is not null;
```

### List findings by account and region

Filter findings by specific AWS account and region to focus on particular environments.

```sql+postgres
select
  id,
  type,
  severity,
  status,
  source_account_id,
  source_region
from
  aws_macie2_finding
where
  source_account_id = '123456789012'
  and source_region = 'us-east-1';
```

```sql+sqlite
select
  id,
  type,
  severity,
  status,
  source_account_id,
  source_region
from
  aws_macie2_finding
where
  source_account_id = '123456789012'
  and source_region = 'us-east-1';
```
