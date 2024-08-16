---
title: "Steampipe Table: aws_qldb_stream - Query AWS QLDB Streams using SQL"
description: "Allows users to query AWS QLDB streams, providing detailed information on stream configurations, statuses, and associated Kinesis Data Streams."
---

# Table: aws_qldb_stream - Query AWS QLDB Streams using SQL

Amazon QLDB (Quantum Ledger Database) is a fully managed ledger database that provides a transparent, immutable, and cryptographically verifiable transaction log. QLDB Streams allow you to capture and process changes to your QLDB ledger's journal in near real-time by streaming them to a Kinesis Data Streams destination. The `aws_qldb_stream` table in Steampipe allows you to query information about these QLDB streams in your AWS environment. This includes details such as stream status, creation time, associated ledger, and Kinesis configuration.

## Table Usage Guide

The `aws_qldb_stream` table enables cloud administrators and DevOps engineers to gather detailed insights into their QLDB streams. You can query various aspects of the stream, such as its status, associated ledger, Kinesis configuration, and error causes. This table is particularly useful for monitoring stream health, ensuring data integrity, and managing stream configurations.

## Examples

### Basic stream information
Retrieve basic information about your AWS QLDB streams, including their name, ARN, status, and region. This can be useful for getting an overview of the streams deployed in your AWS account.

```sql+postgres
select
  stream_name,
  arn,
  status,
  creation_time,
  region
from
  aws_qldb_stream;
```

```sql+sqlite
select
  stream_name,
  arn,
  status,
  creation_time,
  region
from
  aws_qldb_stream;
```

### List active streams
Fetch a list of streams that are currently active. This can help in identifying which streams are operational and available for data streaming.

```sql+postgres
select
  stream_name,
  arn,
  status
from
  aws_qldb_stream
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  stream_name,
  arn,
  status
from
  aws_qldb_stream
where
  status = 'ACTIVE';
```

### List streams with error causes
Identify streams that have encountered errors, which can help in diagnosing issues related to data streaming or permissions.

```sql+postgres
select
  stream_name,
  arn,
  status,
  error_cause
from
  aws_qldb_stream
where
  error_cause is not null;
```

```sql+sqlite
select
  stream_name,
  arn,
  status,
  error_cause
from
  aws_qldb_stream
where
  error_cause is not null;
```

### List streams by ledger
Retrieve streams that are associated with a specific QLDB ledger, which can be useful for auditing and managing data flows from particular ledgers.

```sql+postgres
select
  stream_name,
  arn,
  ledger_name,
  status
from
  aws_qldb_stream
where
  ledger_name = 'your-ledger-name';
```

```sql+sqlite
select
  stream_name,
  arn,
  ledger_name,
  status
from
  aws_qldb_stream
where
  ledger_name = 'your-ledger-name';
```

### List streams by creation date
Fetch streams ordered by their creation date, which can be useful for auditing purposes or understanding the lifecycle of your QLDB streams.

```sql+postgres
select
  stream_name,
  arn,
  creation_time
from
  aws_qldb_stream
order by
  creation_time desc;
```

```sql+sqlite
select
  stream_name,
  arn,
  creation_time
from
  aws_qldb_stream
order by
  creation_time desc;
```

### Get stream Kinesis configuration details
Retrieve detailed information about the Kinesis Data Streams configuration associated with your QLDB streams, which is essential for understanding how data is being streamed and processed.

```sql+postgres
select
  stream_name,
  arn,
  kinesis_configuration
from
  aws_qldb_stream;
```

```sql+sqlite
select
  stream_name,
  arn,
  kinesis_configuration
from
  aws_qldb_stream;
```