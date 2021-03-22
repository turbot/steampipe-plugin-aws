# Table: aws_config_conformance_pack

A conformance pack is a collection of AWS Config rules and remediation actions that can be easily deployed as a single entity in an account and a Region or across an organization in AWS Organizations.

## Examples

### Basic info

```sql
select 
    name, 
    conformance_pack_id, 
    input_parameters, 
    created_by_svc, 
    delivery_bucket, 
    delivery_bucket_prefix, 
    last_update, 
    title, 
    akas
from 
    aws.aws_config_conformance_pack;
```

