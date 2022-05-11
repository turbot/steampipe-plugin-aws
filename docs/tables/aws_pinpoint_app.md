# Table: aws_pinpoint_app

In Amazon Pinpoint, a project is a collection of recipient information, segments, campaigns, and journeys. New Amazon Pinpoint users should start by creating a project. Amazon Pinpoint is an AWS service that you can use to engage with your customers across multiple messaging channels. You can use Amazon Pinpoint to send push notifications, in-app notifications, emails, text messages, voice messages, and messages over custom channels. It includes segmentation, campaign, and journey features that help you send the right message to the right customer at the right time over the right channel.
## Examples

### Basic info

```sql
select
  id,
  name,
  arn,
  limits
from
  aws_pinpoint_app;
```

### Get quiet time for application

```sql
select
  id,
  quiet_time -> 'Start' as start_time,
  quiet_time -> 'End' as end_time
from
  aws_pinpoint_app;
```

### Get campaign hook details for application

```sql
select
  id,
  campaign_hook -> 'LambdaFunctionName' as lambda_function_name,
  campaign_hook -> 'Mode' as mode,
  campaign_hook -> 'WebUrl' as web_url
from
  aws_pinpoint_app;
```

### List the limits of applications

```sql
select
  id,
  limits -> 'Daily' as daily,
  limits -> 'Total' as total,
  limits -> 'Session' as session,
  limits -> 'MaximumDuration' as maximum_duration,
  limits -> 'MessagesPerSecond' as messages_per_second
from
  aws_pinpoint_app;
```