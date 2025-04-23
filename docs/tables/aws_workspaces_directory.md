---
title: "Steampipe Table: aws_workspaces_directory - Query AWS WorkSpaces Directory using SQL"
description: "Allows users to query AWS WorkSpaces Directory information to retrieve details such as directory ID, directory type, workspace creation properties, and more."
folder: "WorkSpaces"
---

# Table: aws_workspaces_directory - Query AWS WorkSpaces Directory using SQL

The AWS WorkSpaces Directory is a feature of Amazon WorkSpaces, a managed, secure Desktop-as-a-Service (DaaS) solution. It allows you to manage your WorkSpaces, provide a directory of users, and integrate with your corporate directory. It also enables you to control access, manage security settings, and monitor your WorkSpaces.

## Table Usage Guide

The `aws_workspaces_directory` table in Steampipe provides you with information about directories within AWS WorkSpaces. This table allows you, as a DevOps engineer, to query directory-specific details, including directory ID, directory type, workspace creation properties, workspace security group ID, and other associated metadata. You can utilize this table to gather insights on directories, such as the state of the directory, registration code, subnet IDs, and more. The schema outlines the various attributes of the AWS WorkSpaces Directory for you, including the self-service permissions, workspace access properties, and associated tags.

## Examples

## Basic info

```sql+postgres
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory;
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory;
```

### List directories that have certificate authority ARN enabled
Determine the areas in which certificate authority ARN is enabled in your AWS Workspaces directories. This helps in identifying directories that have enhanced security measures in place.

```sql+postgres
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  certificate_based_auth_properties ->> 'Status' = 'ENABLED';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  json_extract(certificate_based_auth_properties, '$.Status') = 'ENABLED';
```

### List directories of a particular type
Identify instances where AWS Workspaces directories are of a 'SIMPLE_AD' type. This helps users understand the distribution of directory types and manage their AWS Workspaces more effectively.

```sql+postgres
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  directory_type = 'SIMPLE_AD';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  directory_type = 'SIMPLE_AD';
```

### Get the SAML properties of a particular directory
This query allows you to examine the SAML properties associated with a specific AWS WorkSpaces directory. It's particularly useful when you need to assess the security and access configurations of your virtual desktop infrastructure.

```sql+postgres
select
  name,
  directory_id,
  arn,
  saml_properties ->> 'RelayStateParameterName' as saml_relay_state_parameter_name,
  saml_properties ->> 'Status' as saml_status,
  saml_properties ->> 'UserAccessUrl' as saml_user_access_url
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  json_extract(saml_properties, '$.RelayStateParameterName') as saml_relay_state_parameter_name,
  json_extract(saml_properties, '$.Status') as saml_status,
  json_extract(saml_properties, '$.UserAccessUrl') as saml_user_access_url
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```

### List the directories that have 'SwitchRunningMode' enabled
Determine the areas in which 'SwitchRunningMode' is enabled within your AWS Workspaces. This allows you to identify where users can switch between always-on and auto-stop modes, aiding in resource management and cost control.

```sql+postgres
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state,
  selfservice_permissions ->> 'SwitchRunningMode' as switch_running_mode
from
  aws_workspaces_directory
where
  selfservice_permissions ->> 'SwitchRunningMode' = 'ENABLED';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state,
  json_extract(selfservice_permissions, '$.SwitchRunningMode') as switch_running_mode
from
  aws_workspaces_directory
where
  json_extract(selfservice_permissions, '$.SwitchRunningMode') = 'ENABLED';
```

### Get the workspace creation properties of a particular directory
Analyze the settings to understand the configuration and properties of a specific workspace, such as internet access, maintenance mode, and user administrator status. This can be useful in auditing workspace settings and ensuring they align with company policies and security standards.

```sql+postgres
select
  name,
  directory_id,
  arn,
  workspace_creation_properties ->> 'CustomSecurityGroupId' as custom_security_group_id,
  workspace_creation_properties ->> 'DefaultOu' as default_ou,
  workspace_creation_properties ->> 'EnableInternetAccess' as enable_internet_access,
  workspace_creation_properties ->> 'EnableMaintenanceMode' as enable_maintenance_mode,
  workspace_creation_properties ->> 'EnableWorkDocs' as enable_work_docs,
  workspace_creation_properties ->> 'UserEnabledAsLocalAdministrator' as user_enabled_as_local_administrator
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  json_extract(workspace_creation_properties, '$.CustomSecurityGroupId') as custom_security_group_id,
  json_extract(workspace_creation_properties, '$.DefaultOu') as default_ou,
  json_extract(workspace_creation_properties, '$.EnableInternetAccess') as enable_internet_access,
  json_extract(workspace_creation_properties, '$.EnableMaintenanceMode') as enable_maintenance_mode,
  json_extract(workspace_creation_properties, '$.EnableWorkDocs') as enable_work_docs,
  json_extract(workspace_creation_properties, '$.UserEnabledAsLocalAdministrator') as user_enabled_as_local_administrator
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```

### List all registered directories
Explore which directories are registered in the AWS Workspaces service. This is useful for maintaining an overview of all active directories and ensuring they are in the correct state.

```sql+postgres
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  state = 'REGISTERED';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  alias,
  customer_user_name,
  directory_type,
  state
from
  aws_workspaces_directory
where
  state = 'REGISTERED';
```

### Get the workspace access properties of a particular directory
Explore which devices have access to a specific workspace directory. This can help in understanding the range of device types that can interact with the directory, allowing for better management and security planning.

```sql+postgres
select
  name,
  directory_id,
  arn,
  workspace_access_properties ->> 'DeviceTypeAndroid' as device_type_android,
  workspace_access_properties ->> 'DeviceTypeChromeOs' as device_type_chrome_os,
  workspace_access_properties ->> 'DeviceTypeIos' as device_type_ios,
  workspace_access_properties ->> 'DeviceTypeLinux' as device_type_linux,
  workspace_access_properties ->> 'DeviceTypeOsx' as device_type_osx,
  workspace_access_properties ->> 'DeviceTypeWeb' as device_type_web,
  workspace_access_properties ->> 'DeviceTypeWindows' as device_type_windows,
  workspace_access_properties ->> 'DeviceTypeZeroClient' as device_type_zero_client
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```

```sql+sqlite
select
  name,
  directory_id,
  arn,
  json_extract(workspace_access_properties, '$.DeviceTypeAndroid') as device_type_android,
  json_extract(workspace_access_properties, '$.DeviceTypeChromeOs') as device_type_chrome_os,
  json_extract(workspace_access_properties, '$.DeviceTypeIos') as device_type_ios,
  json_extract(workspace_access_properties, '$.DeviceTypeLinux') as device_type_linux,
  json_extract(workspace_access_properties, '$.DeviceTypeOsx') as device_type_osx,
  json_extract(workspace_access_properties, '$.DeviceTypeWeb') as device_type_web,
  json_extract(workspace_access_properties, '$.DeviceTypeWindows') as device_type_windows,
  json_extract(workspace_access_properties, '$.DeviceTypeZeroClient') as device_type_zero_client
from
  aws_workspaces_directory
where
  directory_id = 'd-96676995ea';
```