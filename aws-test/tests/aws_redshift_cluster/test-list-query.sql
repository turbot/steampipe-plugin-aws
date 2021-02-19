select cluster_identifier, arn
from aws.aws_redshift_cluster
where arn = '{{ output.resource_aka.value }}'
