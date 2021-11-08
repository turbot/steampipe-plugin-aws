# Table: aws_lambda_layer

A Lambda layer is a .zip file archive that can contain additional code or data. A layer can contain libraries, a custom runtime, data, or configuration files. Layers promote code sharing and separation of responsibilities so that you can iterate faster on writing business logic.

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
