---
title: "Steampipe Table: aws_kinesis_stream - Query AWS Kinesis Stream using SQL"
description: "Allows users to query AWS Kinesis Stream data, including stream name, status, creation time, and associated tags."
folder: "Kinesis"
---

# Table: aws_kinesis_stream - Query AWS Kinesis Stream using SQL

The AWS Kinesis Stream is a resource in Amazon Kinesis Data Streams that allows you to build custom applications that process or analyze streaming data for specialized needs. It can continuously capture and store terabytes of data per hour from hundreds of thousands of sources. This real-time data stream processing makes it easy to analyze and process data as it arrives.

## Table Usage Guide

The `aws_kinesis_stream` table in Steampipe provides you with information about Kinesis streams within AWS Kinesis. This table allows you, as a DevOps engineer, to query stream-specific details, including the stream name, status, creation time, and associated metadata. You can utilize this table to gather insights on streams, such as stream health, data throughput, and more. The schema outlines the various attributes of the Kinesis stream for you, including the stream ARN, creation timestamp, number of shards, and associated tags.

## Examples

### Basic info
Explore which AWS Kinesis streams are active and how many consumers each has, to better manage resource allocation and optimize data flow. This information can also provide insights into stream usage patterns over time and across different regions.

```sql+postgres
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream;
```

```sql+sqlite
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream;
```


### List streams that are not active
Determine the areas in which streams are inactive to manage resources better and optimize performance. This can be particularly useful when auditing system activity or troubleshooting issues.

```sql+postgres
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  stream_status != 'ACTIVE';
```

```sql+sqlite
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  stream_status <> 'ACTIVE';
```


### List streams that have no consumers
Explore which data streams are currently not being used by any consumers in your AWS Kinesis setup. This can help identify unused resources for potential clean up or reallocation.

```sql+postgres
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  consumer_count = 0;
```

```sql+sqlite
select
  stream_name,
  stream_arn,
  stream_status,
  consumer_count,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  consumer_count = 0;
```


### List streams that are not encrypted
Discover the segments that are transmitting data without any encryption, which could potentially expose sensitive information and pose a security risk. This query is useful for identifying these unprotected streams and improving your data security measures.

```sql+postgres
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type = 'NONE';
```

```sql+sqlite
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type = 'NONE';
```


### List streams that are not encrypted using CMK
Discover the segments that are not secured using Customer Master Key (CMK) in your Kinesis streams. This is useful for ensuring all your data streams are adequately protected, maintaining your data's privacy and security.

```sql+postgres
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type != 'NONE'
  and key_id = 'alias/aws/kinesis';
```

```sql+sqlite
select
  stream_name,
  stream_arn,
  encryption_type,
  key_id,
  stream_creation_timestamp,
  region
from
  aws_kinesis_stream
where
  encryption_type != 'NONE'
  and key_id = 'alias/aws/kinesis';
```