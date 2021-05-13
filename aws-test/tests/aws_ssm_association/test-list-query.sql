select association_id, association_name, arn, association_version, compliance_severity, overview, partition, region, schedule_expression
from aws.aws_ssm_association
where akas::text = '["{{ output.resource_aka.value }}"]';
