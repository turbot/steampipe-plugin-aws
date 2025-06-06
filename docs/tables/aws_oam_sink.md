---
title: "Steampipe Table: aws_oam_sink - Query AWS OAM Sink using SQL"
description: "Allows users to query AWS OAM Sink data, providing detailed information about each AWS OAM Sink in your AWS account."
folder: "OAM"
---

# Table: aws_oam_sink - Query AWS OAM Sink using SQL

The AWS OAM (Operational Application Management) Sink is a component of the AWS Distro for OpenTelemetry (ADOT) that collects, processes, and exports telemetry data. It enables you to understand the performance of your applications and services, troubleshoot operational issues, and gain insights into customer behavior. As a managed service, OAM Sink simplifies the process of gathering and analyzing high-volume, high-velocity data generated by your AWS resources and applications.

## Table Usage Guide

The `aws_oam_sink` table in Steampipe provides you with information about each OAM Sink within AWS Operational Application Manager (OAM). This table allows you as a DevOps engineer to query Sink-specific details, including the Sink ARN, creation date, Sink status, and associated metadata. You can utilize this table to gather insights on Sinks, such as Sinks with specific statuses, verification of Sink properties, and more. The schema outlines the various attributes of the OAM Sink for you, including the Sink ARN, creation date, status, and associated tags.

## Examples

### Basic info
Determine the areas in which specific resources in your AWS environment are tagged or titled, allowing for efficient organization and management of resources.

```sql+postgres
select
  name,
  id,
  arn,
  tags,
  title
from
  aws_oam_sink;
```

```sql+sqlite
select
  name,
  id,
  arn,
  tags,
  title
from
  aws_oam_sink;
```

### Get sink by ID
Determine the specific details of a unique resource within your AWS environment. This is particularly useful when you need to quickly access information about a single resource based on its unique identifier.

```sql+postgres
select
  name,
  id,
  arn
from
  aws_oam_sink
where
  id = 'hfj44c81-7bdf-3847-r7i3-5dfc61b17483';
```

```sql+sqlite
select
  name,
  id,
  arn
from
  aws_oam_sink
where
  id = 'hfj44c81-7bdf-3847-r7i3-5dfc61b17483';
```