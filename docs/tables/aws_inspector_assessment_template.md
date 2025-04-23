---
title: "Steampipe Table: aws_inspector_assessment_template - Query AWS Inspector Assessment Templates using SQL"
description: "Allows users to query AWS Inspector Assessment Templates to gain insights into each template's configuration, including ARN, duration, rules package ARNs, and user attributes for findings."
folder: "Inspector"
---

# Table: aws_inspector_assessment_template - Query AWS Inspector Assessment Templates using SQL

The AWS Inspector Assessment Template is a resource within AWS Inspector that helps you analyze the behavior of the applications you run on AWS and helps identify potential security issues. It automatically assesses applications for exposure, vulnerabilities, and deviations from best practices. After performing an assessment, AWS Inspector produces a detailed list of security findings prioritized by level of severity.

## Table Usage Guide

The `aws_inspector_assessment_template` table in Steampipe provides you with information about assessment templates within AWS Inspector. This table allows you, as a DevOps engineer, security analyst, or other technical professional, to query template-specific details, including the ARN, duration, rules package ARNs, and user attributes for findings. You can utilize this table to gather insights on assessment templates, such as identifying templates with specific rules, verifying template configurations, and more. The schema outlines the various attributes of the assessment template for you, including the template ARN, duration, rules package ARNs, user attributes for findings, and associated tags.

## Examples

### Basic info
Explore which AWS Inspector assessment templates are in use to understand their distribution across regions and assess how frequently they are run. This can help identify potential areas for optimizing resource usage and improving security assessment practices.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  region
from
  aws_inspector_assessment_template;
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  region
from
  aws_inspector_assessment_template;
```


### List assessment templates that have no assigned finding attributes
Determine the areas in which assessment templates in AWS Inspector have not been assigned any finding attributes. This is useful for identifying potential gaps in your security assessment configuration.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  user_attributes_for_findings = '[]';
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  user_attributes_for_findings = '[]';
```

### List assessment templates that have no assessment runs
Identify instances where certain assessment templates in your AWS Inspector setup have not been used for any assessment runs. This can help pinpoint unused resources and optimize your security assessment process.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```

### List assessment templates with run duration less than 1 hour
Determine the areas in which assessment templates have a run duration of less than an hour, helpful for identifying any quick assessments in your AWS Inspector setup.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```

### List assessment templates that have no assessment runs
Identify assessment templates that are yet to be used for any assessment runs. This could be useful to clean up unused resources or to pinpoint areas where assessments are not being conducted.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```


### List assessment templates with run duration less than 1 hour
Determine the areas in which AWS Inspector Assessment templates have a run duration of less than an hour. This can be useful for identifying templates that may be completing their run too quickly, potentially missing out on thorough inspections.

```sql+postgres
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```

```sql+sqlite
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```