select akas, association_id, association_version, partition
from aws.aws_ssm_association
where association_id = '{{ output.resource_id.value }}';
