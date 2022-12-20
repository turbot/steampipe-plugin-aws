# Table: aws_eks_addon_version

Amazon EKS add-ons help to automate the provisioning and lifecycle management of common operational software for Amazon EKS clusters.

Amazon EKS add-on versions describe the Kubernetes versions that the add-on can be used with.

## Examples

### Basic info

```sql
select
  addon_name,
  addon_version,
  type
from
  aws_eks_addon_version;
```

### Count the number of add-on versions by add-on

```sql
select
  addon_name,
  count(addon_version) as addon_version_count
from
  aws_eks_addon_version
group by
  addon_name;
```

### Get configuration details of each add-on version

```sql
select
  addon_name,
  addon_version,
  addon_configuration -> '$defs' -> 'extraVolumeTags' ->> 'description' as configuration_def_description,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'propertyNames' as configuration_def_property_names,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'propertyNames' as configuration_def_pattern_properties,
  addon_configuration -> 'properties' as configuration_properties
from
  aws_eks_addon_version limit 10;
```