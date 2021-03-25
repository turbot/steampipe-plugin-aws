select association_id, association_version, association_name, document_name, document_version, targets, partition, region, account_id
from aws.aws_ssm_association
where association_id = '{{ output.resource_id.value }}';
