select registry_id, title, region
from aws.aws_ecr_registry
where region = '{{ output.aws_region.value }}';
