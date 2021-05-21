select akas, id, domain_name, is_ipv6_enabled, in_progress_invalidation_batches, e_tag, enabled, comment
from aws.aws_cloudfront_distribution
where akas = '["{{ output.resource_aka.value }}"]';