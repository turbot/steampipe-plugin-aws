---
title: "Table: aws_auditmanager_framework - Query AWS Audit Manager Framework using SQL"
description: "Allows users to query AWS Audit Manager Frameworks"
---

# Table: aws_auditmanager_framework - Query AWS Audit Manager Framework using SQL

The `aws_auditmanager_framework` table in Steampipe provides information about frameworks within AWS Audit Manager. This table allows DevOps engineers to query framework-specific details, including the framework's ARN, ID, type, and associated metadata. Users can utilize this table to gather insights on frameworks, such as the number of controls associated with each framework, the compliance type, and more. The schema outlines the various attributes of the Audit Manager Framework, including the framework ARN, creation date, last updated date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_auditmanager_framework` table, you can use the `.inspect aws_auditmanager_framework` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the framework. It can be used to join this table with other tables that contain framework ARNs.
- `id`: The unique identifier of the framework. It can be used to join this table with other tables that contain framework IDs.
- `type`: The type of the framework. It can be useful to join this table with other tables that contain framework types, allowing for more granular insights.

## Examples

### Basic info

```sql
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework;
```

### List custom audit manager frameworks

```sql
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework
where
  type = 'Custom';
```
