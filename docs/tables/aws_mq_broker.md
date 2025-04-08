---
title: "Steampipe Table: aws_mq_broker - Query AWS MQ Brokers using SQL"
description: "Allows users to query AWS MQ Brokers."
folder: "Amazon MQ"
---

# Table: aws_mq_broker - Query AWS MQ Brokers using SQL

Amazon MQ is a managed message broker service provided by AWS (Amazon Web Services). It supports popular messaging protocols such as MQTT, AMQP, and STOMP, making it compatible with a variety of applications. Amazon MQ simplifies the setup, deployment, and maintenance of message brokers, allowing you to focus on developing your applications.

## Table Usage Guide

The `aws_mq_broker` table in Steampipe provides you with information about MQ brokers within AWS. This table allows you, as a DevOps engineer, to query broker specific details, including the boker ARN, creation time, and associated metadata. You can utilize this table to gather insights on nrokers, such as the number of broker nodes, the version and type of the engine used, the state of the broker, and more. The schema outlines the various attributes of the MQ broker for you, including the encryption info, authentication strategy, and associated tags.

## Examples

### Basic Info
Explore the status and details of your AWS MQ broker to understand their configuration and operational state.

```sql+postgres
select
  arn,
  broker_name,
  broker_state,
  deployment_mode,
  created,
  host_instance_type,
  engine_type,
  engine_version
  tags
from
  aws_mq_broker;
```

```sql+sqlite
select
  arn,
  broker_name,
  broker_state,
  deployment_mode,
  created,
  host_instance_type,
  engine_type,
  engine_version
  tags
from
  aws_mq_broker;
```

### List brokers that are in rebooting state
Identify certain brokers within AWS MQ service that are in reboot state. This could be useful for system administrators who need to manage resources.

```sql+postgres
select
  arn,
  broker_name,
  broker_state,
  created
  data_replication_mode,
  authentication_strategy
from
  aws_mq_broker
where
  broker_state = 'REBOOT_IN_PROGRESS';
```

```sql+sqlite
select
  arn,
  broker_name,
  broker_state,
  created
  data_replication_mode,
  authentication_strategy
from
  aws_mq_broker
where
  broker_state = 'REBOOT_IN_PROGRESS';
```

### List brokers that allow public access
Determine the areas in which public access is allowed for broker. This is useful for identifying potential security risks and ensuring that access to sensitive data is appropriately restricted.

```sql+postgres
select
  arn,
  broker_name,
  broker_state,
  created
from
  aws_mq_broker
where
  publicly_accessible;
```

```sql+sqlite
select
  arn,
  broker_name,
  broker_state,
  created
from
  aws_mq_broker
where
  publicly_accessible;
```

### List brokers that encrypted with customer managed key
Identify the specific domains or components within the system where data is secured through encryption using keys managed by the customer.

```sql+postgres
select
  arn,
  broker_name,
  encryption_options ->> 'UseAwsOwnedKey' as use_aws_owned_key,
  created
from
  aws_mq_broker
where
  encryption_options ->> 'UseAwsOwnedKey' = 'false';
```

```sql+sqlite
select
  arn,
  broker_name,
  json_extract(encryption_options, '$.EncryptionInfo.EncryptionAtRest') as use_aws_owned_key,
  created
from
  aws_mq_broker
where
  json_extract(encryption_options, '$.UseAwsOwnedKey') = 'false';
```

### Get maintenance window details of brokers
During the Maintenance Window, the broker instances might be briefly unavailable or experience reduced capacity as updates are applied. This scheduled approach helps minimize the impact on your applications and users, as these activities are carried out during a designated time frame, allowing for predictability and coordination.

```sql+postgres
select
  arn,
  broker_name,
  maintenance_window_start_time -> 'DayOfWeek' as day_of_week,
  maintenance_window_start_time -> 'TimeOfDay' as time_of_day,
  maintenance_window_start_time -> 'TimeZone' as time_zone
from
  aws_mq_broker;
```

```sql+sqlite
select
  arn,
  broker_name,
  json_extract(maintenance_window_start_time, '$.DayOfWeek') as day_of_week,
  json_extract(maintenance_window_start_time, '$.TimeOfDay') as time_of_day,
  json_extract(maintenance_window_start_time, '$.TimeZone') as time_zone
from
  aws_mq_broker;
```