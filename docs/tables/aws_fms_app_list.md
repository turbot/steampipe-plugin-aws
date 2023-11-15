# Table: aws_fms_app_list

An AWS Network Firewall Application List is a custom list of domain names or IP addresses that you can use in Network Firewall rules. These lists are useful for allowing or blocking specific traffic based on the defined applications or addresses.

## Examples

### Basic info

```sql
select
  list_name,
  list_id,
  arn,
  create_time,
  creation_time
from
  aws_fms_app_list;
```

### List of apps created in last 30 days

```sql
select
  list_name,
  list_id,
  arn,
  create_time,
  creation_time
from
  aws_fms_app_list
where
  create_time >= now() - interval '30' day;
```

### Get application details of each app list

```sql
select
  list_name,
  list_id,
  a ->> 'AppName' as app_name,
  a ->> 'Port' as port,
  a ->> 'Protocol' as protocol
from
  aws_fms_app_list,
  jsonb_array_elements(apps_list -> 'AppsList') as a;

```