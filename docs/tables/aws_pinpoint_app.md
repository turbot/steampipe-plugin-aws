# Table: aws_pinpoint_app

In Amazon Pinpoint, a project is a collection of recipient information, segments, campaigns, and journeys. New Amazon Pinpoint users should start by creating a project.

## Examples

### Basic info

```sql
select
  id,
  name,
  arn,
  application_settings
from
  aws_pinpoint_app;
```


### Application settings details

```sql
select
  id,
  application_settings ->> 'Limits' as limits,
  application_settings ->> 'QuietTime' as quiet_time,
  application_settings ->> 'CampaignHook' as campaign_hook,
  application_settings ->> 'LastModifiedDate' as last_modified_date
from
  aws_pinpoint_app;
```

### List limit of an application

```sql
select
  id,
  application_settings -> 'Limits' ->> 'Daily' as daily,
  application_settings -> 'Limits' ->> 'Total' as total,
  application_settings -> 'Limits' ->> 'Session' as session,
  application_settings -> 'Limits' ->> 'MaximumDuration' as maximum_duration,
  application_settings -> 'Limits' ->> 'MessagesPerSecond' as messages_per_second
from
  aws_pinpoint_app;
```