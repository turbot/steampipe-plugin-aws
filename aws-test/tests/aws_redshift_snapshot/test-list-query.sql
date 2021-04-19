select snapshot_identifier, cluster_identifier
from aws.aws_redshift_snapshot
where akas::text = '["{{output.resource_aka.value}}"]';
