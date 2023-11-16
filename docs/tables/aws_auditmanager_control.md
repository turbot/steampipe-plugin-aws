---
title: "Table: aws_auditmanager_control - Query AWS Audit Manager Control using SQL"
description: "Allows users to query AWS Audit Manager Control data, providing information about controls within AWS Audit Manager. This table enables users to access detailed information about controls, such as control source, control type, description, and associated metadata."
---

# Table: aws_auditmanager_control - Query AWS Audit Manager Control using SQL

The `aws_auditmanager_control` table in Steampipe provides information about controls within AWS Audit Manager. This table allows DevOps engineers to query control-specific details, including control source, control type, description, and associated metadata. Users can utilize this table to gather insights on controls, such as their sources, types, descriptions, and more. The schema outlines the various attributes of the control, including the control id, name, type, source, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_auditmanager_control` table, you can use the `.inspect aws_auditmanager_control` command in Steampipe.

### Key columns:

- `id`: The unique identifier of the control. This can be used to join this table with other tables to get more detailed information about the control.
- `name`: The name of the control. This provides a human-readable identifier for the control, useful for understanding its purpose.
- `control_sources`: The sources of the control. This can provide context about the origin of the control and can be used to join with other tables that contain more details about the control sources.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  type
from
  aws_auditmanager_control;
```


### List custom audit manager controls

```sql
select
  name,
  id,
  type
from
  aws_auditmanager_control
where
  type = 'Custom';
```
