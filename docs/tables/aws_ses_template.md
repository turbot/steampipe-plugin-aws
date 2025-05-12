---
title: "Steampipe Table: aws_ses_template - Query AWS SES Templates using SQL"
description: "Allows users to query AWS SES Templates and retrieve detailed information about each template, including its name, subject, content, and timestamps."
folder: "SES"
---

# Table: aws_ses_template - Query AWS SES Templates using SQL

The AWS SES Template is a feature of Amazon Simple Email Service (SES) that allows you to create and manage email templates. These templates can include both HTML and text content, along with variables that can be replaced with actual values when sending emails. Templates help maintain consistent branding and messaging across your email communications while simplifying the email sending process.

## Table Usage Guide

The `aws_ses_template` table in Steampipe provides you with information about email templates within AWS Simple Email Service (SES). This table allows you, as a DevOps engineer or email administrator, to query template-specific details, including the template name, subject line, HTML and text content, creation timestamp, and last update timestamp. You can utilize this table to gather insights on your email templates, such as their content, creation dates, and usage patterns. The schema outlines the various attributes of the SES template for you, including the template name, subject part, text part, HTML part, and associated metadata.

## Examples

### Basic info
Explore the basic information of your AWS SES templates, including their names, subjects, and creation timestamps. This can help you maintain an overview of your email templates and their metadata.

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
Retrieve detailed information about a specific template, including its subject line and content. This is useful for reviewing or updating specific templates.

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
  name = 'MyTemplate';
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
  name = 'MyTemplate';
```

### List templates created in the last 30 days
Identify recently created templates to track new additions to your email template library. This can help in maintaining an up-to-date inventory of your email templates.

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
Find templates that include HTML formatting, which is useful for identifying templates that require special rendering or may need additional testing across different email clients.

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
Identify templates that contain only plain text content, which can be useful for ensuring accessibility and compatibility with all email clients.

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