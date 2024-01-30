---
title: "Steampipe Table: aws_iot_greengrass_core_device - Query AWS IoT Greengrass Core Devices using SQL"
description: "Allows users to query AWS IoT Greengrass Core Devices. This table provides information about Greengrass Core Devices within AWS IoT Greengrass, enabling users to gather insights on core devices, including their thing name, architecture, core version, and platform."
---

# Table: aws_iot_greengrass_core_device - Query AWS IoT Greengrass Core Devices using SQL

A Greengrass core device is a fundamental component that runs the AWS IoT Greengrass Core software. This software allows the core device to establish direct communication with AWS IoT Core and the AWS IoT Greengrass service. Each core device possesses its unique device certificate for authentication with AWS IoT Core. Additionally, core devices have a device shadow and are registered in the AWS IoT Core registry. Greengrass core devices operate a local Lambda runtime, deployment agent, and IP address tracker. The IP address tracker sends IP address information to the AWS IoT Greengrass service, enabling client devices to automatically discover their group and core connection information.

## Table Usage Guide

The `aws_iot_greengrass_core_device` table in Steampipe provides IoT engineers and DevOps professionals with the ability to query specific details about AWS IoT Greengrass Core Devices. This includes essential information like core device thing name, architecture, core version, and platform. You can use this table to monitor the health and status of Greengrass core devices, identify unhealthy devices, and filter devices by their platform for efficient management.

## Examples

### Basic info
Retrieve basic information about Greengrass Core Devices, such as their thing names, architecture, core versions, and platforms. This query provides an overview of the core devices in your environment.

```sql+postgres
select
  core_device_thing_name,
  architecture,
  core_version,
  platform
from
  aws_iot_greengrass_core_device;
```

```sql+sqlite
select
  core_device_thing_name,
  architecture,
  core_version,
  platform
from
  aws_iot_greengrass_core_device;
```

### List devices that are unhealthy
Identify Greengrass Core Devices that are marked as 'UNHEALTHY' in their status. This helps you monitor and manage devices with potential issues.

```sql+postgres
select
  core_device_thing_name,
  architecture,
  core_version,
  last_status_update_timestamp,
  status
from
  aws_iot_greengrass_core_device
where
  status = 'UNHEALTHY';
```

```sql+sqlite
select
  core_device_thing_name,
  architecture,
  core_version,
  last_status_update_timestamp,
  status
from
  aws_iot_greengrass_core_device
where
  status = 'UNHEALTHY';
```

### List devices by platform
Filter Greengrass Core Devices by platform, focusing on devices with a platform that matches 'linux.' This query allows you to group and manage devices based on their platform type.

```sql+postgres
select
  core_device_thing_name,
  architecture,
  core_version,
  platform
from
  aws_iot_greengrass_core_device
where
  platform ILIKE 'linux';
```

```sql+sqlite
select
  core_device_thing_name,
  architecture,
  core_version,
  platform
from
  aws_iot_greengrass_core_device
where
  platform LIKE 'linux';
```
