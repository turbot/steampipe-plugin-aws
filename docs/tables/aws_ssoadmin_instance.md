---
title: "Steampipe Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL"
description: "Allows users to query AWS SSO Admin Instance, providing information about each AWS SSO instance in your AWS account."
folder: "SSO"
---

# Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL

The AWS SSO Admin Instance is a component of AWS Single Sign-On (SSO) service that enables you to manage SSO access to multiple AWS accounts and business applications. It simplifies the management of SSO access by centrally managing access to all of your AWS accounts and cloud applications. AWS SSO also includes built-in SAML integrations to many business applications.

## Table Usage Guide

The `aws_ssoadmin_instance` table in Steampipe provides you with information about each AWS SSO instance in your AWS account. This table allows you, as a DevOps engineer, to query instance-specific details, including the instance ARN, identity store ID, and associated metadata. You can utilize this table to gather insights on instances, such as instance status, instance creation time, and more. The schema outlines the various attributes of the SSO admin instance for you, including the instance ARN, identity store ID, and instance status.

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