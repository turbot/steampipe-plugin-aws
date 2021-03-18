select cluster_identifier, akas
from aws.aws_redshift_cluster
where cluster_identifier = '{{ output.resource_name.value }}1p3';
