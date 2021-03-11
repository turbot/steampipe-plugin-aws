select title, tags, akas
from aws_new.aws_emr_cluster
where name = '{{ resourceName }}';