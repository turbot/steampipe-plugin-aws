# Table: aws_workspaces_directory

Amazon WorkSpaces Directory refers to the service that manages user identities and access control for Amazon WorkSpaces. Amazon WorkSpaces is a cloud-based desktop computing service that allows users to access virtual desktops from anywhere, on any device. The directory in Amazon WorkSpaces serves as the central source for user authentication and management. It can be thought of as the repository of user accounts and permissions that control who can access and use WorkSpaces instances.

## Examples

## Basic info

```sql
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

```sql
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

### List directories of a particular type

```sql
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

```sql
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

### List the directories that have 'SwitchRunningMode' enabled

```sql
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

### Get the workspace creation properties of a particular directory

```sql
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

### List all registered directories

```sql
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

```sql
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
