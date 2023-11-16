---
title: "Table: aws_securityhub_standards_subscription - Query AWS Security Hub Standards Subscriptions using SQL"
description: "Allows users to query AWS Security Hub Standards Subscriptions, providing detailed information about each standard subscription in AWS Security Hub."
---

# Table: aws_securityhub_standards_subscription - Query AWS Security Hub Standards Subscriptions using SQL

The `aws_securityhub_standards_subscription` table in Steampipe provides information about standards subscriptions within AWS Security Hub. This table allows DevOps engineers to query subscription-specific details, including the standard's ARN, name, description, and compliance status. Users can utilize this table to gather insights on standards, such as their status, updates, and the regions in which they are enabled. The schema outlines the various attributes of the standards subscription, including the standards ARN, status, and enabled timestamp.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_standards_subscription` table, you can use the `.inspect aws_securityhub_standards_subscription` command in Steampipe.

**Key columns**:

- `standards_arn`: The ARN of a standard. This can be used to join with other tables that reference AWS Security Hub standards.
- `standards_subscription_arn`: The ARN of a standards subscription. This can be used to join with other tables that reference AWS Security Hub standards subscriptions.
- `standards_status`: The status of the standards subscription. This is useful for querying the compliance status of various standards.

## Examples

### Basic info

```sql
select
  name,
  standards_arn,
  description,
  region
from
  aws_securityhub_standards_subscription;
```


### List enabled security hub standards

```sql
select
  name,
  standards_arn,
  enabled_by_default
from
  aws_securityhub_standards_subscription
where
  enabled_by_default;
```


### List standards whose status is not ready

```sql
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

### List standards that are not managed by AWS

```sql
select
  name,
  standards_arn,
  standards_managed_by ->> 'Company' as standards_managed_by_company
from
  aws_securityhub_standards_subscription
where
  standards_managed_by ->> 'Company' <> 'AWS';
```