---
title: "Steampipe Table: aws_kms_key - Query AWS KMS Key using SQL"
description: "Allows users to query AWS KMS Key data including cryptographic details, key usage, key state, and associated metadata."
folder: "KMS"
---

# Table: aws_kms_key - Query AWS KMS Key using SQL

The AWS Key Management Service (KMS) is a managed service that makes it easy for you to create and control the cryptographic keys used to encrypt your data. It provides a highly available key storage, management, and auditing solution for you to encrypt data within your own applications and control the encryption of stored data across AWS services. With AWS KMS, you can protect your keys from unauthorized use by defining key policies and IAM policies.

## Table Usage Guide

The `aws_kms_key` table in Steampipe provides you with information about Key Management Service (KMS) keys within AWS. This table allows you, as a DevOps engineer, to query key-specific details, including cryptographic details, key usage, key state, and associated metadata. You can utilize this table to gather insights on keys, such as keys rotation status, key type, key state, and more. The schema outlines the various attributes of the KMS key for you, including the key ARN, creation date, key state, key usage, and associated tags.

## Examples

### Basic info
Explore which AWS Key Management Service (KMS) keys have been created and who manages them. This is useful for auditing security practices and understanding the distribution of access control within your AWS environment.

```sql+postgres
select
  id,
  title,
  arn,
  key_manager,
  creation_date
from
  aws_kms_key;
```

```sql+sqlite
select
  id,
  title,
  arn,
  key_manager,
  creation_date
from
  aws_kms_key;
```

### List of KMS keys where key rotation is not enabled
Identify instances where key rotation is not enabled for your AWS KMS keys. This can help enhance your security posture by revealing keys that may be vulnerable due to lack of regular rotation.

```sql+postgres
select
  id,
  key_rotation_enabled
from
  aws_kms_key
where
  not key_rotation_enabled;
```

```sql+sqlite
select
  id,
  key_rotation_enabled
from
  aws_kms_key
where
  key_rotation_enabled = 0;
```


### List of KMS Customer Managed keys scheduled for deletion
Identify instances where customer-managed keys are scheduled for deletion in the AWS Key Management Service. This can help in managing key lifecycle and preventing accidental loss of access to AWS resources.

```sql+postgres
select
  id,
  key_state,
  deletion_date
from
  aws_kms_key
where
  key_state = 'PendingDeletion';
```

```sql+sqlite
select
  id,
  key_state,
  deletion_date
from
  aws_kms_key
where
  key_state = 'PendingDeletion';
```


### List of unused Customer Managed Keys
Discover the segments that consist of inactive customer-managed keys, enabling you to identify potential areas for resource optimization and enhanced security management.

```sql+postgres
select
  id,
  enabled as key_enabled
from
  aws_kms_key
where
  not enabled;
```

```sql+sqlite
select
  id,
  enabled as key_enabled
from
  aws_kms_key
where
  enabled = 0;
```


### Count of AWS KMS keys by Key manager
Discover the segments that utilize AWS Key Management Service (KMS) keys by grouping them according to their key manager. This can provide insights into the distribution and management of your encryption keys, aiding in security audits and compliance reviews.

```sql+postgres
select
  key_manager,
  count(key_manager) as count
from
  aws_kms_key
group by
  key_manager;
```

```sql+sqlite
select
  key_manager,
  count(key_manager) as count
from
  aws_kms_key
group by
  key_manager;
```