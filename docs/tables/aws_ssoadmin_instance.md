---
title: "Steampipe Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL"
description: "Allows users to query AWS SSO Admin Instance, providing information about each AWS SSO instance in your AWS account."
folder: "SSO"
---

# Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL

The AWS SSO Admin Instance is a component of AWS Single Sign-On (SSO) service that enables you to manage SSO access to multiple AWS accounts and business applications. It simplifies the management of SSO access by centrally managing access to all of your AWS accounts and cloud applications. AWS SSO also includes built-in SAML integrations to many business applications.

## Table Usage Guide

The `aws_ssoadmin_instance` table in Steampipe provides you with information about each AWS SSO instance in your AWS account. This table allows you, as a DevOps engineer, to query instance-specific details, including the instance ARN, name, identity store ID, creation date, owner account ID, status, and associated metadata. You can utilize this table to gather insights on instances, such as instance status, instance creation time, ownership, and more. The schema outlines the various attributes of the SSO admin instance for you, including the instance ARN, name, identity store ID, owner account ID, creation date, and instance status.

## Examples

### Basic info
Explore the instances where AWS Single Sign-On (SSO) admin resources are utilized to gain insights into their associated identity store IDs. This can be beneficial in managing access permissions and understanding user identities across different AWS accounts.

```sql+postgres
select
  arn,
  identity_store_id
from
  aws_ssoadmin_instance
```

```sql+sqlite
select
  arn,
  identity_store_id
from
  aws_ssoadmin_instance
```

### List instances with status and ownership details
Explore which AWS Identity Center (SSO) instances exist in your account, including their operational status, creation date, and ownership information. This is useful for auditing instance configurations and understanding your Identity Center deployment.

```sql+postgres
select
  name,
  status,
  created_date,
  owner_account_id,
  arn,
  identity_store_id,
  region
from
  aws_ssoadmin_instance;
```

```sql+sqlite
select
  name,
  status,
  created_date,
  owner_account_id,
  arn,
  identity_store_id,
  region
from
  aws_ssoadmin_instance;
```

### List instances ordered by creation date
Identify the most recently created Identity Center instances to track deployment history and understand when SSO was enabled in your organization.

```sql+postgres
select
  name,
  status,
  created_date,
  owner_account_id,
  region
from
  aws_ssoadmin_instance
order by
  created_date desc;
```

```sql+sqlite
select
  name,
  status,
  created_date,
  owner_account_id,
  region
from
  aws_ssoadmin_instance
order by
  created_date desc;
```