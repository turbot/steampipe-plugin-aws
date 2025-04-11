---
title: "Steampipe Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL"
description: "Allows users to query Access Analyzer Analyzer in AWS IAM to retrieve information about analyzers."
folder: "Access Analyzer"
---

# Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL

The AWS Access Analyzer is a service that helps to identify resources in your organization and accounts, such as S3 buckets or IAM roles, that are shared with an external entity. It uses logic-based reasoning to analyze the resource-based policies in your AWS environment, allowing you to identify unintended access to your resources and data. This helps in mitigating potential security risks.

## Table Usage Guide

The `aws_accessanalyzer_analyzer` table in Steampipe provides you with information about analyzers within AWS IAM Access Analyzer. This table allows you, as a DevOps engineer, to query analyzer-specific details, including the analyzer ARN, type, status, and associated metadata. You can utilize this table to gather insights on analyzers, such as the status of each analyzer, the type of analyzer, and the resource that was analyzed. The schema outlines the various attributes of the Access Analyzer for you, including the analyzer ARN, creation time, last resource scanned, and associated tags.

## Examples

### Basic info
Explore the status and type of your AWS Access Analyzer to understand when the last resource was analyzed. This could be beneficial for maintaining security and compliance in your AWS environment.The query provides an overview of AWS Access Analyzer analyzers in a user's environment. It helps in monitoring the current status and types of analyzers, along with the details of the most recent resources analyzed. This is useful for administrators and security personnel to ensure that their AWS environment is continuously scanned for compliance and security risks, and to stay informed about the analyzer's activities and findings.

```sql+postgres
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  status,
  type
from
  aws_accessanalyzer_analyzer;
```

```sql+sqlite
select
  name,
  last_resource_analyzed,
  last_resource_analyzed_at,
  status,
  type
from
  aws_accessanalyzer_analyzer;
```

### List analyzers which are enabled
Determine the areas in which AWS Access Analyzer is active to gain insights into potential security and access control issues. This is useful for maintaining optimal security practices and ensuring that all analyzers are functioning as expected.The query identifies and provides details on all active AWS Access Analyzer analyzers. It is particularly useful for ensuring that the necessary analyzers are operational and actively scanning resources. This information aids in maintaining continuous compliance and security oversight by highlighting only those analyzers currently in an active state, along with their last analyzed resources and associated tags. This enables efficient tracking and management of security analysis tools within the AWS environment.

```sql+postgres
select
  name,
  status
  last_resource_analyzed,
  last_resource_analyzed_at,
  tags
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  name,
  status,
  last_resource_analyzed,
  last_resource_analyzed_at,
  tags
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE';
```

### List analyzers with findings that need to be resolved
Explore which active AWS Access Analyzer instances have findings that require resolution. This is useful in identifying potential security risks that need immediate attention.The query focuses on identifying active AWS Access Analyzer analyzers that have unresolved findings. It serves as a tool for security and compliance teams to pinpoint which analyzers have detected potential issues, needing immediate attention. By filtering for active analyzers with existing findings, it streamlines the process of addressing security or compliance concerns within the AWS environment, ensuring that no critical issues are overlooked. This aids in maintaining a secure and compliant cloud infrastructure.

```sql+postgres
select
  a.arn as analyzer_arn,
  a.name as analyzer_name,
  a.region as analyzer_region,
  a.account_id,
  count(f.id) as findings_count
from
  aws_accessanalyzer_analyzer as a
  join aws_accessanalyzer_finding as f on f.access_analyzer_arn = a.arn
where
  a.status = 'ACTIVE'
group by
  a.arn,
  a.name,
  a.region,
  a.account_id
having
  count(f.id) > 0;
```

```sql+sqlite
select
  a.arn as analyzer_arn,
  a.name as analyzer_name,
  a.region as analyzer_region,
  a.account_id,
  count(f.id) as findings_count
from
  aws_accessanalyzer_analyzer as a
  join aws_accessanalyzer_finding as f on f.access_analyzer_arn = a.arn
where
  a.status = 'ACTIVE'
group by
  a.arn,
  a.name,
  a.region,
  a.account_id
having
  count(f.id) > 0;
```