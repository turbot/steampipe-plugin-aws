select account_id, akas, region, title
from aws.aws_synthetics_canary
where name = '{{ output.resource_id.value }}' and region = '{{ output.aws_region.value }}';
