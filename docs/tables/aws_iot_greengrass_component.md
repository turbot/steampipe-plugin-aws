# Table: aws_iot_greengrass_component

AWS IoT Greengrass components are software modules that you deploy to Greengrass core devices. Components can represent applications, runtime installers, libraries, or any code that you would run on a device. You can define components that depend on other components. For example, you might define a component that installs Python, and then define that component as a dependency of your components that run Python applications. When you deploy your components to your fleets of devices, Greengrass deploys only the software modules that your devices require.

## Examples

### Basic info

```sql
select
  component_name,
  arn,
  recipe,
  recipe_output_format
from
  aws_iot_greengrass_component;
```

### Get the latest version details of components

```sql
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

### Filter the components by tag

```sql
select
  component_name,
  arn,
  recipe
from
  aws_iot_greengrass_component
where
  tags -> 'owner' is null;
```
