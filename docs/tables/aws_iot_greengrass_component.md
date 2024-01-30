---
title: "Steampipe Table: aws_iot_greengrass_component - Query AWS IoT Greengrass Components using SQL"
description: "Allows users to query AWS IoT Greengrass Components. This table provides information about Greengrass Components, which are software modules deployed to Greengrass core devices. Components can represent applications, runtime installers, libraries, or any code that runs on a device. Users can define components with dependencies, and Greengrass deploys only the required software modules when deploying to devices."
---

# Table: aws_iot_greengrass_component - Query AWS IoT Greengrass Components using SQL

AWS IoT Greengrass components are software modules that you deploy to Greengrass core devices. Components can represent applications, runtime installers, libraries, or any code that you would run on a device. You can define components that depend on other components. For example, you might define a component that installs Python, and then define that component as a dependency of your components that run Python applications. When you deploy your components to your fleets of devices, Greengrass deploys only the software modules that your devices require.

## Table Usage Guide

The `aws_iot_greengrass_component` table in Steampipe enables users to query information about AWS IoT Greengrass Components. This includes details such as component name, ARN, recipe, and recipe output format. Users can also retrieve the latest version details of components, filter components by tags, and gain insights into component dependencies.

## Examples

### Basic info
Retrieve basic information about Greengrass Components, such as their names, ARNs, recipes, and recipe output formats. This query provides an overview of the components in your environment.

```sql+postgres
select
  component_name,
  arn,
  recipe,
  recipe_output_format
from
  aws_iot_greengrass_component;
```

```sql+sqlite
select
  component_name,
  arn,
  recipe,
  recipe_output_format
from
  aws_iot_greengrass_component;
```

### Get the latest version details of components
Retrieve the latest version details of Greengrass Components, including the ARN, component version, description, creation timestamp, supported platforms, and publisher. This query helps you stay up-to-date with the latest component versions.

```sql+postgres
select
  component_name,
  latest_version ->> 'Arn' as latest_version_arn,
  latest_version ->> 'ComponentVersion' as component_version,
  latest_version ->> 'Description' as description,
  latest_version ->> 'CreationTimestamp' as creation_timestamp,
  latest_version -> 'Platforms' as platforms,
  latest_version ->> 'Publisher' as publisher
from
  aws_iot_greengrass_component;
```

```sql+sqlite
select
  component_name,
  json_extract(latest_version, '$.Arn') as latest_version_arn,
  json_extract(latest_version, '$.ComponentVersion') as component_version,
  json_extract(latest_version, '$.Description') as description,
  json_extract(latest_version, '$.CreationTimestamp') as creation_timestamp,
  json_extract(latest_version, '$.Platforms') as platforms,
  json_extract(latest_version, '$.Publisher') as publisher
from
  aws_iot_greengrass_component;
```

### Filter the components by tag
Filter Greengrass Components based on their tags. In this example, components with a missing 'owner' tag are retrieved. This query allows you to categorize and manage components effectively.

```sql+postgres
select
  component_name,
  arn,
  recipe
from
  aws_iot_greengrass_component
where
  tags -> 'owner' is null;
```

```sql+sqlite
select
  component_name,
  arn,
  recipe
from
  aws_iot_greengrass_component
where
  json_extract(tags, '$.owner') is null;
```
