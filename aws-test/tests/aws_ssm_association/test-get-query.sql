select akas, association_id, association_version, name, partition
from aws.aws_ssm_association
where association_id = '{{ output.resource_id.value }}'
