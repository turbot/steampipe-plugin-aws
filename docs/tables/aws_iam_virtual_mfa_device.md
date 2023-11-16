---
title: "Table: aws_iam_virtual_mfa_device - Query AWS Identity and Access Management (IAM) Virtual MFA Devices using SQL"
description: "Allows users to query Virtual MFA Devices in AWS Identity and Access Management (IAM)."
---

# Table: aws_iam_virtual_mfa_device - Query AWS Identity and Access Management (IAM) Virtual MFA Devices using SQL

The `aws_iam_virtual_mfa_device` table in Steampipe provides information about Virtual MFA Devices within AWS Identity and Access Management (IAM). This table allows security administrators and compliance auditors to query device-specific details, including the device ARN, enablement status, and associated users. Users can utilize this table to gather insights on MFA devices, such as which devices are enabled, which users are associated with a particular device, and more. The schema outlines the various attributes of the Virtual MFA Device, including the device ARN, enable date, serial number, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_virtual_mfa_device` table, you can use the `.inspect aws_iam_virtual_mfa_device` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the virtual MFA device. This column can be used to join this table with other tables that contain IAM resource information.

- `serial_number`: The serial number that uniquely identifies a virtual MFA device. This can be used to join with other tables that may store device-specific data.

- `user`: The IAM user associated with the virtual MFA device. This column can be used to join this table with other tables that contain user-specific information, providing a comprehensive view of user security configurations.

## Examples

### Basic info

```sql
select
  serial_number,
  enable_date,
  user_name
from
  aws_iam_virtual_mfa_device;
```

### User details for users with a virtual MFA device assigned

```sql
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
