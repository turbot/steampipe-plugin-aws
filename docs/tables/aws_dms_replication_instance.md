---
title: "Steampipe Table: aws_dms_replication_instance - Query AWS Database Migration Service Replication Instances using SQL"
description: "Allows users to query AWS Database Migration Service Replication Instances and provides information about each replication instance in an AWS DMS (Database Migration Service)."
folder: "DMS"
---

# Table: aws_dms_replication_instance - Query AWS Database Migration Service Replication Instances using SQL

The AWS Database Migration Service Replication Instances are fully managed, serverless instances that enable the migration of data from one type of database to another. They facilitate homogeneous or heterogeneous migrations and can handle continuous data replication with high availability and consolidated auditing. This service significantly simplifies the process of migrating existing data to AWS in a secure and efficient manner.

## Table Usage Guide

The `aws_dms_replication_instance` table in Steampipe provides you with information about each replication instance in an AWS Database Migration Service. This table allows you, as a database administrator, to query replication-specific details, including engine version, instance class, allocated storage, and associated metadata. You can utilize this table to gather insights on replication instances, such as their current state, multi-AZ mode, publicly accessible status, and more. The schema outlines the various attributes of the replication instance, including the replication instance ARN, replication instance identifier, availability zone, and associated tags for you.

## Examples

### Basic info
Explore which replication instances in your AWS Database Migration Service have public accessibility. This can help identify potential security risks and ensure that your data is properly protected.

```sql+postgres
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  kms_key_id,
  publicly_accessible,
  region
from
  aws_dms_replication_instance;
```

```sql+sqlite
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  kms_key_id,
  publicly_accessible,
  region
from
  aws_dms_replication_instance;
```


### List replication instances with auto minor version upgrades disabled
Determine the areas in which replication instances have automatic minor version upgrades turned off. This is useful for identifying potential security risks or outdated systems that may require manual updates.

```sql+postgres
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  auto_minor_version_upgrade,
  region
from
  aws_dms_replication_instance
where
  not auto_minor_version_upgrade;
```

```sql+sqlite
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  auto_minor_version_upgrade,
  region
from
  aws_dms_replication_instance
where
  auto_minor_version_upgrade = 0;
```

### List replication instances provisioned with undesired (for example, dms.r5.16xlarge and dms.r5.24xlarge are not desired) instance classes
Determine the areas in which replication instances are provisioned with instance classes that are not preferred, such as dms.r5.16xlarge and dms.r5.24xlarge. This enables you to identify and rectify instances that may not meet your specific requirements or standards.

```sql+postgres
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  replication_instance_class,
  region
from
  aws_dms_replication_instance
where
  replication_instance_class not in ('dms.r5.16xlarge', 'dms.r5.24xlarge');
```

```sql+sqlite
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  replication_instance_class,
  region
from
  aws_dms_replication_instance
where
  replication_instance_class not in ('dms.r5.16xlarge', 'dms.r5.24xlarge');
```


### List publicly accessible replication instances
Determine the areas in which replication instances are publicly accessible. This can help enhance security by identifying potential vulnerabilities in your system.

```sql+postgres
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  region
from
  aws_dms_replication_instance
where
  publicly_accessible;
```

```sql+sqlite
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  region
from
  aws_dms_replication_instance
where
  publicly_accessible = 1;
```


### List replication instances not using multi-AZ deployment configurations
Identify instances where the replication process is not utilizing multi-AZ deployment configurations. This query is beneficial for pinpointing potential areas of vulnerability in your system, as it highlights where redundancies may not be in place to prevent data loss in the event of an AZ outage.

```sql+postgres
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  multi_az,
  region
from
  aws_dms_replication_instance
where
  not multi_az;
```

```sql+sqlite
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  multi_az,
  region
from
  aws_dms_replication_instance
where
  multi_az = 0;
```