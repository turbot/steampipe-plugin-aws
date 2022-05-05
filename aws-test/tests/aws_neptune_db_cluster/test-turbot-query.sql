select db_cluster_identifier, title, tags, akas
from aws_neptune_db_cluster
where db_cluster_identifier = '{{ resourceName }}';
