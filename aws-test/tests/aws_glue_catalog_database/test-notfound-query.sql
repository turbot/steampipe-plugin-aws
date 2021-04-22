select title, akas, region, account_id
from aws.aws_glue_catalog_database
where name = '{{ resourceName }}::xzq';
