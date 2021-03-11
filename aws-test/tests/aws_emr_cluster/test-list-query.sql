select id, cluster_arn, name, auto_terminate
from aws_new.aws_emr_cluster
where name = '{{ resourceName }}';