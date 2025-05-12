---
title: "Steampipe Table: aws_quicksight_account_setting - Query AWS QuickSight Account Settings using SQL"
description: "Allows users to query AWS QuickSight Account Settings, providing details about the QuickSight account configuration, edition, notification settings, and security features."
folder: "QuickSight"
---

# Table: aws_quicksight_account_setting - Query AWS QuickSight Account Settings using SQL

AWS QuickSight Account Settings represent the configuration and preferences set for your QuickSight account. These settings include account name, edition type, notification email, security features, and various other configurations that determine how QuickSight operates within your AWS environment.

## Table Usage Guide

The `aws_quicksight_account_setting` table in Steampipe provides you with information about account settings within AWS QuickSight. This table allows you, as a data analyst or administrator, to query account-specific details, including edition type, notification settings, security configurations, and various other account-level settings. You can utilize this table to gather insights on account configurations, such as termination protection status, public sharing settings, and notification preferences.

**Important Notes**
- You **_must_** specify `region` in a `where` clause in order to use this table.
- Account setting information for QuickSight is only available from the **identity region** (i.e., the region where the QuickSight account was initially created or enabled).
- Since there is no direct API to retrieve the identity region, users must provide it manually in the query to retrieve data successfully.

## Examples

### Basic info
Explore the basic configuration of your AWS QuickSight account to understand its setup and security settings.

```sql+postgres
select
  account_name,
  edition,
  default_namespace,
  notification_email,
  termination_protection_enabled,
  public_sharing_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1';
```

```sql+sqlite
select
  account_name,
  edition,
  default_namespace,
  notification_email,
  termination_protection_enabled,
  public_sharing_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1';
```

### Check accounts with termination protection disabled
Identify AWS QuickSight accounts that might be at risk due to disabled termination protection.

```sql+postgres
select
  account_name,
  edition,
  notification_email,
  termination_protection_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and not termination_protection_enabled;
```

```sql+sqlite
select
  account_name,
  edition,
  notification_email,
  termination_protection_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and termination_protection_enabled = 0;
```

### List accounts with public sharing enabled
Identify accounts that have public sharing enabled to assess potential data exposure risks.

```sql+postgres
select
  account_name,
  edition,
  notification_email,
  public_sharing_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and public_sharing_enabled;
```

```sql+sqlite
select
  account_name,
  edition,
  notification_email,
  public_sharing_enabled
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and public_sharing_enabled = 1;
```

### Get enterprise edition accounts
List all QuickSight accounts that are using the Enterprise edition to understand feature availability.

```sql+postgres
select
  account_name,
  edition,
  default_namespace,
  notification_email
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and edition = 'ENTERPRISE';
```

```sql+sqlite
select
  account_name,
  edition,
  default_namespace,
  notification_email
from
  aws_quicksight_account_setting
where
  region = 'us-east-1'
  and edition = 'ENTERPRISE';
```
