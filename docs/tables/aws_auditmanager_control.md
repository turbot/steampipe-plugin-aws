---
title: "Steampipe Table: aws_auditmanager_control - Query AWS Audit Manager Control using SQL"
description: "Allows users to query AWS Audit Manager Control data, providing information about controls within AWS Audit Manager. This table enables users to access detailed information about controls, such as control source, control type, description, and associated metadata."
folder: "Audit Manager"
---

# Table: aws_auditmanager_control - Query AWS Audit Manager Control using SQL

The AWS Audit Manager Control is a feature within AWS Audit Manager that allows you to evaluate how well your AWS resource configurations align with established best practices. It helps you to simplify the compliance process and reduce risk by automating the collection of evidence of your AWS resource compliance with regulations and standards. The control feature allows for continuous auditing to ensure ongoing compliance.

## Table Usage Guide

The `aws_auditmanager_control` table in Steampipe provides you with information about controls within AWS Audit Manager. This table allows you, as a DevOps engineer, to query control-specific details, including control source, control type, description, and associated metadata. You can utilize this table to gather insights on controls, such as their sources, types, descriptions, and more. The schema outlines the various attributes of the control for you, including the control id, name, type, source, description, and associated tags.

**Important Notes**
- This table by default returns the `Standard` controls.
- You **must** specify `type` in a `where` clause to retrieve other control types. For more information, please refer to the [list of controls by specific type](https://docs.aws.amazon.com/audit-manager/latest/APIReference/API_ListControls.html#API_ListControls_RequestSyntax).

## Examples

### Basic info
Explore the basic information about the controls in AWS Audit Manager to understand their purpose and type. This can help in managing and assessing your AWS resources and environment effectively.

```sql+postgres
select
  name,
  id,
  description,
  type
from
  aws_auditmanager_control;
```

```sql+sqlite
select
  name,
  id,
  description,
  type
from
  aws_auditmanager_control;
```


### List custom audit manager controls
Discover the segments that consist of custom audit manager controls in your AWS environment. This can be particularly useful for understanding and managing your custom security and compliance configurations.

```sql+postgres
select
  name,
  id,
  type
from
  aws_auditmanager_control
where
  type = 'Custom';
```

```sql+sqlite
select
  name,
  id,
  type
from
  aws_auditmanager_control
where
  type = 'Custom';
```