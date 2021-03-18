select cluster_identifier, akas, allow_version_upgrade
from aws.aws_redshift_cluster
where cluster_identifier = '{{ output.resource_name.value }}';
