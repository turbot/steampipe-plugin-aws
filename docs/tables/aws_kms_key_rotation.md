---
title: "Steampipe Table: aws_kms_key_rotation - Query AWS KMS Key Rotation using SQL"
description: "Allows users to query AWS KMS Key Rotation data including the rotation schedule, type, and associated key details."
folder: "KMS"
---

# Table: aws_kms_key_rotation - Query AWS KMS Key Rotation using SQL

The AWS Key Management Service (KMS) includes functionalities for rotating encryption keys, which is crucial for maintaining the security of cryptographic keys over time. The `aws_kms_key_rotation` table provides access to detailed information about the rotation status and history of these keys, enabling enhanced security management and compliance with best practices.

## Table Usage Guide

The `aws_kms_key_rotation` table in Steampipe is useful for security analysts and DevOps engineers to monitor and audit the rotation of AWS KMS keys. It includes key details such as the rotation date, type, and associated key ARN. This table allows you to query information efficiently for regular audits and compliance reporting.

## Examples

### Basic info
Retrieve basic information about key rotations including ARN, rotation date, and type. This can be useful for regular audits of key management practices.

```sql+postgres
select
  key_id,
  key_arn,
  rotation_date,
  rotation_type
from
  aws_kms_key_rotation;
```

```sql+sqlite
select
  key_id,
  key_arn,
  rotation_date,
  rotation_type
from
  aws_kms_key_rotation;
```

### Keys with recent rotations
List details of keys that have undergone rotation within the last 30 days, helping to ensure recent key rotations are tracked for security compliance.

```sql+postgres
select
  key_id,
  key_arn,
  rotation_date
from
  aws_kms_key_rotation
where
  rotation_date >= current_date - interval '30 days';
```

```sql+sqlite
select
  key_id,
  key_arn,
  rotation_date
from
  aws_kms_key_rotation
where
  strftime('%s', 'now') - strftime('%s', rotation_date) <= 2592000;
```

### Join with aws_kms_key to get complete key details
Provide a comprehensive overview of key rotation along with key management details.

```sql+postgres
select
  akr.key_id,
  ak.title,
  akr.rotation_date,
  akr.rotation_type,
  ak.key_manager
from
  aws_kms_key_rotation akr
join
  aws_kms_key ak
on
  akr.key_id = ak.id;
```

```sql+sqlite
select
  akr.key_id,
  ak.title,
  akr.rotation_date,
  akr.rotation_type,
  ak.key_manager
from
  aws_kms_key_rotation akr
join
  aws_kms_key ak
on
  akr.key_id = ak.id;
```

### Count of key rotations by type
This query groups keys by rotation type, providing insights into how many keys are rotated automatically versus on-demand.

```sql+postgres
select
  rotation_type,
  count(key_id) as count
from
  aws_kms_key_rotation
group by
  rotation_type;
```

```sql+sqlite
select
  rotation_type,
  count(key_id) as count
from
  aws_kms_key_rotation
group by
  rotation_type;
```