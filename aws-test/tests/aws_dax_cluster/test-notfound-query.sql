select title, akas, region, account_id
from aws.aws_dax_cluster
where cluster_name = '{{ output.resource_name.value }}1p';