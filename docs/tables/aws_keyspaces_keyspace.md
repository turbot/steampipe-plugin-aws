---
title: "Steampipe Table: aws_keyspaces_keyspace - Query AWS Keyspaces Keyspaces using SQL"
description: "Allows users to query AWS Keyspaces keyspaces, providing detailed information on replication strategies, creation timestamps, and regions."
folder: "Keyspaces"
---

# Table: aws_keyspaces_keyspace - Query AWS Keyspaces Keyspaces using SQL

Amazon Keyspaces (for Apache Cassandra) is a scalable, highly available, and managed Apache Cassandra-compatible database service. It enables you to run Cassandra workloads on AWS without managing the underlying infrastructure. The `aws_keyspaces_keyspace` table in Steampipe allows you to query information about your Keyspaces keyspaces in AWS, including their name, ARN, replication strategy, and other details.

## Table Usage Guide

The `aws_keyspaces_keyspace` table enables cloud administrators and DevOps engineers to gather detailed insights into their Keyspaces keyspaces. You can query various aspects of the keyspace, such as its replication strategy, creation timestamp, and replication regions. This table is useful for monitoring keyspaces configurations and ensuring they align with organizational requirements.

## Examples

### Basic keyspace information
Retrieve basic information about your AWS Keyspaces keyspaces, including their name, ARN, and region.

```sql+postgres
select
  keyspace_name,
  arn,
  region
from
  aws_keyspaces_keyspace;
```

```sql+sqlite
select
  keyspace_name,
  arn,
  region
from
  aws_keyspaces_keyspace;
```

### List keyspaces by region
Retrieve a list of keyspaces grouped by their region. This is useful for understanding your keyspaces distribution across AWS regions.

```sql+postgres
select
  keyspace_name,
  arn,
  region
from
  aws_keyspaces_keyspace
order by
  region;
```

```sql+sqlite
select
  keyspace_name,
  arn,
  region
from
  aws_keyspaces_keyspace
order by
  region;
```

### Identify keyspaces with specific replication strategy
Find keyspaces using a specific replication strategy, which can help ensure compliance with replication policies.

```sql+postgres
select
  keyspace_name,
  arn,
  replication_strategy
from
  aws_keyspaces_keyspace
where
  replication_strategy = 'SINGLE_REGION';
```

```sql+sqlite
select
  keyspace_name,
  arn,
  replication_strategy
from
  aws_keyspaces_keyspace
where
  replication_strategy = 'SINGLE_REGION';
```

### List keyspaces with their replication regions
Retrieve a list of keyspaces along with their replication regions to understand where data is being replicated.

```sql+postgres
select
  keyspace_name,
  arn,
  replication_regions
from
  aws_keyspaces_keyspace;
```

```sql+sqlite
select
  keyspace_name,
  arn,
  replication_regions
from
  aws_keyspaces_keyspace;
```