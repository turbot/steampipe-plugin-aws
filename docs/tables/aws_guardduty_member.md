---
title: "Steampipe Table: aws_guardduty_member - Query AWS GuardDuty Member using SQL"
description: "Allows users to query AWS GuardDuty Member data, including member account details, detector ID, invitation status, and relationship status."
folder: "GuardDuty"
---

# Table: aws_guardduty_member - Query AWS GuardDuty Member using SQL

The AWS GuardDuty Member is a component of the Amazon GuardDuty service which is a threat detection service that continuously monitors for malicious activity and unauthorized behavior to protect your AWS accounts and workloads. A member represents the accounts added to the GuardDuty service from the master account. It helps in managing and organizing multiple AWS accounts for threat detection and notifications.

## Table Usage Guide

The `aws_guardduty_member` table in Steampipe provides you with information about member accounts within AWS GuardDuty. This table allows you, as a security analyst, to query member-specific details, including account details, detector ID, invitation status, and relationship status. You can utilize this table to gather insights on member accounts, such as the status of invitations sent to these accounts, the relationship status between the master and member accounts, and more. The schema outlines the various attributes of the GuardDuty member for you, including the account ID, email, detector ID, and associated tags.

## Examples

### Basic info
Explore which member accounts are linked to your AWS GuardDuty detectors and when they were invited, to understand the security relationships within your network. This can be useful in assessing the overall security posture and identifying any potential weak points.

```sql+postgres
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member;
```

```sql+sqlite
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member;
```

### List members that failed email verification
Uncover the details of members who have not successfully completed the email verification process. This is particularly useful for identifying potential security issues and ensuring all users have been properly validated.

```sql+postgres
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  relationship_status = 'EmailVerificationFailed';
```

```sql+sqlite
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  relationship_status = 'EmailVerificationFailed';
```

### List uninvited members
Identify instances where members of the AWS GuardDuty service have not been invited. This is useful for maintaining security and ensuring all members have been properly onboarded.

```sql+postgres
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at is null;
```

```sql+sqlite
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at is null;
```

### List members which were invited within the last 10 days
Identify newly invited members in the AWS GuardDuty service over the past ten days to monitor recent additions and their relationship status.

```sql+postgres
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at >= (now() - interval '10' day);
```

```sql+sqlite
select
  member_account_id,
  detector_id,
  invited_at,
  relationship_status
from
  aws_guardduty_member
where
  invited_at >= datetime('now', '-10 days');
```