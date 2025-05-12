---
title: "Steampipe Table: aws_iot_thing - Query AWS Internet of Things using SQL"
description: "Allows users to query AWS Internet of Things to retrieve detailed information about the the virtual model of a physical device with in an AWS account."
folder: "IoT Core"
---

# Table: aws_iot_thing - Query AWS IoT Things using SQL

AWS IoT Core allows for the management and connection of Internet of Things (IoT) devices to the AWS cloud. In this service, an IoT "Thing" represents the virtual model of a physical device or asset. These things, which can range from sensors and appliances to machinery, are capable of collecting and exchanging data over the internet or other networks.

## Table Usage Guide

The `aws_iot_thing` table facilitates the exploration and management of AWS IoT Things. Utilize this table to access comprehensive details about each IoT Thing, including its name, ID, type, and attributes. This is particularly beneficial for IoT administrators and developers who need to oversee IoT device configurations and statuses within the AWS ecosystem.

## Examples

### Basic Info
Retrieve fundamental details about your AWS IoT Things. This query is crucial for a general overview of IoT Things in your environment, including their identifiers and types.

```sql+postgres
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing;
```

```sql+sqlite
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing;
```

### Filter Things by Attribute Name
Identify specific IoT Things based on a particular attribute name. This query is useful for segmenting IoT Things according to custom-defined criteria or characteristics.

```sql+postgres
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing
where
  attribute_name = 'foo';
```

```sql+sqlite
select
  thing_name,
  thing_id,
  arn,
  thing_type_name,
  version
from
  aws_iot_thing
where
  attribute_name = 'foo';
```

### List Things for a Given Type Name
Find all IoT Things that belong to a specific type. This query aids in understanding the distribution and categorization of IoT Things by type within your AWS IoT environment.

```sql+postgres
select
  thing_name,
  arn,
  thing_id,
  thing_type_name,
  attribute_value
from
  aws_iot_thing
where
  thing_type_name = 'foo';
```

```sql+sqlite
select
  thing_name,
  arn,
  thing_id,
  thing_type_name,
  attribute_value
from
  aws_iot_thing
where
  thing_type_name = 'foo';
```
