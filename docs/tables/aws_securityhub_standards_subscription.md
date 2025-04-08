---
title: "Steampipe Table: aws_securityhub_standards_subscription - Query AWS Security Hub Standards Subscriptions using SQL"
description: "Allows users to query AWS Security Hub Standards Subscriptions, providing detailed information about each standard subscription in AWS Security Hub."
folder: "RDS"
---

# Table: aws_securityhub_standards_subscription - Query AWS Security Hub Standards Subscriptions using SQL

The AWS Security Hub Standards Subscriptions is a feature of AWS Security Hub that allows you to manage and implement security standards in your AWS environments. It provides a comprehensive view of your high-priority security alerts and compliance status across AWS accounts. This enables you to quickly assess your security and compliance status, identify potential issues, and take necessary actions to maintain a secure and compliant environment.

## Table Usage Guide

The `aws_securityhub_standards_subscription` table in Steampipe provides you with information about standards subscriptions within AWS Security Hub. This table allows you, as a DevOps engineer, to query subscription-specific details, including the standard's ARN, name, description, and compliance status. You can utilize this table to gather insights on standards, such as their status, updates, and the regions in which they are enabled. The schema outlines the various attributes of the standards subscription for you, including the standards ARN, status, and enabled timestamp.

## Examples

### Basic info
Explore which security standards are currently subscribed to within your AWS SecurityHub across different regions. This can help you assess your security posture and ensure compliance with necessary standards.

```sql+postgres
select
  name,
  standards_arn,
  description,
  region
from
  aws_securityhub_standards_subscription;
```

```sql+sqlite
select
  name,
  standards_arn,
  description,
  region
from
  aws_securityhub_standards_subscription;
```


### List enabled security hub standards
Discover the segments that are automatically safeguarded by identifying the activated security standards within AWS Security Hub. This is beneficial for understanding which areas of your system are already protected by default settings, helping to inform decisions on additional security measures.

```sql+postgres
select
  name,
  standards_arn,
  enabled_by_default
from
  aws_securityhub_standards_subscription
where
  enabled_by_default;
```

```sql+sqlite
select
  name,
  standards_arn,
  enabled_by_default
from
  aws_securityhub_standards_subscription
where
  enabled_by_default = 1;
```


### List standards whose status is not ready
Explore which security standards are not in a 'ready' state. This is useful to identify potential issues or delays in your security configuration.

```sql+postgres
select
  name,
  standards_arn,
  standards_subscription_arn,
  standards_status,
  standards_status_reason_code
from
  aws_securityhub_standards_subscription
where
  standards_status <> 'READY';
```

```sql+sqlite
select
  name,
  standards_arn,
  standards_subscription_arn,
  standards_status,
  standards_status_reason_code
from
  aws_securityhub_standards_subscription
where
  standards_status != 'READY';
```

### List standards that are not managed by AWS
Determine the areas in which security standards are not managed by AWS, allowing you to pinpoint specific locations where alternative management strategies are in effect.

```sql+postgres
select
  name,
  standards_arn,
  standards_managed_by ->> 'Company' as standards_managed_by_company
from
  aws_securityhub_standards_subscription
where
  standards_managed_by ->> 'Company' <> 'AWS';
```

```sql+sqlite
select
  name,
  standards_arn,
  json_extract(standards_managed_by, '$.Company') as standards_managed_by_company
from
  aws_securityhub_standards_subscription
where
  json_extract(standards_managed_by, '$.Company') <> 'AWS';
```