# Table: aws_lambda_layer_version

A Lambda layer is a .zip file archive that can contain additional code or data. A layer can contain libraries, a custom runtime, data, or configuration files. Layers promote code sharing and separation of responsibilities so that you can iterate faster on writing business logic.

Layers can have one or more version. When you create a layer, Lambda sets the layer version to version 1. You can configure permissions on an existing layer version, but to update the code or make other configuration changes, you must create a new version of the layer.

## Examples

### Basic Info

```sql
select
  layer_arn,
  layer_name,
  layer_version_arn,
  created_date,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std,
  version
from
  aws_lambda_layer_version;
```
