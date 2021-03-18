select cluster_identifier, title, tags, akas
from aws.aws_redshift_cluster
where cluster_identifier = '{{ output.resource_name.value }}';
