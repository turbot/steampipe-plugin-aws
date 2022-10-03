select p ->> 'KmsArn' as kms_arn, p ->> 'DigestAlgorithmMnemonic' as digest_algorithm_mnemonic, name, tags_src, query_logging_configs, p ->> 'SigningAlgorithmMnemonic' as signing_algorithm_mnemonic
from aws_route53_zone, jsonb_array_elements(dnssec_key_signing_keys) as p
where id = '{{ output.zone_id.value }}';