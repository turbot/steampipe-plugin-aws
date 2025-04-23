---
title: "Steampipe Table: aws_iot_thing_type - Query AWS IoT Thing Type using SQL"
description: "Allows users to query AWS IoT Thing Type to gain insights into each thing type's configuration, including ARN, name, creation date, and deprecation status."
folder: "IoT Core"
---

# Table: aws_iot_thing_type - Query AWS IoT Thing Types using SQL

In AWS IoT Core, the IoT Thing Type feature enables the categorization of IoT devices (Things) into different types based on shared characteristics or use cases. A Thing Type acts as a template, defining attributes and properties common to a specific class or category of devices. Defining Thing Types facilitates more effective and consistent management and interaction with groups of Things.

## Table Usage Guide

The `aws_iot_thing_type` table can be used to access detailed information about the types of IoT Things, including their names, IDs, descriptions, and creation dates. This table is particularly useful for IoT administrators and developers who need to oversee the classification and properties of IoT devices within AWS.

## Examples

### Basic info
Retrieve essential details about AWS IoT Thing Types. This query is vital for understanding the various types of IoT Things, their descriptions, and when they were created.

```sql+postgres
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date
from
  aws_iot_thing_type;
```

```sql+sqlite
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date
from
  aws_iot_thing_type;
```

### List deprecated thing types
Identify Thing Types that have been marked as deprecated. This query helps track Thing Types that are no longer recommended for use.

```sql+postgres
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated
from
  aws_iot_thing_type
where
  deprecated;
```

```sql+sqlite
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated
from
  aws_iot_thing_type
where
  deprecated;
```

### List thing types created in the last 30 days
Find Thing Types that were created within the last 30 days. This query is useful for monitoring recent additions and updates in your IoT environment.

```sql+postgres
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated,
  searchable_attributes
from
  aws_iot_thing_type
where
  creation_date >= now() - interval '30 days';
```

```sql+sqlite
select
  thing_type_name,
  arn,
  thing_type_id,
  thing_type_description,
  creation_date,
  deprecated,
  searchable_attributes
from
  aws_iot_thing_type
where
  datetime(creation_date) >= datetime('now', '-30 days');
```

### List thing types scheduled for deprecation within the next 30 days
Discover Thing Types scheduled for deprecation in the next 30 days. This query assists in proactive planning for transitioning away from soon-to-be deprecated Thing Types.

```sql+postgres
select
  thing_type_name,
  arn,
  thing_type_id,
  creation_date,
  tags,
  deprecation_date
from
  aws_iot_thing_type
where
  deprecation_date <= now() - interval '30 days';
```

```sql+sqlite
select
  thing_type_name,
  arn,
  thing_type_id,
  creation_date,
  tags,
  deprecation_date
from
  aws_iot_thing_type
where
  datetime(deprecation_date) <= datetime('now', '-30 days');
```
