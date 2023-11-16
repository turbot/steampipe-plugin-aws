---
title: "Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL"
description: "Allows users to query Access Analyzer Analyzer in AWS IAM to retrieve information about analyzers."
---

# Table: aws_accessanalyzer_analyzer - Query AWS Access Analyzer using SQL

The `aws_accessanalyzer_analyzer` table in Steampipe provides information about analyzers within AWS IAM Access Analyzer. This table allows DevOps engineers to query analyzer-specific details, including the analyzer ARN, type, status, and associated metadata. Users can utilize this table to gather insights on analyzers, such as the status of each analyzer, the type of analyzer, and the resource that was analyzed. The schema outlines the various attributes of the Access Analyzer, including the analyzer ARN, creation time, last resource scanned, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_accessanalyzer_analyzer` table, you can use the `.inspect aws_accessanalyzer_analyzer` command in Steampipe.

**Key columns**:

- `arn`: The ARN (Amazon Resource Name) of the analyzer. This can be used to join with other tables where the analyzer ARN is required.
- `name`: The name of the analyzer. This can be used to join with other tables where the analyzer name is required.
- `type`: The type of analyzer. This can be useful for filtering the data based on the analyzer type.

## Examples

### Basic info
The query provides an overview of AWS Access Analyzer analyzers in a user's environment. It helps in monitoring the current status and types of analyzers, along with the details of the most recent resources analyzed. This is useful for administrators and security personnel to ensure that their AWS environment is continuously scanned for compliance and security risks, and to stay informed about the analyzer's activities and findings.

```sql
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
The query identifies and provides details on all active AWS Access Analyzer analyzers. It is particularly useful for ensuring that the necessary analyzers are operational and actively scanning resources. This information aids in maintaining continuous compliance and security oversight by highlighting only those analyzers currently in an active state, along with their last analyzed resources and associated tags. This enables efficient tracking and management of security analysis tools within the AWS environment.

```sql
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

### List analyzers with findings that need to be resolved
The query focuses on identifying active AWS Access Analyzer analyzers that have unresolved findings. It serves as a tool for security and compliance teams to pinpoint which analyzers have detected potential issues, needing immediate attention. By filtering for active analyzers with existing findings, it streamlines the process of addressing security or compliance concerns within the AWS environment, ensuring that no critical issues are overlooked. This aids in maintaining a secure and compliant cloud infrastructure.

```sql
select
  name,
  status,
  type,
  last_resource_analyzed
from
  aws_accessanalyzer_analyzer
where
  status = 'ACTIVE'
  and findings is not null;
```
