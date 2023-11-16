---
title: "Table: aws_docdb_cluster_instance - Query Amazon DocumentDB Cluster Instances using SQL"
description: "Allows users to query Amazon DocumentDB Cluster Instances to gather detailed information such as instance identifier, cluster identifier, instance class, availability zone, engine version, and more."
---

# Table: aws_docdb_cluster_instance - Query Amazon DocumentDB Cluster Instances using SQL

The `aws_docdb_cluster_instance` table in Steampipe provides information about Amazon DocumentDB Cluster Instances. This table allows DevOps engineers, database administrators, and other technical professionals to query detailed information about each cluster instance, such as its identifier, associated cluster identifier, instance class, availability zone, engine version, and other relevant metadata. Users can utilize this table to gather insights on the configuration, performance, and status of their DocumentDB cluster instances. The schema outlines the various attributes of the DocumentDB cluster instance, including instance ARN, creation time, instance status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_docdb_cluster_instance` table, you can use the `.inspect aws_docdb_cluster_instance` command in Steampipe.

Key columns:

- `db_instance_identifier`: This is the unique identifier for the instance. It's useful for joining with other tables that contain instance-specific data.
- `db_cluster_identifier`: This is the identifier of the cluster that the instance is part of. It can be used to join with other tables that contain cluster-specific information.
- `availability_zone`: This column contains the name of the availability zone the instance is located in. It's useful for joining with other tables that contain availability zone-specific data.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  engine,
  engine_version,
  db_instance_class,
  availability_zone
from
  aws_docdb_cluster_instance;
```

### List instances which are publicly accessible

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  engine,
  engine_version,
  db_instance_class,
  availability_zone
from
  aws_docdb_cluster_instance
where
  publicly_accessible;
```

### Get DB subnet group information of each instance

```sql
select
  db_subnet_group_arn,
  db_subnet_group_name,
  db_subnet_group_description,
  db_subnet_group_status
from
  aws_docdb_cluster_instance;
```

### Get VPC and subnet information of each instance

```sql
select
  db_instance_identifier as attached_vpc,
  vsg ->> 'VpcSecurityGroupId' as vpc_security_group_id,
  vsg ->> 'Status' as status,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub -> 'SubnetOutpost' ->> 'Arn' as subnet_outpost,
  sub ->> 'SubnetStatus' as subnet_status
from
  aws_docdb_cluster_instance
  cross join jsonb_array_elements(vpc_security_groups) as vsg
  cross join jsonb_array_elements(subnets) as sub;
```

### List instances with unecrypted storage

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  not storage_encrypted;
```

### List instances with cloudwatch logs disabled

```sql
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  enabled_cloudwatch_logs_exports is null;
```

### Get network endpoint information of each instance

```sql
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_docdb_cluster_instance;
```
