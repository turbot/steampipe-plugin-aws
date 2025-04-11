---
title: "Steampipe Table: aws_docdb_cluster_instance - Query Amazon DocumentDB Cluster Instances using SQL"
description: "Allows users to query Amazon DocumentDB Cluster Instances to gather detailed information such as instance identifier, cluster identifier, instance class, availability zone, engine version, and more."
folder: "DocumentDB"
---

# Table: aws_docdb_cluster_instance - Query Amazon DocumentDB Cluster Instances using SQL

The Amazon DocumentDB Cluster Instance is a part of Amazon DocumentDB, a fast, scalable, highly available, and fully managed document database service that supports MongoDB workloads. It provides the performance, scalability, and availability you need when operating mission-critical MongoDB workloads at scale. With DocumentDB, you can store, query, and index JSON data.

## Table Usage Guide

The `aws_docdb_cluster_instance` table in Steampipe provides you with information about Amazon DocumentDB Cluster Instances. This table allows you as a DevOps engineer, database administrator, or other technical professional to query detailed information about each cluster instance, such as its identifier, associated cluster identifier, instance class, availability zone, engine version, and other relevant metadata. You can utilize this table to gather insights on the configuration, performance, and status of your DocumentDB cluster instances. The schema outlines the various attributes of the DocumentDB cluster instance, including instance ARN, creation time, instance status, and associated tags for you.

## Examples

### Basic info
Gain insights into the specifics of your AWS DocumentDB Cluster instances, such as the engine type, version, and instance class. This can be useful for assessing your current configuration and identifying potential areas for optimization or upgrade.

```sql+postgres
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

```sql+sqlite
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
Identify instances that are accessible to the public, allowing you to review and manage your data's exposure and security. This query is useful for maintaining control over your data privacy and ensuring that only authorized users have access.

```sql+postgres
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

```sql+sqlite
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
  publicly_accessible = 1;
```

### Get DB subnet group information of each instance
Explore the status and details of your database subnet groups across instances to understand their configuration and ensure optimal database management. This is beneficial for maintaining network efficiency and security in your AWS DocumentDB clusters.

```sql+postgres
select
  db_subnet_group_arn,
  db_subnet_group_name,
  db_subnet_group_description,
  db_subnet_group_status
from
  aws_docdb_cluster_instance;
```

```sql+sqlite
select
  db_subnet_group_arn,
  db_subnet_group_name,
  db_subnet_group_description,
  db_subnet_group_status
from
  aws_docdb_cluster_instance;
```

### Get VPC and subnet information of each instance
Determine the areas in which each instance of your database is connected to a VPC and its associated subnet. This is useful for understanding your database's network configuration and ensuring it aligns with your security and performance requirements.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier as attached_vpc,
  json_extract(vsg.value, '$.VpcSecurityGroupId') as vpc_security_group_id,
  json_extract(vsg.value, '$.Status') as status,
  json_extract(json_extract(sub.value, '$.SubnetAvailabilityZone'), '$.Name') as subnet_availability_zone,
  json_extract(sub.value, '$.SubnetIdentifier') as subnet_identifier,
  json_extract(json_extract(sub.value, '$.SubnetOutpost'), '$.Arn') as subnet_outpost,
  json_extract(sub.value, '$.SubnetStatus') as subnet_status
from
  aws_docdb_cluster_instance,
  json_each(vpc_security_groups) as vsg,
  json_each(subnets) as sub;
```

### List instances with unecrypted storage
Identify instances where storage is not encrypted to understand potential vulnerabilities in your database security. This is crucial for ensuring data protection and compliance with security regulations.

```sql+postgres
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  not storage_encrypted;
```

```sql+sqlite
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  storage_encrypted = 0;
```

### List instances with cloudwatch logs disabled
Identify instances where DocumentDB clusters in AWS might be vulnerable due to disabled CloudWatch logs. This query is beneficial for improving security and compliance by ensuring that all instances have logging enabled.

```sql+postgres
select
  db_instance_identifier,
  db_cluster_identifier,
  db_instance_class
from
  aws_docdb_cluster_instance
where
  enabled_cloudwatch_logs_exports is null;
```

```sql+sqlite
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
Gain insights into the network connectivity of each instance by identifying the network endpoint details. This can be beneficial in diagnosing connectivity issues or planning network configurations.

```sql+postgres
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_docdb_cluster_instance;
```

```sql+sqlite
select
  db_instance_identifier,
  endpoint_address,
  endpoint_hosted_zone_id,
  endpoint_port
from
  aws_docdb_cluster_instance;
```