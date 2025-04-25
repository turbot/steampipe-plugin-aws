---
title: "Steampipe Table: aws_securityhub_hub - Query AWS Security Hub using SQL"
description: "Allows users to query Security Hub to retrieve information about the Hub resources."
folder: "Security Hub"
---

# Table: aws_securityhub_hub - Query AWS Security Hub using SQL

The AWS Security Hub provides a comprehensive view of your high-priority security alerts and compliance status across your AWS accounts. It aggregates, organizes, and prioritizes your security alerts, or findings, from multiple AWS services, such as Amazon GuardDuty, Amazon Inspector, and Amazon Macie, as well as from AWS Partner solutions. The findings are then visually summarized on integrated dashboards with actionable graphs and tables.

## Table Usage Guide

The `aws_securityhub_hub` table in Steampipe provides you with information about Hub resources within AWS Security Hub. This table allows you, as a DevOps engineer, to query Hub-specific details, including the ARN, subscription status, and auto-enable controls. You can utilize this table to gather insights on Hub resources, such as their subscription status, whether auto-enable controls are activated, and more. The schema outlines the various attributes of the Security Hub for you, including the Hub ARN, auto-enable controls status, and associated tags.

## Examples

### Basic info
Explore which AWS Security Hub settings are automatically enabling controls and when they were subscribed to, across different regions. This can help in managing security protocols and ensuring timely compliance across the organization's AWS infrastructure.

```sql+postgres
select
  hub_arn,
  auto_enable_controls,
  subscribed_at,
  region
from
  aws_securityhub_hub;
```

```sql+sqlite
select
  hub_arn,
  auto_enable_controls,
  subscribed_at,
  region
from
  aws_securityhub_hub;
```

### List hubs that do not automatically enable new controls
Identify hubs within the AWS Security Hub service that have not been configured to automatically enable new controls. This can be useful in assessing the level of manual intervention required in your security management process.

```sql+postgres
select
  hub_arn,
  auto_enable_controls
from
  aws_securityhub_hub
where
  not auto_enable_controls;
```

```sql+sqlite
select
  hub_arn,
  auto_enable_controls
from
  aws_securityhub_hub
where
  auto_enable_controls = 0;
```

### List administrator account details for the hub 
Explore the details of administrator accounts in the security hub, including invitation status and time, to manage and monitor account usage. This is particularly useful in tracking the status of administrator invitations and maintaining security controls.

```sql+postgres
select
  hub_arn,
  auto_enable_controls,
  administrator_account ->> 'AccountId' as administrator_account_id,
  administrator_account ->> 'InvitationId' as administrator_invitation_id,
  administrator_account ->> 'InvitedAt' as administrator_invitation_time,
  administrator_account ->> 'MemberStatus' as administrator_status
from
  aws_securityhub_hub
where
  administrator_account is not null;
```

```sql+sqlite
select
  hub_arn,
  auto_enable_controls,
  json_extract(administrator_account, '$.AccountId') as administrator_account_id,
  json_extract(administrator_account, '$.InvitationId') as administrator_invitation_id,
  json_extract(administrator_account, '$.InvitedAt') as administrator_invitation_time,
  json_extract(administrator_account, '$.MemberStatus') as administrator_status
from
  aws_securityhub_hub
where
  administrator_account is not null;
```