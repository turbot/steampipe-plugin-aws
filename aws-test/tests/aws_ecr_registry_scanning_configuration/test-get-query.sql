select registry_id, scanning_configuration, title, region
from aws.aws_ecr_registry_scanning_configuration
where region = '{{ output.aws_region.value }}'