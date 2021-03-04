# Table: aws_service_region

An AWS Region is a collection of AWS resources in a geographic area. Each AWS Region is isolated and independent of the other Regions. Regions provide fault tolerance, stability, and resilience, and can also reduce latency. An AWS Service is a named product that is available to users of Amazon Web Services.
This table provides a cross-reference of what AWS Services are available in each Region, and provided the available endpoints and protocols for that service.

## Examples

### AWS region info

```sql
select
  name,
  region_name,
  endpoint,
  protocols
from
  aws_service_region
where
  region_name = 'us-east-1'
  and name = 'sqs';
```

### Join with aws_service table

```sql
select
  sr.name,
  sr.region_name,
  s.long_name,
  sr.endpoint,
  sr.protocols,
  s.marketing_url
from
  aws_service_region sr,
  aws_service s
where
  sr.name = s.name
  and sr.region_name = 'us-east-1'
  and sr.name = 'sqs';
```
