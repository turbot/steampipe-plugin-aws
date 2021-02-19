select cluster_identifier, arn, allow_version_upgrade
from aws.aws_redshift_cluster
where cluster_identifier = '{{ resourceName }}'
