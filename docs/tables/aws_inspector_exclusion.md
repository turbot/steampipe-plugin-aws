---
title: "Steampipe Table: aws_inspector_exclusion - Query AWS Inspector Exclusions using SQL"
description: "Allows users to query AWS Inspector Exclusions and retrieve data about AWS Inspector Exclusions, including their ARNs, descriptions, and recommendations."
folder: "Inspector"
---

# Table: aws_inspector_exclusion - Query AWS Inspector Exclusions using SQL

The AWS Inspector Exclusion is a feature of AWS Inspector, a service that helps you improve the security and compliance of applications deployed on AWS. Exclusions are defined during the assessment template creation process, and they represent a scope that is excluded from the assessment run. They are used to exclude false positives from the assessment findings, allowing you to focus on truly significant security findings.

## Table Usage Guide

The `aws_inspector_exclusion` table in Steampipe provides you with information about exclusions within AWS Inspector. This table allows you, as a DevOps engineer, to query exclusion-specific details, including the ARN, description, and recommendation. You can utilize this table to gather insights on exclusions, such as their status, the reasons behind their exclusions, and more. The schema outlines the various attributes of the AWS Inspector exclusion for you, including the ARN, description, recommendation, and associated metadata.

## Examples

### Basic info
Determine the areas in which AWS Inspector exclusions are applied to gain insights into your AWS security setup. This can help in understanding the scope and impact of these exclusions within your infrastructure.

```sql+postgres
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion;
```

```sql+sqlite
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion;
```

### List exclusions associated with an assessment run
Identify the exclusions linked to a specific assessment run to understand the areas that were omitted during the assessment. This can be helpful in reviewing the comprehensiveness of the assessment or identifying potential blind spots.

```sql+postgres
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion
where
  assessment_run_arn = 'arn:aws:inspector:us-east-1:012345678912:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz';
```

```sql+sqlite
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion
where
  assessment_run_arn = 'arn:aws:inspector:us-east-1:012345678912:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz';
```

### Get the attribute and scope details for each exclusion
Explore the specifics of each exclusion in your AWS Inspector to understand the nature and extent of what is excluded. This can be useful in auditing your security setup, ensuring that no critical resources are accidentally excluded from inspections.

```sql+postgres
select
  arn,
  jsonb_pretty(attributes) as attributes,
  jsonb_pretty(scopes) as scopes
from
  aws_inspector_exclusion;
```

```sql+sqlite
select
  arn,
  attributes,
  scopes
from
  aws_inspector_exclusion;
```

### Count the number of exclusions whose type is 'Agent not found'
Determine the areas in which the number of 'Agent not found' exclusions are highest. This helps in identifying regions that might have connectivity issues or where agents are not deployed properly.

```sql+postgres
select
  arn,
  region,
  title,
  count(arn)
from
  aws_inspector_exclusion
group by
  arn,
  region,
  title
order by
  count desc;
```

```sql+sqlite
select
  arn,
  region,
  title,
  count(arn)
from
  aws_inspector_exclusion
group by
  arn,
  region,
  title
order by
  count(arn) desc;
```

### Get the exclusion details of each assessment template that have run at least once
Identify instances where specific assessment templates have been used at least once, and gain insights into the exclusions related to each of these templates. This is useful to understand which templates are commonly used and to review the exclusions associated with them for better resource management.

```sql+postgres
select 
  e.arn, 
  e.title, 
  jsonb_pretty(e.attributes) as attributes, 
  e.recommendation 
from 
  aws_inspector_exclusion e, 
  aws_inspector_assessment_run r, 
  aws_inspector_assessment_template t 
where 
  e.assessment_run_arn = r.arn 
and 
  r.assessment_template_arn = t.arn;
```

```sql+sqlite
select 
  e.arn, 
  e.title, 
  e.attributes, 
  e.recommendation 
from 
  aws_inspector_exclusion e, 
  aws_inspector_assessment_run r, 
  aws_inspector_assessment_template t 
where 
  e.assessment_run_arn = r.arn 
and 
  r.assessment_template_arn = t.arn;
```