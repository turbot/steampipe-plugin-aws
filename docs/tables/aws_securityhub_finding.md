# Table: aws_securityhub_finding

AWS Security Hub eliminates the complexity of addressing large volumes of findings from multiple providers. It reduces the effort required to manage and improve the security of all of your AWS accounts, resources, and workloads.

## Examples

### Basic info

```sql
select
  name,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from
  aws_securityhub_finding;
```

### List findings with high severity

```sql
select
  name,
  product_arn,
  product_name,
  severity ->> 'Original' as severity_original
from
  aws_securityhub_finding
where
  severity ->> 'Original' = 'HIGH'
```

### List findings with failed compliance

```sql
select
  name,
  product_arn,
  product_name,
  compliance ->> 'Status' as compliance_status,
  compliance ->> 'StatusReasons'as compliance_status_reason
from
  aws_securityhub_finding
where
  compliance ->> 'Status' = 'FAILED'
```