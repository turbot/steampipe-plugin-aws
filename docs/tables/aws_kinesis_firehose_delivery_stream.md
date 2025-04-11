---
title: "Steampipe Table: aws_kinesis_firehose_delivery_stream - Query AWS Kinesis Firehose Delivery Stream using SQL"
description: "Allows users to query AWS Kinesis Firehose Delivery Stream data, providing detailed information about each delivery stream in the AWS account."
folder: "Kinesis"
---

# Table: aws_kinesis_firehose_delivery_stream - Query AWS Kinesis Firehose Delivery Stream using SQL

The AWS Kinesis Firehose Delivery Stream is a fully managed service that makes it easy to capture, transform, and load data streams into AWS data stores for near real-time analytics with existing business intelligence tools. It can capture, transform, and deliver streaming data to Amazon S3, Amazon Redshift, Amazon Elasticsearch Service, and Splunk, enabling near real-time analytics with existing business intelligence tools and dashboards. It is a key data streaming solution as part of the AWS ecosystem.

## Table Usage Guide

The `aws_kinesis_firehose_delivery_stream` table in Steampipe provides you with information about each Kinesis Firehose Delivery Stream within AWS. This table allows you, as a DevOps engineer, to query delivery stream-specific details, including the delivery stream name, status, creation time, and associated metadata. You can utilize this table to gather insights on delivery streams, such as the status of the stream, its creation time, and the destinations configured for the stream. The schema outlines the various attributes of the Kinesis Firehose Delivery Stream, including the delivery stream ARN, delivery stream type, delivery stream status, and more for you.

## Examples

### Basic info
Explore which AWS Kinesis Firehose streams are active and when they were created to better manage and monitor data flow. This can be useful for auditing purposes or to optimize resource allocation.

```sql+postgres
select
  delivery_stream_name,
  arn,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream;
```

```sql+sqlite
select
  delivery_stream_name,
  arn,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream;
```


### List inactive delivery streams
Determine the areas in which delivery streams are inactive within the AWS Kinesis Firehose service. This is useful for identifying potential issues or bottlenecks in your data flow.

```sql+postgres
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_status != 'ACTIVE';
```

```sql+sqlite
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_status != 'ACTIVE';
```


### List delivery streams that are not encrypted
Identify instances where your Kinesis Firehose delivery streams are not encrypted. This is crucial for ensuring the security of your data in transit.

```sql+postgres
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type,
  delivery_stream_encryption_configuration ->> 'Status' as encryption_status
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_encryption_configuration ->> 'Status' = 'DISABLED';
```

```sql+sqlite
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type,
  json_extract(delivery_stream_encryption_configuration, '$.Status') as encryption_status
from
  aws_kinesis_firehose_delivery_stream
where
  json_extract(delivery_stream_encryption_configuration, '$.Status') = 'DISABLED';
```


### List delivery streams for a specific delivery stream type
Determine the areas in which specific types of delivery streams are being used within the AWS Kinesis Firehose service. This can help in managing and optimizing the use of different delivery stream types.

```sql+postgres
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_type = 'DirectPut';
```

```sql+sqlite
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  create_timestamp,
  delivery_stream_type
from
  aws_kinesis_firehose_delivery_stream
where
  delivery_stream_type = 'DirectPut';
```


### List delivery streams with at least one failure
Identify instances where delivery streams have experienced at least one failure. This is useful for diagnosing issues and ensuring reliable data delivery within AWS Kinesis.

```sql+postgres
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  delivery_stream_type,
  failure_description
from
  aws_kinesis_firehose_delivery_stream
where
  failure_description is not null;
```

```sql+sqlite
select
  delivery_stream_name,
  arn,
  delivery_stream_status,
  delivery_stream_type,
  failure_description
from
  aws_kinesis_firehose_delivery_stream
where
  failure_description is not null;
```