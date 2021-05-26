select id, akas, tags, status, domain_name, default_root_object, e_tag, price_class
from aws.aws_cloudfront_distribution
where id = '{{ output.resource_id.value }}';