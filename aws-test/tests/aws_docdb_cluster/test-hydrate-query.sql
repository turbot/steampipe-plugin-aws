select db_cluster_identifier, arn, deletion_protection, endpoint, engine, master_user_name, multi_az, storage_encrypted, tags_src
from aws.aws_docdb_cluster
where db_cluster_identifier = '{{ resourceName }}'
