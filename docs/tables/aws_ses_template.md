---
title: "Steampipe Table: aws_ses_template - Query AWS SES Templates using SQL"
description: "Allows users to query AWS SES Templates for detailed information about email templates used in Amazon Simple Email Service."
folder: "SES"
---

# Table: aws_ses_template - Query AWS SES Templates using SQL

Amazon Simple Email Service (SES) Templates are reusable email templates that can be used to send personalized emails. These templates can include both HTML and text content, along with subject lines, making it easier to maintain consistent email formatting and content across your application.

## Table Usage Guide

The `aws_ses_template` table provides information about email templates in Amazon SES. You can use this table to:

- List all available email templates
- Get detailed information about a specific template, including its HTML and text content
- Monitor template creation and update timestamps
- Track template usage across your AWS account

## Examples

### Basic info
```sql+postgres
select
  name,
  subject_part,
  created_timestamp,
  last_updated_timestamp
from
  aws_ses_template;
```

```sql+sqlite
select
  name,
  subject_part,
  created_timestamp,
  last_updated_timestamp
from
  aws_ses_template;
```

### Get template details by name
```sql+postgres
select
  name,
  subject_part,
  text_part,
  html_part,
  created_timestamp
from
  aws_ses_template
where
  name = 'my-template';
```

```sql+sqlite
select
  name,
  subject_part,
  text_part,
  html_part,
  created_timestamp
from
  aws_ses_template
where
  name = 'my-template';
```

### List templates created in the last 30 days
```sql+postgres
select
  name,
  subject_part,
  created_timestamp
from
  aws_ses_template
where
  created_timestamp >= now() - interval '30 days'
order by
  created_timestamp desc;
```

```sql+sqlite
select
  name,
  subject_part,
  created_timestamp
from
  aws_ses_template
where
  created_timestamp >= datetime('now', '-30 days')
order by
  created_timestamp desc;
```

### List templates with HTML content
```sql+postgres
select
  name,
  subject_part,
  html_part
from
  aws_ses_template
where
  html_part is not null;
```

```sql+sqlite
select
  name,
  subject_part,
  html_part
from
  aws_ses_template
where
  html_part is not null;
```

### List templates with text-only content
```sql+postgres
select
  name,
  subject_part,
  text_part
from
  aws_ses_template
where
  html_part is null
  and text_part is not null;
```

```sql+sqlite
select
  name,
  subject_part,
  text_part
from
  aws_ses_template
where
  html_part is null
  and text_part is not null;
``` 