---
title: "Steampipe Table: aws_securityhub_member - Query AWS Security Hub Members using SQL"
description: "Allows users to query AWS Security Hub Members for detailed information about each member's account, including its ID, email, status, and more."
folder: "Security Hub"
---

# Table: aws_securityhub_member - Query AWS Security Hub Members using SQL

The AWS Security Hub Members are a part of AWS Security Hub service that allows you to manage and improve the security of your AWS accounts. It aggregates, organizes, and prioritizes your security alerts, or findings, from multiple AWS services, such as Amazon GuardDuty, Amazon Inspector, and Amazon Macie, as well as from AWS Partner solutions. The Members feature specifically enables you to add accounts to and manage accounts in your Security Hub administrator account.

## Table Usage Guide

The `aws_securityhub_member` table in Steampipe provides you with information about each member account within AWS Security Hub. This table allows you, as a DevOps engineer, to query member-specific details, including account ID, email, status, and the timestamp of the invitation. You can utilize this table to gather insights on member accounts, such as their invitation and verification status, the email associated with each account, and more. The schema outlines the various attributes of the Security Hub member for you, including the member account ID, email, status, invited at timestamp, and updated at timestamp.

## Examples

### Basic info
This query allows you to gain insights into the status and administrative details of member accounts within AWS Security Hub. It's useful for monitoring account activity and managing security settings.

```sql+postgres
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member;
```

```sql+sqlite
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member;
```

### List members which are enabled
Explore which members in your AWS Security Hub have their status enabled. This can be useful in maintaining security standards by ensuring only authorized members have access.

```sql+postgres
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled';
```

```sql+sqlite
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled';
```

### List members which are invited but did not accept
Discover the segments of your network where members have been invited but have yet to accept. This can be useful for tracking pending invitations and identifying potential issues with user engagement or notification delivery.

```sql+postgres
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member
where
  member_status = 'Created';
```

```sql+sqlite
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at
from
  aws_securityhub_member
where
  member_status = 'Created';
```

### List members which were invited within the last 10 days
Determine the members who have been enabled and invited to your AWS Security Hub within the past ten days. This can help keep track of recent additions and manage your security operations effectively.

```sql+postgres
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled'
and
  invited_at <= (now() - interval '10' day);
```

```sql+sqlite
select
  member_account_id,
  email,
  administrator_id,
  member_status,
  updated_at,
  invited_at
from
  aws_securityhub_member
where
  member_status = 'Enabled'
and
  invited_at <= datetime('now','-10 day');
```