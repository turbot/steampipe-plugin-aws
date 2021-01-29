# Table: aws_iam_virtual_mfa_device

A software app that runs on a phone or other device and emulates a physical device. The device generates a six-digit numeric code based upon a time-synchronized one-time password algorithm.
MFA adds extra security because it requires users to provide unique authentication from an AWS supported MFA mechanism in addition to their regular sign-in credentials when they access AWS websites or services.

## Examples

### Basic Virtual MFA Device info

```sql
select
  serial_number,
  enable_date,
  user_name
from
  aws_iam_virtual_mfa_device;
```

### User details to which the Virtual MFA device is assigned

```sql
select
  name,
  u.user_id,
  mfa.serial_number,
  path,
  create_date,
  password_last_used
from
  aws_iam_user  u
inner join aws_iam_virtual_mfa_device mfa
  on u.name = mfa.user_name;
```
