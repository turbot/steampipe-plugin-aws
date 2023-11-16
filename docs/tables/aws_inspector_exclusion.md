---
title: "Table: aws_inspector_exclusion - Query AWS Inspector Exclusions using SQL"
description: "Allows users to query AWS Inspector Exclusions and retrieve data about AWS Inspector Exclusions, including their ARNs, descriptions, and recommendations."
---

# Table: aws_inspector_exclusion - Query AWS Inspector Exclusions using SQL

The `aws_inspector_exclusion` table in Steampipe provides information about exclusions within AWS Inspector. This table allows DevOps engineers to query exclusion-specific details, including the ARN, description, and recommendation. Users can utilize this table to gather insights on exclusions, such as their status, the reasons behind their exclusions, and more. The schema outlines the various attributes of the AWS Inspector exclusion, including the ARN, description, recommendation, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector_exclusion` table, you can use the `.inspect aws_inspector_exclusion` command in Steampipe.

**Key columns**:

- `arn`: The ARN of the exclusion. This is a unique identifier for the exclusion and can be used to join this table with other tables.
- `description`: The description of the exclusion. This provides context about the exclusion and can be useful when analyzing the data.
- `recommendation`: The recommendation for the exclusion. This provides additional information about the exclusion and can be used to understand the reasons behind the exclusion.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  arn,
  jsonb_pretty(attributes) as attributes,
  jsonb_pretty(scopes) as scopes
from
  aws_inspector_exclusion;
```

### Count the number of exclusions whose type is 'Agent not found'

```sql
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

### Get the exclusion details of each assessment template that have run at least once

```sql
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