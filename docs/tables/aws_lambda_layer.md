# Table: aws_lambda_layer

A Lambda layer is an archive containing additional code, such as libraries, dependencies, or even custom runtimes. When you include a layer in a function, the contents are extracted to the /opt directory in the execution environment.

## Examples

### Basic Info

```sql
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  jsonb_pretty(compatible_runtimes) as compatible_runtimes,
  jsonb_pretty(compatible_architectures) as compatible_architectures,
  version
from
  aws_lambda_layer;
```
