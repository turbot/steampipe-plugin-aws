---
title: "Table: aws_oam_sink - Query AWS OAM Sink using SQL"
description: "Allows users to query AWS OAM Sink data, providing detailed information about each AWS OAM Sink in your AWS account."
---

# Table: aws_oam_sink - Query AWS OAM Sink using SQL

The `aws_oam_sink` table in Steampipe provides information about each OAM Sink within AWS Operational Application Manager (OAM). This table allows DevOps engineers to query Sink-specific details, including the Sink ARN, creation date, Sink status, and associated metadata. Users can utilize this table to gather insights on Sinks, such as Sinks with specific statuses, verification of Sink properties, and more. The schema outlines the various attributes of the OAM Sink, including the Sink ARN, creation date, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_oam_sink` table, you can use the `.inspect aws_oam_sink` command in Steampipe.

### Key columns:

- `sink_arn`: The Amazon Resource Name (ARN) of the Sink. This can be used to join this table with other tables as it uniquely identifies the Sink.
- `sink_name`: The name of the Sink. This is an important identifier and can be used for filtering specific Sinks.
- `status`: The status of the Sink. This is useful for monitoring and managing the state of the Sink.

## Examples

### Basic info

```sql
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

```sql
select
  name,
  id,
  arn
from
  aws_oam_sink
where
  id = 'hfj44c81-7bdf-3847-r7i3-5dfc61b17483';
```