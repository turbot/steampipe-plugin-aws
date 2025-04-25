---
title: "Steampipe Table: aws_shield_emergency_contact - Query AWS Shield Advanced Emergency Contacts using SQL"
description: "Allows users to query AWS Shield Advanced Emergency Contacts and retrieve the phone number, email address and notes for each contact."
folder: "Shield"
---

# Table: aws_shield_emergency_contact - Query AWS Shield Advanced Emergency Contacts using SQL

AWS Shield Advanced is a DDoS protection service from AWS. The Emergency Contacts settings allow you to configure the phone numbers, email addresses and notes for the emergency contacts the Shield Response Team (SRT) should reach out to in case of a DDoS attack.

## Table Usage Guide

The `aws_shield_emergency_contact` table in Steampipe allows you to query the AWS Shield Advanced Emergency Contacts and retrieve the phone number, email address and notes for each contact. For more details about the individual fields, please refer to the [AWS Shield Advanced API documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeEmergencyContactSettings.html).

## Examples

### Basic info

```sql+postgres
select
  email_address,
  phone_number,
  contact_notes
from
  aws_shield_emergency_contact;
```

```sql+sqlite
select
  email_address,
  phone_number,
  contact_notes
from
  aws_shield_emergency_contact;
```
