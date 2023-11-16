---
title: "Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL"
description: "Allows users to query AWS SSO Admin Instance, providing information about each AWS SSO instance in your AWS account."
---

# Table: aws_ssoadmin_instance - Query AWS SSO Admin Instance using SQL

The `aws_ssoadmin_instance` table in Steampipe provides information about each AWS SSO instance in your AWS account. This table allows DevOps engineers to query instance-specific details, including the instance ARN, identity store ID, and associated metadata. Users can utilize this table to gather insights on instances, such as instance status, instance creation time, and more. The schema outlines the various attributes of the SSO admin instance, including the instance ARN, identity store ID, and instance status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssoadmin_instance` table, you can use the `.inspect aws_ssoadmin_instance` command in Steampipe.

**Key columns**:

- `instance_arn`: The ARN (Amazon Resource Name) of the SSO instance. This can be used to join this table with other tables that include SSO instance ARNs.
- `identity_store_id`: The ID of the identity store. This can be useful for joining with other tables that include identity store IDs, allowing for more detailed queries across multiple tables.
- `status`: The status of the SSO instance. This can be useful for filtering instances based on their status, providing insights into the operational state of your SSO instances.

## Examples

### Basic info

```sql
select
  arn,
  identity_store_id
from
  aws_ssoadmin_instance
```
