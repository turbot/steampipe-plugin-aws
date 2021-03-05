select association_id, association_name, association_version, compliance_severity, document_version, instance_id, name, overview, partition, region, schedule_expression, targets
from aws.aws_ssm_association
where akas::text = '["{{ output.resource_aka.value }}"]'
