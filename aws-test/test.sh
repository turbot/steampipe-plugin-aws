#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"
 
 
# Define your function here
run_test () {
   echo -e "${BLACK}Running $1 ${ENDCOLOR}"
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
# rm -rf output.txt failed_tests.txt passed_tests.txt
 date >> output.txt
 date >> failed_tests.txt
 date >> passed_tests.txt

run_test accessanalyzer_analyzer
run_test acm_certificate
run_test api_gateway_api_key
run_test api_gateway_authorizer
run_test api_gateway_rest_api
run_test api_gateway_stage
run_test api_gateway_usage_plan
run_test api_gatewayv2_api
run_test api_gatewayv2_domain_name
run_test api_gatewayv2_integration
run_test api_gatewayv2_stage
run_test appautoscaling_target
run_test auditmanager_assessment
run_test auditmanager_control
run_test auditmanager_framework
run_test backup_plan
run_test backup_selection
run_test backup_vault
run_test cloudformation_stack
run_test cloudfront_cache_policy
run_test cloudfront_distribution
run_test cloudfront_origin_access_identity
run_test cloudfront_origin_request_policy
run_test cloudtrail_trail
run_test cloudwatch_alarm
run_test cloudwatch_log_group
run_test cloudwatch_log_stream
run_test codebuild_project
run_test codebuild_source_credential
run_test codecommit_repository
run_test codepipeline_pipeline
run_test config_configuration_recorder
run_test config_conformance_pack
run_test config_rule
run_test dax_cluster
run_test directory_service_directory
run_test dms_replication_instance
run_test dynamodb_table
run_test ebs_snapshot
run_test ebs_volume
run_test ec2_ami
run_test ec2_ami_shared
run_test ec2_application_load_balancer
run_test ec2_autoscaling_group
run_test ec2_capacity_reservation
run_test ec2_classic_load_balancer
run_test ec2_gateway_load_balancer
run_test ec2_instance
run_test ec2_key_pair
run_test ec2_launch_configuration
run_test ec2_load_balancer_listener
run_test ec2_network_interface
run_test ec2_network_load_balancer
run_test ec2_regional_settings
run_test ec2_reserved_instance
run_test ec2_ssl_policy
run_test ec2_target_group
run_test ec2_transit_gateway
run_test ec2_transit_gateway_route
run_test ec2_transit_gateway_route_table
run_test ec2_transit_gateway_vpc_attachment
run_test ecr_repository
run_test ecrpublic_repository
run_test ecs_cluster
run_test ecs_service
run_test ecs_task_definition
run_test efs_access_point
run_test efs_file_system
run_test efs_mount_target
run_test eks_addon
run_test eks_addon_version
run_test eks_cluster
run_test eks_identity_provider_config
run_test elastic_beanstalk_application
run_test elastic_beanstalk_environment
run_test elasticache_cluster
run_test elasticache_parameter_group
run_test elasticache_replication_group
run_test elasticache_subnet_group
run_test elasticsearch_domain
run_test emr_cluster
run_test emr_instance_group
run_test eventbridge_bus
run_test eventbridge_rule
run_test fsx_file_system
run_test glacier_vault
run_test glue_catalog_database
run_test guardduty_detector
run_test guardduty_finding
run_test guardduty_ipset
run_test guardduty_threat_intel_set
run_test iam_access_advisor
run_test iam_access_key
run_test iam_account_password_policy
run_test iam_account_summary
run_test iam_group
run_test iam_policy
run_test iam_policy_simulator
run_test iam_role
run_test iam_server_certificate
run_test iam_user
run_test identitystore_group
run_test identitystore_user
run_test inspector_assessment_target
run_test inspector_assessment_template
run_test kinesis_consumer
run_test kinesis_firehose_delivery_stream
run_test kinesis_stream
run_test kinesis_video_stream
run_test kinesisanalyticsv2_application
run_test kms_key
run_test lambda_alias
run_test lambda_function
run_test lambda_version
run_test macie2_classification_job
run_test organizations_account
run_test rds_db_cluster
run_test rds_db_cluster_parameter_group
run_test rds_db_cluster_snapshot
run_test rds_db_event_subscription
run_test rds_db_instance
run_test rds_db_option_group
run_test rds_db_parameter_group
run_test rds_db_snapshot
run_test rds_db_subnet_group
run_test redshift_cluster
run_test redshift_event_subscription
run_test redshift_parameter_group
run_test redshift_snapshot
run_test redshift_subnet_group
run_test region
run_test route53_record
run_test route53_resolver_endpoint
run_test route53_resolver_rule
run_test route53_zone
run_test s3_access_point
run_test s3_bucket
run_test sagemaker_endpoint_configuration
run_test sagemaker_model
run_test sagemaker_notebook_instance
run_test secretsmanager_secret
run_test securityhub_hub
run_test securityhub_product
run_test securityhub_standards_subscription
run_test sfn_state_machine
run_test sfn_state_machine_execution
run_test sfn_state_machine_execution_history
run_test sns_topic
run_test sns_topic_subscription
run_test sqs_queue
run_test ssm_association
run_test ssm_document
run_test ssm_maintenance_window
run_test ssm_managed_instance
run_test ssm_managed_instance_compliance
run_test ssm_parameter
run_test ssm_patch_baseline
run_test ssoadmin_instance
run_test ssoadmin_managed_policy_attachment
run_test ssoadmin_permission_set
run_test tagging_resource
run_test vpc
run_test vpc_customer_gateway
run_test vpc_dhcp_options
run_test vpc_egress_only_internet_gateway
run_test vpc_eip
run_test vpc_endpoint
run_test vpc_endpoint_service
run_test vpc_flow_log
run_test vpc_internet_gateway
run_test vpc_nat_gateway
run_test vpc_network_acl
run_test vpc_route
run_test vpc_route_table
run_test vpc_security_group
run_test vpc_subnet
run_test vpc_vpn_connection
run_test vpc_vpn_gateway
run_test waf_rate_based_rule
run_test waf_rule
run_test wafv2_ip_set
run_test wafv2_regex_pattern_set
run_test wafv2_rule_group
run_test wafv2_web_acl
run_test wellarchitected_workload

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt