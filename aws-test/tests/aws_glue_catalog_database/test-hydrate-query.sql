select akas, title
from aws.aws_glue_catalog_database
where name = '{{ resourceName }}';
