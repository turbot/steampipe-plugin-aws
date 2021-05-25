# Table: aws_audit_manager_control

AWS Audit Manager helps to continuously audit AWS usage to simplify how you assess risk and compliance with regulations and industry standards. With Audit Manager, it is easy to assess if policies, procedures, and activities – also known as controls – are operating effectively.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  type
from
  aws_audit_manager_control;
```


### List custom audit manager controls

```sql
select
  name,
  id,
  type
from
  aws_audit_manager_control
where
  type = 'Custom';
```
