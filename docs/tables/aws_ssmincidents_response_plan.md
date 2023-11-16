# Table: aws_ssmincidents_response_plan

AWS SSM Incidents response plan automates the initial response to incidents. A response plan engages contacts, starts chat channel collaboration, and initiates runbooks at the beginning of an incident.

## Examples

### Basic info

```sql
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

```sql
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

### Get incident template details of a partuicular response plan

```sql
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

### Get the details of integrations associated to the response plans

```sql
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

### Get details of engagements associated to the response plans

```sql
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