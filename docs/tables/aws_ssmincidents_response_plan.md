---
title: "Steampipe Table: aws_ssmincidents_response_plan - Query AWS SSM Incidents Response Plan using SQL"
description: "Allows users to query AWS SSM Incidents Response Plan data to retrieve information about each resource plan in your AWS account."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssmincidents_response_plan - Query AWS SSM Incidents Response Plan using SQL

AWS SSM Incidents response plan automates the initial response to incidents. A response plan engages contacts, starts chat channel collaboration, and initiates runbooks at the beginning of an incident.

## Table Usage Guide

The `aws_ssmincidents_response_plan` table in Steampipe allows you to query information about each response plan in your AWS account. This table provides you, as a DevOps engineer or system administrator, with response plan specific details, including the ARN, name, chat channel, incident template, and more. You can utilize this table to gather insights on response plans.

## Examples

### Basic info
Analyze the settings to understand the comprehensive overview of the response plans configured in AWS Systems Manager Incident Manager, aiding in effective incident response and management. This information is particularly useful for assessing the setup of your incident response infrastructure, ensuring that all necessary components are in place and properly configured.

```sql+postgres
select
  name,
  arn,
  display_name,
  chat_channel,
  incident_template,
  integrations,
  title
from
  aws_ssmincidents_response_plan;
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  chat_channel,
  incident_template,
  integrations,
  title
from
  aws_ssmincidents_response_plan;
```

### List response plans with chat channel configured
Focuses on retrieving information about AWS Systems Manager Incident Manager response plans that have an associated chat channel. It helps organizations to assess and manage their incident response strategies, especially in the context of communication readiness and efficiency.

```sql+postgres
select
  name,
  arn,
  display_name,
  chat_channel,
  incident_template,
  integrations,
  title
from
  aws_ssmincidents_response_plan
where
  chat_channel is not null;
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  chat_channel,
  incident_template,
  integrations,
  title
from
  aws_ssmincidents_response_plan
where
  chat_channel is not null;
```

### Get incident template details of a response plan
Retrieve detailed information about a specific AWS Systems Manager Incident Manager response plan, particularly focusing on various aspects of the incident template.

```sql+postgres
select
  name,
  arn,
  display_name,
  incident_template -> 'Impact' as incident_template_impact,
  incident_template -> 'Title' as incident_template_title,
  incident_template -> 'DedupeString',
  incident_template -> 'IncidentTags' as incident_template_tags,
  incident_template -> 'NotificationTargets' as incident_notification_targets,
  incident_template -> 'Summary' as incident_template_summary,
  title
from
  aws_ssmincidents_response_plan
where
  incident_template is not null
  and arn = 'arn:aws:ssm-incidents::111111111111:response-plan/response-plan-test';
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  json_extract(incident_template, '$.Impact') as incident_template_impact,
  json_extract(incident_template, '$.Title') as incident_template_title,
  json_extract(incident_template, '$.DedupeString'),
  json_extract(incident_template, '$.IncidentTags') as incident_template_tags,
  json_extract(incident_template, '$.NotificationTargets') as incident_notification_targets,
  json_extract(incident_template, '$.Summary') as incident_template_summary,
  title
from
  aws_ssmincidents_response_plan
where
  incident_template is not null
  and arn = 'arn:aws:ssm-incidents::111111111111:response-plan/response-plan-test';

```

### Get the details of integrations associated to the response plans
Involves collecting and understanding all pertinent information about the external services and tools that are linked to specific incident response plans, with the aim of understanding how these integrations support the overall incident management process.

```sql+postgres
select
  name,
  arn,
  display_name,
  jsonb_pretty(integrations),
  title
from
  aws_ssmincidents_response_plan
where
  integrations is not null;
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  json_extract(integrations, '$') as integrations,
  title
from
  aws_ssmincidents_response_plan
where
  integrations is not null;

```

### Get details of engagements associated to the response plans
Analyzing the engagements associated with response plans, you can gain insights into how an organization prepares for and manages incident responses. This includes understanding the communication protocols, roles and responsibilities, and other coordination strategies embedded in the response plans.

```sql+postgres
select
  name,
  arn,
  display_name,
  jsonb_pretty(engagements),
  title
from
  aws_ssmincidents_response_plan
where
  engagements is not null;
```

```sql+sqlite
select
  name,
  arn,
  display_name,
  json_extract(engagements, '$') as engagements,
  title
from
  aws_ssmincidents_response_plan
where
  engagements is not null;
```