---
title: "Table: aws_securityhub_hub - Query AWS Security Hub using SQL"
description: "Allows users to query Security Hub to retrieve information about the Hub resources."
---

# Table: aws_securityhub_hub - Query AWS Security Hub using SQL

The `aws_securityhub_hub` table in Steampipe provides information about Hub resources within AWS Security Hub. This table allows DevOps engineers to query Hub-specific details, including the ARN, subscription status, and auto-enable controls. Users can utilize this table to gather insights on Hub resources, such as their subscription status, whether auto-enable controls are activated, and more. The schema outlines the various attributes of the Security Hub, including the Hub ARN, auto-enable controls status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_hub` table, you can use the `.inspect aws_securityhub_hub` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the Hub. This is a key column as it uniquely identifies each Hub resource and can be used to join this table with other tables.
- `auto_enable_controls`: Indicates whether the Hub automatically enables new controls when they are added to standards that are enabled. This is useful for understanding the configuration of the Hub.
- `subscription_status`: The subscription status of the Hub. This is important for understanding the status of the Hub in the AWS Security Hub service.

## Examples

### Basic info

```sql
select
  hub_arn,
  auto_enable_controls,
  subscribed_at,
  region
from
  aws_securityhub_hub;
```


### List hubs that do not automatically enable new controls

```sql
select
  hub_arn,
  auto_enable_controls
from
  aws_securityhub_hub
where
  not auto_enable_controls;
```

### List administrator account details for the hub 

```sql
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