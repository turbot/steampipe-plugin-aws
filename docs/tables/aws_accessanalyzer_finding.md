---
title: "Steampipe Table: aws_accessanalyzer_finding - Query AWS Access Analyzer Findings using SQL"
description: "Allows users to query Access Analyzer findings in AWS IAM to retrieve detailed information about potential security risks."
folder: "Access Analyzer"
---

# Table: aws_accessanalyzer_finding - Query AWS Access Analyzer Findings using SQL

AWS Access Analyzer findings provide detailed information about potential security risks in your AWS environment. These findings are generated when Access Analyzer identifies resources that are shared with an external entity, highlighting potential unintended access. By analyzing the resource-based policies, Access Analyzer helps you understand how access to your resources is granted and suggests modifications to achieve desired access policies, enhancing your security posture.

## Table Usage Guide

The `aws_accessanalyzer_finding` table in Steampipe allows you to query information related to findings from the AWS IAM Access Analyzer. This table is essential for security and compliance teams, enabling them to identify, analyze, and manage findings related to resource access policies. Through this table, users can access detailed information about each finding, including the actions involved, the condition that led to the finding, the resource and principal involved, and the finding's status. By leveraging this table, you can efficiently address security and compliance issues in your AWS environment.

## Examples

### Basic Info

Retrieve essential details of findings to understand potential access issues and their current status. This query helps in identifying the nature of each finding, the resources involved, and the actions recommended or taken to resolve these issues.

```sql+postgres
select
  id,
  access_analyzer_arn,
  analyzed_at,
  resource_type,
  status,
  is_public
from
  aws_accessanalyzer_finding;
```

```sql+sqlite
select
  id,
  analyzed_at,
  resource_type,
  status,
  is_public
from
  aws_accessanalyzer_finding;
```

### Findings involving public access

Identify findings where resources are potentially exposed to public access. Highlighting such findings is critical for prioritizing issues that may lead to unauthorized access. This query helps in swiftly identifying and addressing potential vulnerabilities, ensuring that resources are adequately secured against public exposure.

```sql+postgres
select
  id,
  resource_type,
  access_analyzer_arn,
  status,
  is_public
from
  aws_accessanalyzer_finding
where
  is_public = true;
```

```sql+sqlite
select
  id,
  resource_type,
  access_analyzer_arn,
  status,
  is_public
from
  aws_accessanalyzer_finding
where
  is_public = true;
```

### Findings by resource type

Aggregate findings by resource type to focus remediation efforts on specific types of resources. This categorization helps in streamlining the security review process by allowing teams to prioritize resources based on their sensitivity and exposure.

```sql+postgres
select
  resource_type,
  count(*) as findings_count
from
  aws_accessanalyzer_finding
group by
  resource_type;
```
  
```sql+sqlite
select
  resource_type,
  count(*) as findings_count
from
  aws_accessanalyzer_finding
group by
  resource_type;
```

### Recent findings

Focus on findings that have been identified recently to address potentially new security risks. This query aids in maintaining an up-to-date security posture by ensuring that recent findings are promptly reviewed and addressed.

```sql+postgres
select
  id,
  resource,
  status,
  analyzed_at
from
  aws_accessanalyzer_finding
where
  analyzed_at > current_date - interval '30 days';
```

```sql+sqlite
select
  id,
  resource,
  status,
  analyzed_at
from
  aws_accessanalyzer_finding
where
  analyzed_at > date('now', '-30 day');
```
