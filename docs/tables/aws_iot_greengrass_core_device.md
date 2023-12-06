# Table: aws_iot_greengrass_core_device

A Greengrass core is a device that runs the AWS IoT Greengrass Core software, which allows it to communicate directly with AWS IoT Core and the AWS IoT Greengrass service. A core has its own device certificate used for authenticating with AWS IoT Core. It has a device shadow and an entry in the AWS IoT Core registry. Greengrass cores run a local Lambda runtime, deployment agent, and IP address tracker that sends IP address information to the AWS IoT Greengrass service to allow client devices to automatically discover their group and core connection information.

## Examples

### Basic info

```sql
select
  core_device_thing_name,
  architecture,
  core_version,
  platform
from
  aws_iot_greengrass_core_device;
```

### List devices that are unhealthy

```sql
select
  core_device_thing_name,
  architecture,
  core_version
  last_status_update_timestamp,
  status
from
  aws_iot_greengrass_core_device
where
  status = 'UNHEALTHY';
```

### List devices by platform

```sql
select
  core_device_thing_name,
  architecture,
  core_version
  platform
from
  aws_iot_greengrass_core_device
where
  platform ILIKE 'linux';
```
