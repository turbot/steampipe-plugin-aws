---
title: "Steampipe Table: aws_qldb_ledger - Query AWS QLDB Ledgers using SQL"
description: "Allows users to query AWS QLDB ledgers, providing detailed information on ledger configurations, encryption, and permissions."
---

# Table: aws_qldb_ledger - Query AWS QLDB Ledgers using SQL

Amazon QLDB (Quantum Ledger Database) is a fully managed ledger database that provides a transparent, immutable, and cryptographically verifiable transaction log. The `aws_qldb_ledger` table in Steampipe allows you to query information about QLDB ledgers in your AWS environment. This includes details such as ledger state, encryption settings, creation time, and more.

## Table Usage Guide

The `aws_qldb_ledger` table enables cloud administrators and DevOps engineers to gather detailed insights into their QLDB ledgers. You can query various aspects of the ledger, such as its encryption status, deletion protection, and permissions mode. This table is particularly useful for monitoring ledger health, ensuring security compliance, and managing ledger configurations.

## Examples

### Basic ledger information
Retrieve basic information about your AWS QLDB ledgers, including their name, ARN, state, and region. This can be useful for getting an overview of the ledgers deployed in your AWS account.

```sql+postgres
select
  name,
  arn,
  state,
  creation_time,
  region
from
  aws_qldb_ledger;
```

```sql+sqlite
select
  name,
  arn,
  state,
  creation_time,
  region
from
  aws_qldb_ledger;
```

### List active ledgers
Fetch a list of ledgers that are currently active. This can help in identifying which ledgers are operational and available for use.

```sql+postgres
select
  name,
  arn,
  state
from
  aws_qldb_ledger
where
  state = 'ACTIVE';
```

```sql+sqlite
select
  name,
  arn,
  state
from
  aws_qldb_ledger
where
  state = 'ACTIVE';
```

### List ledgers with deletion protection enabled
Identify ledgers that have deletion protection enabled, ensuring that critical data is not accidentally deleted.

```sql+postgres
select
  name,
  arn,
  deletion_protection
from
  aws_qldb_ledger
where
  deletion_protection = true;
```

```sql+sqlite
select
  name,
  arn,
  deletion_protection
from
  aws_qldb_ledger
where
  deletion_protection = 1;
```

### List ledgers by encryption status
Retrieve ledgers based on their encryption status, which can be useful for ensuring compliance with data security standards.

```sql+postgres
select
  name,
  arn,
  encryption_status,
  kms_key_arn
from
  aws_qldb_ledger
where
  encryption_status = 'ENABLED';
```

```sql+sqlite
select
  name,
  arn,
  encryption_status,
  kms_key_arn
from
  aws_qldb_ledger
where
  encryption_status = 'ENABLED';
```

### List ledgers with inaccessible KMS keys
Identify ledgers where the KMS key has become inaccessible, which could indicate potential issues with encryption or key management.

```sql+postgres
select
  name,
  arn,
  inaccessible_kms_key_date_time
from
  aws_qldb_ledger
where
  inaccessible_kms_key_date_time is not null;
```

```sql+sqlite
select
  name,
  arn,
  inaccessible_kms_key_date_time
from
  aws_qldb_ledger
where
  inaccessible_kms_key_date_time is not null;
```

### List ledgers by permissions mode
Fetch a list of ledgers along with their permissions mode, which is useful for understanding the access control configuration of your ledgers.

```sql+postgres
select
  name,
  arn,
  permissions_mode
from
  aws_qldb_ledger;
```

```sql+sqlite
select
  name,
  arn,
  permissions_mode
from
  aws_qldb_ledger;
```