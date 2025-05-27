---
title: "Steampipe Table: aws_eks_addon_version - Query AWS EKS Add-On Versions using SQL"
description: "Allows users to query AWS EKS Add-On Versions."
folder: "EKS"
---

# Table: aws_eks_addon_version - Query AWS EKS Add-On Versions using SQL

The AWS EKS Add-On Versions are a part of the Amazon Elastic Kubernetes Service (EKS), which is a managed service that makes it easy for you to run Kubernetes on AWS without needing to install and operate your own Kubernetes control plane. Add-Ons help to automate the process of installing, upgrading, and operating additional Kubernetes software. They are versions of Kubernetes software components that can be installed onto your Amazon EKS clusters.

## Table Usage Guide

The `aws_eks_addon_version` table in Steampipe provides you with information about Add-On versions within Amazon Elastic Kubernetes Service (EKS). This table allows you, as a DevOps engineer, to query add-on specific details, including addon name, addon version, architecture, and associated metadata. You can utilize this table to gather insights on add-ons, such as the add-on version status, the specific architectures it supports, and more. The schema outlines the various attributes of the EKS add-on for you, including the add-on name, add-on version, and supported architectures.

## Examples

### Basic info
Explore which addons are available for your AWS EKS service and identify their versions to ensure compatibility and optimal performance. This can be beneficial in maintaining an updated and efficient system.

```sql+postgres
select
  addon_name,
  addon_version,
  type
from
  aws_eks_addon_version;
```

```sql+sqlite
select
  addon_name,
  addon_version,
  type
from
  aws_eks_addon_version;
```

### Count the number of add-on versions by add-on
Determine the areas in which various versions of add-ons are being used within your AWS Elastic Kubernetes Service (EKS). This can help in understanding the distribution and usage of different add-on versions, aiding in effective management and potential upgrades.

```sql+postgres
select
  addon_name,
  count(addon_version) as addon_version_count
from
  aws_eks_addon_version
group by
  addon_name;
```

```sql+sqlite
select
  addon_name,
  count(addon_version) as addon_version_count
from
  aws_eks_addon_version
group by
  addon_name;
```

### Get configuration details of each add-on version
Explore the specific configuration details for each version of an add-on to understand how it's set up and functions. This can help to identify any potential issues or areas for improvement in your AWS EKS environment.

```sql+postgres
select
  addon_name,
  addon_version,
  addon_configuration -> '$defs' -> 'extraVolumeTags' ->> 'description' as addon_configuration_def_description,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'propertyNames' as addon_configuration_def_property_names,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'patternProperties' as addon_configuration_def_pattern_properties,
  addon_configuration -> 'properties' as addon_configuration_properties
from
  aws_eks_addon_version limit 10;
```

```sql+sqlite
select
  addon_name,
  addon_version,
  json_extract(addon_configuration, '$.$defs.extraVolumeTags.description') as addon_configuration_def_description,
  json_extract(addon_configuration, '$.$defs.extraVolumeTags.propertyNames') as addon_configuration_def_property_names,
  json_extract(addon_configuration, '$.$defs.extraVolumeTags.patternProperties') as addon_configuration_def_pattern_properties,
  json_extract(addon_configuration, '$.properties') as addon_configuration_properties
from
  aws_eks_addon_version limit 10;
```