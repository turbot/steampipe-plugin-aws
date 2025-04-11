---
title: "Steampipe Table: aws_iam_virtual_mfa_device - Query AWS Identity and Access Management (IAM) Virtual MFA Devices using SQL"
description: "Allows users to query Virtual MFA Devices in AWS Identity and Access Management (IAM)."
folder: "IAM"
---

# Table: aws_iam_virtual_mfa_device - Query AWS Identity and Access Management (IAM) Virtual MFA Devices using SQL

An AWS Identity and Access Management (IAM) Virtual Multi-Factor Authentication (MFA) Device is a two-step authentication process that enhances security by requiring users to authenticate using a virtual or hardware device that produces a six-digit numeric code. The virtual MFA device uses a software application that generates one-time passwords (OTP) on a device such as a smartphone. This provides an additional layer of security for AWS service access, making it more difficult for attackers to compromise accounts.

## Table Usage Guide

The `aws_iam_virtual_mfa_device` table in Steampipe provides you with information about Virtual MFA Devices within AWS Identity and Access Management (IAM). This table allows you, as a security administrator or compliance auditor, to query device-specific details, including the device ARN, enablement status, and associated users. You can utilize this table to gather insights on MFA devices, such as which devices are enabled, which users are associated with a particular device, and more. The schema outlines the various attributes of the Virtual MFA Device for you, including the device ARN, enable date, serial number, and associated tags.

## Examples

### Basic info
Explore which multi-factor authentication devices are enabled on your AWS account and who they are assigned to. This information can help you manage security settings and ensure that only authorized users have access.

```sql+postgres
select
  serial_number,
  enable_date,
  user_name
from
  aws_iam_virtual_mfa_device;
```

```sql+sqlite
select
  serial_number,
  enable_date,
  user_name
from
  aws_iam_virtual_mfa_device;
```

### User details for users with a virtual MFA device assigned
Explore which users have a virtual Multi-Factor Authentication (MFA) device assigned. This is useful to ensure all users are following security best practices and have an additional layer of security enabled.

```sql+postgres
select
  name,
  u.user_id,
  mfa.serial_number,
  path,
  create_date,
  password_last_used
from
  aws_iam_user u
  inner join aws_iam_virtual_mfa_device mfa on u.name = mfa.user_name;
```

```sql+sqlite
select
  name,
  u.user_id,
  mfa.serial_number,
  path,
  create_date,
  password_last_used
from
  aws_iam_user u
  join aws_iam_virtual_mfa_device mfa on u.name = mfa.user_name;
```