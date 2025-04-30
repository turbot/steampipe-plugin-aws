---
title: "Steampipe Table: aws_sagemaker_device_fleet - Query AWS SageMaker Device Fleet using SQL"
description: "Allows users to query AWS SageMaker Device Fleet for detailed information about device fleets, including their configuration, status, and associated metadata."
folder: "SageMaker"
---

# Table: aws_sagemaker_device_fleet - Query AWS SageMaker Device Fleet using SQL

AWS SageMaker Device Fleet is a service that enables you to manage and monitor fleets of edge devices for machine learning inference. A device fleet is a collection of edge devices that run SageMaker Edge Manager to manage ML models and collect data. This service helps you deploy, manage, and monitor machine learning models on edge devices at scale.

## Table Usage Guide

The `aws_sagemaker_device_fleet` table in Steampipe provides you with information about device fleets within AWS SageMaker. This table allows you, as a DevOps engineer, data scientist, or machine learning engineer, to query device fleet-specific details, including fleet configuration, IoT role settings, output configurations, and associated metadata. You can utilize this table to gather insights on device fleets, such as their creation time, last modification time, description, and associated tags. The schema outlines the various attributes of the device fleet for you, including the fleet name, ARN, creation time, and associated tags.

## Examples

### Basic info
Explore the basic information of your AWS SageMaker device fleets to understand their configuration and status. This can help in managing and monitoring your edge device deployments effectively.

```sql+postgres
select
  device_fleet_name,
  arn,
  creation_time,
  last_modified_time,
  description
from
  aws_sagemaker_device_fleet;
```

```sql+sqlite
select
  device_fleet_name,
  arn,
  creation_time,
  last_modified_time,
  description
from
  aws_sagemaker_device_fleet;
```

### List device fleets with IoT role configuration
Identify device fleets and their associated IoT role configurations to ensure proper access control and security settings are in place for your edge devices.

```sql+postgres
select
  device_fleet_name,
  iot_role_alias,
  role_arn
from
  aws_sagemaker_device_fleet;
```

```sql+sqlite
select
  device_fleet_name,
  iot_role_alias,
  role_arn
from
  aws_sagemaker_device_fleet;
```

### Get output configuration details for device fleets
Review the output configuration settings for your device fleets to understand how data is being collected and stored from your edge devices.

```sql+postgres
select
  device_fleet_name,
  output_config
from
  aws_sagemaker_device_fleet;
```

```sql+sqlite
select
  device_fleet_name,
  output_config
from
  aws_sagemaker_device_fleet;
```

### List device fleets with tags
Explore device fleets along with their associated tags to better organize and manage your edge device deployments.

```sql+postgres
select
  device_fleet_name,
  tags
from
  aws_sagemaker_device_fleet;
```

```sql+sqlite
select
  device_fleet_name,
  tags
from
  aws_sagemaker_device_fleet;
```

### Find device fleets created in the last 30 days
Identify recently created device fleets to track new deployments and ensure they are properly configured and monitored.

```sql+postgres
select
  device_fleet_name,
  creation_time
from
  aws_sagemaker_device_fleet
where
  creation_time >= now() - interval '30 days';
```

```sql+sqlite
select
  device_fleet_name,
  creation_time
from
  aws_sagemaker_device_fleet
where
  creation_time >= datetime('now', '-30 days');
``` 