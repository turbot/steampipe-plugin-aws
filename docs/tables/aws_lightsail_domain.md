# aws_lightsail_domain

The `aws_lightsail_domain` table provides information about domains in AWS Lightsail. Lightsail domains are used to manage DNS records for your Lightsail resources.

## Examples

### Basic info
```sql
select
  domain_name,
  arn,
  created_at,
  is_managed,
  resource_type
from
  aws_lightsail_domain;
```

### List managed domains
```sql
select
  domain_name,
  arn,
  created_at
from
  aws_lightsail_domain
where
  is_managed;
```

### Get domain details by name
```sql
select
  domain_name,
  arn,
  created_at,
  is_managed,
  location,
  resource_type,
  support_code
from
  aws_lightsail_domain
where
  domain_name = 'example.com';
```

### List domains with tags
```sql
select
  domain_name,
  tags
from
  aws_lightsail_domain
where
  tags is not null;
```

## Schema

| Name | Type | Description |
|------|------|-------------|
| domain_name | string | The name of the domain |
| arn | string | The Amazon Resource Name (ARN) of the domain |
| created_at | timestamp | The timestamp when the domain was created |
| is_managed | boolean | Indicates whether the domain is managed by Lightsail |
| location | json | The AWS Region and Availability Zones where the domain is located |
| resource_type | string | The resource type |
| support_code | string | The support code for the domain |
| tags_src | json | A list of tags assigned to the domain |
| title | string | The title of the resource |
| tags | json | A map of tags for the resource |
| akas | json | Array of globally unique identifier strings (also known as) for the resource |

## References

* [AWS Lightsail Domains Documentation](https://docs.aws.amazon.com/lightsail/latest/userguide/understanding-lightsail-domains.html)
* [GetDomain API Documentation](https://docs.aws.amazon.com/lightsail/2016-11-28/api-reference/API_GetDomain.html)
* [GetDomains API Documentation](https://docs.aws.amazon.com/lightsail/2016-11-28/api-reference/API_GetDomains.html) 