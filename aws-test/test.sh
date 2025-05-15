 #!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 YELLOW="\e[93m"
 ENDCOLOR="\e[0m"

# Define your function here
run_test () {
   echo -e "${YELLOW}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
 }

 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
rm -rf output.txt failed_tests.txt passed_tests.txt
date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt

# run_test aws_accessanalyzer_analyzer
# run_test aws_acm_certificate
# run_test aws_acmpca_certificate_authority
# run_test aws_amplify_app
# run_test aws_api_gateway_api_key
# run_test aws_api_gateway_authorizer
# run_test aws_api_gateway_method
# run_test aws_api_gateway_rest_api
# run_test aws_api_gateway_stage
# run_test aws_api_gateway_usage_plan
# run_test aws_api_gatewayv2_api
# run_test aws_api_gatewayv2_domain_name
# run_test aws_api_gatewayv2_integration
# run_test aws_api_gatewayv2_route
# run_test aws_api_gatewayv2_stage
# run_test aws_app_runner_service
# run_test aws_appautoscaling_policy
# run_test aws_appautoscaling_target
# run_test aws_appconfig_application
# run_test aws_appstream_fleet
# run_test aws_appsync_graphql_api
# run_test aws_auditmanager_assessment
# run_test aws_auditmanager_control
# run_test aws_auditmanager_framework
# run_test aws_backup_framework
# run_test aws_backup_plan
# run_test aws_backup_report_plan
# run_test aws_backup_selection
# run_test aws_backup_vault
# run_test aws_cloudformation_stack
# run_test aws_cloudformation_stack_resource
# run_test aws_cloudformation_stack_set
# run_test aws_cloudfront_cache_policy
# run_test aws_cloudfront_distribution
# run_test aws_cloudfront_function
# run_test aws_cloudfront_origin_access_identity
# run_test aws_cloudfront_origin_request_policy
# run_test aws_cloudfront_response_headers_policy
# run_test aws_cloudsearch_domain
# run_test aws_cloudtrail_event_data_store
# run_test aws_cloudtrail_trail
# run_test aws_cloudwatch_alarm
# run_test aws_cloudwatch_log_group
# run_test aws_cloudwatch_log_resource_policy
# run_test aws_cloudwatch_log_stream
# run_test aws_codeartifact_domain
# run_test aws_codeartifact_repository
# run_test aws_codebuild_project
# run_test aws_codebuild_source_credential
# run_test aws_codecommit_repository
# run_test aws_codedeploy_app
# run_test aws_codedeploy_deployment_config
# run_test aws_codedeploy_deployment_group
# run_test aws_codepipeline_pipeline
# run_test aws_config_aggregate_authorization
# run_test aws_config_configuration_recorder
# run_test aws_config_conformance_pack
# run_test aws_config_rule
# run_test aws_dax_cluster
# run_test aws_dax_parameter
# run_test aws_dax_parameter_group
# run_test aws_dax_subnet_group
# run_test aws_directory_service_directory
# run_test aws_dlm_lifecycle_policy
# run_test aws_dms_certificate
# run_test aws_dms_endpoint
# run_test aws_dms_replication_instance
# run_test aws_dms_replication_task
# run_test aws_docdb_cluster
# run_test aws_docdb_cluster_instance
# run_test aws_docdb_cluster_snapshot
# run_test aws_dynamodb_table
# run_test aws_ebs_snapshot
# run_test aws_ebs_volume
# run_test aws_ec2_ami
# run_test aws_ec2_ami_shared
# run_test aws_ec2_application_load_balancer
# run_test aws_ec2_autoscaling_group
# run_test aws_ec2_capacity_reservation
# run_test aws_ec2_classic_load_balancer
# run_test aws_ec2_gateway_load_balancer
# run_test aws_ec2_instance
# run_test aws_ec2_key_pair
# run_test aws_ec2_launch_configuration
# run_test aws_ec2_launch_template
# run_test aws_ec2_launch_template_version
# run_test aws_ec2_load_balancer_listener
# run_test aws_ec2_managed_prefix_list
# run_test aws_ec2_managed_prefix_list_entry
# run_test aws_ec2_network_interface
# run_test aws_ec2_network_load_balancer
# run_test aws_ec2_regional_settings
# run_test aws_ec2_ssl_policy
# run_test aws_ec2_target_group
# run_test aws_ec2_transit_gateway
# run_test aws_ec2_transit_gateway_route
# run_test aws_ec2_transit_gateway_route_table
# run_test aws_ec2_transit_gateway_vpc_attachment
# run_test aws_ecr_registry_scanning_configuration
# run_test aws_ecr_repository
# run_test aws_ecrpublic_repository
# run_test aws_ecs_cluster
# run_test aws_ecs_service
# run_test aws_ecs_task_definition
# run_test aws_efs_access_point
# run_test aws_efs_file_system
# run_test aws_efs_mount_target
# run_test aws_eks_addon
# run_test aws_eks_addon_version
# run_test aws_eks_cluster
# run_test aws_eks_fargate_profile
# run_test aws_eks_identity_provider_config
# run_test aws_elastic_beanstalk_application
# run_test aws_elastic_beanstalk_application_version
# run_test aws_elastic_beanstalk_environment
# run_test aws_elasticache_cluster
# run_test aws_elasticache_parameter_group
# run_test aws_elasticache_replication_group
# run_test aws_elasticache_subnet_group
# run_test aws_elasticsearch_domain
# run_test aws_emr_cluster
# run_test aws_emr_instance_fleet
# run_test aws_emr_instance_group
# run_test aws_emr_security_configuration
# run_test aws_eventbridge_bus
# run_test aws_eventbridge_rule
# run_test aws_fsx_file_system
# run_test aws_glacier_vault
# run_test aws_glue_catalog_database
# run_test aws_glue_catalog_table
# run_test aws_glue_connection
# run_test aws_glue_crawler
# run_test aws_glue_data_catalog_encryption_settings
# run_test aws_glue_dev_endpoint
# run_test aws_glue_job
# run_test aws_glue_security_configuration
# run_test aws_guardduty_detector
# run_test aws_guardduty_filter
# run_test aws_guardduty_finding
# run_test aws_guardduty_ipset
# run_test aws_guardduty_member
# run_test aws_guardduty_publishing_destination
# run_test aws_guardduty_threat_intel_set
# run_test aws_iam_access_advisor
# run_test aws_iam_access_key
# run_test aws_iam_account_password_policy
# run_test aws_iam_account_summary
# run_test aws_iam_group
# run_test aws_iam_open_id_connect_provider
# run_test aws_iam_policy
# run_test aws_iam_policy_simulator
# run_test aws_iam_role
# run_test aws_iam_saml_provider
# run_test aws_iam_server_certificate
# run_test aws_iam_service_specific_credential
# run_test aws_iam_user
# run_test aws_identitystore_group
# run_test aws_identitystore_user
# run_test aws_inspector_assessment_target
# run_test aws_inspector_assessment_template
# run_test aws_iot_thing
# run_test aws_iot_thing_group
# run_test aws_iot_thing_type
# run_test aws_keyspaces_table
# run_test aws_kinesis_consumer
# run_test aws_kinesis_firehose_delivery_stream
# run_test aws_kinesis_stream
# run_test aws_kinesis_video_stream
# run_test aws_kinesisanalyticsv2_application
# run_test aws_kms_alias
# run_test aws_kms_key
# run_test aws_lambda_alias
# run_test aws_lambda_function
# run_test aws_lambda_layer_version
# run_test aws_lambda_version
# run_test aws_lightsail_bucket
# run_test aws_lightsail_instance
# run_test aws_macie2_classification_job
# run_test aws_media_store_container
# run_test aws_memorydb_cluster
# run_test aws_mq_broker
# run_test aws_msk_serverless_cluster
# run_test aws_neptune_db_cluster
# run_test aws_neptune_db_cluster_snapshot
# run_test aws_networkfirewall_firewall
# run_test aws_networkfirewall_firewall_policy
# run_test aws_opensearch_domain
# run_test aws_organizations_account
# run_test aws_pinpoint_app
# run_test aws_ram_principal_association
# run_test aws_ram_resource_association
# run_test aws_rds_db_cluster
# run_test aws_rds_db_cluster_parameter_group
# run_test aws_rds_db_cluster_snapshot
# run_test aws_rds_db_event_subscription
# run_test aws_rds_db_instance
# run_test aws_rds_db_option_group
# run_test aws_rds_db_parameter_group
# run_test aws_rds_db_proxy
# run_test aws_rds_db_snapshot
# run_test aws_rds_db_subnet_group
# run_test aws_redshift_cluster
# run_test aws_redshift_event_subscription
# run_test aws_redshift_parameter_group
# run_test aws_redshift_snapshot
# run_test aws_redshift_subnet_group
# run_test aws_redshiftserverless_namespace
# run_test aws_redshiftserverless_workgroup
# run_test aws_region
# run_test aws_route53_health_check
# run_test aws_route53_query_log
# run_test aws_route53_record
# run_test aws_route53_resolver_endpoint
# run_test aws_route53_resolver_query_log_config
# run_test aws_route53_resolver_rule
# run_test aws_route53_traffic_policy
# run_test aws_route53_traffic_policy_instance
# run_test aws_route53_vpc_association_authorization
# run_test aws_route53_zone
# run_test aws_s3_access_point
# run_test aws_s3_bucket
# run_test aws_sagemaker_app
# run_test aws_sagemaker_domain
# run_test aws_sagemaker_endpoint_configuration
# run_test aws_sagemaker_model
# run_test aws_sagemaker_notebook_instance
# run_test aws_secretsmanager_secret
# run_test aws_securityhub_action_target
# run_test aws_securityhub_hub
# run_test aws_securityhub_insight
# run_test aws_securityhub_member
# run_test aws_securityhub_product
# run_test aws_securityhub_standards_control
# run_test aws_securityhub_standards_subscription
# run_test aws_serverlessapplicationrepository_application
# run_test aws_service_discovery_instance
# run_test aws_service_discovery_namespace
run_test aws_service_discovery_service
run_test aws_servicecatalog_portfolio
run_test aws_servicecatalog_product
run_test aws_servicequotas_default_service_quota
run_test aws_servicequotas_service_quota
run_test aws_ses_email_identity
run_test aws_sfn_state_machine
run_test aws_sfn_state_machine_execution
run_test aws_sfn_state_machine_execution_history
run_test aws_sns_subscription
run_test aws_sns_topic
run_test aws_sns_topic_subscription
run_test aws_sqs_queue
run_test aws_ssm_association
run_test aws_ssm_document
run_test aws_ssm_maintenance_window
run_test aws_ssm_managed_instance
run_test aws_ssm_managed_instance_compliance
run_test aws_ssm_parameter
run_test aws_ssm_patch_baseline
run_test aws_ssoadmin_account_assignment
run_test aws_ssoadmin_instance
run_test aws_ssoadmin_managed_policy_attachment
run_test aws_ssoadmin_permission_set
run_test aws_tagging_resource
run_test aws_timestreamwrite_database
run_test aws_timestreamwrite_table
run_test aws_transfer_server
run_test aws_transfer_user
run_test aws_vpc
run_test aws_vpc_customer_gateway
run_test aws_vpc_dhcp_options
run_test aws_vpc_egress_only_internet_gateway
run_test aws_vpc_eip
run_test aws_vpc_endpoint
run_test aws_vpc_endpoint_service
run_test aws_vpc_flow_log
run_test aws_vpc_internet_gateway
run_test aws_vpc_nat_gateway
run_test aws_vpc_network_acl
run_test aws_vpc_peering_connection
run_test aws_vpc_route
run_test aws_vpc_route_table
run_test aws_vpc_security_group
run_test aws_vpc_subnet
run_test aws_vpc_vpn_connection
run_test aws_vpc_vpn_gateway
run_test aws_waf_rate_based_rule
run_test aws_waf_rule
run_test aws_waf_rule_group
run_test aws_waf_web_acl
run_test aws_wafregional_rule
run_test aws_wafregional_rule_group
run_test aws_wafregional_web_acl
run_test aws_wafv2_ip_set
run_test aws_wafv2_regex_pattern_set
run_test aws_wafv2_rule_group
run_test aws_wafv2_web_acl
run_test aws_wellarchitected_workload

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt