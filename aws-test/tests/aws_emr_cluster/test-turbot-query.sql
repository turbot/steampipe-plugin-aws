select title, tags, akas
from aws.aws_emr_cluster
where name = '{{ resourceName }}';