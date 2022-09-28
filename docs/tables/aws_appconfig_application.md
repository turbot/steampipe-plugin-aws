# Table: aws_appconfig_application

In AWS AppConfig, an application is simply an organizational construct like a folder. This organizational construct has a relationship with some unit of executable code.

## Examples

### Basic info

```sql
select
  arn,
  id,
  name,
  description,
  tags
from
  aws_appconfig_application;
```