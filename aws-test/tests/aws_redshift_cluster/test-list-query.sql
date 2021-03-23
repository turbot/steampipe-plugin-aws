select cluster_identifier, akas
from aws.aws_redshift_cluster
where akas::text = '["{{output.resource_aka.value}}"]';
