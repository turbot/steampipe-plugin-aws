select replication_instance_identifier, arn, replication_instance_class,publicly_accessible, title, region
from aws_dms_replication_instance
where akas::text = '["{{ output.resource_aka.value }}"]';