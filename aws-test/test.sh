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



run_test aws_accessanalyzer_analyzer		
run_test aws_directory_service_directory		
run_test aws_eks_addon_version			
run_test aws_kinesisanalyticsv2_application	
run_test aws_sns_topic_subscription
run_test aws_acm_certificate			
run_test aws_dms_replication_instance		
run_test aws_eks_cluster				
run_test aws_kms_key				
run_test aws_sqs_queue
run_test aws_api_gateway_api_key			
run_test aws_dynamodb_table			
run_test aws_eks_identity_provider_config	
run_test aws_lambda_alias			
run_test aws_ssm_association
run_test aws_api_gateway_authorizer		
run_test aws_ebs_snapshot			
run_test aws_elastic_beanstalk_application	
run_test aws_lambda_function			
run_test aws_ssm_document
run_test aws_api_gateway_rest_api		
run_test aws_ebs_volume				
run_test aws_elastic_beanstalk_environment	
run_test aws_lambda_version			
run_test aws_ssm_maintenance_window
run_test aws_api_gateway_stage			
run_test aws_ec2_ami				
run_test aws_elasticache_cluster			
run_test aws_macie2_classification_job		
run_test aws_ssm_managed_instance
run_test aws_api_gateway_usage_plan		
run_test aws_ec2_ami_shared			
run_test aws_elasticache_parameter_group		
run_test aws_rds_db_cluster			
run_test aws_ssm_managed_instance_compliance
run_test aws_api_gatewayv2_api			
run_test aws_ec2_application_load_balancer	
run_test aws_elasticache_replication_group	
run_test aws_rds_db_cluster_parameter_group	
run_test aws_ssm_parameter
run_test aws_api_gatewayv2_domain_name		
run_test aws_ec2_autoscaling_group		
run_test aws_elasticache_subnet_group		
run_test aws_rds_db_cluster_snapshot		
run_test aws_ssm_patch_baseline
run_test aws_api_gatewayv2_integration		
run_test aws_ec2_capacity_reservation		
run_test aws_elasticsearch_domain		
run_test aws_rds_db_event_subscription		
run_test aws_tagging_resource
run_test aws_api_gatewayv2_stage			
run_test aws_ec2_classic_load_balancer		
run_test aws_emr_cluster				
run_test aws_rds_db_instance			
run_test aws_vpc
run_test aws_appautoscaling_target		
run_test aws_ec2_gateway_load_balancer		
run_test aws_emr_instance_group			
run_test aws_rds_db_option_group			
run_test aws_vpc_customer_gateway
run_test aws_auditmanager_assessment		
run_test aws_ec2_instance			
run_test aws_eventbridge_rule			
run_test aws_rds_db_parameter_group		
run_test aws_vpc_dhcp_options
run_test aws_auditmanager_control		
run_test aws_ec2_key_pair			
run_test aws_glacier_vault			
run_test aws_rds_db_snapshot			
run_test aws_vpc_egress_only_internet_gateway
run_test aws_auditmanager_framework		
run_test aws_ec2_launch_configuration		
run_test aws_glue_catalog_database		
run_test aws_rds_db_subnet_group			
run_test aws_vpc_eip
run_test aws_backup_plan				
run_test aws_ec2_load_balancer_listener		
run_test aws_guardduty_detector			
run_test aws_redshift_cluster			
run_test aws_vpc_endpoint
run_test aws_backup_selection			
run_test aws_ec2_network_interface		
run_test aws_guardduty_finding			
run_test aws_redshift_event_subscription		
run_test aws_vpc_endpoint_service
run_test aws_backup_vault			
run_test aws_ec2_network_load_balancer		
run_test aws_guardduty_ipset			
run_test aws_redshift_parameter_group		
run_test aws_vpc_flow_log
run_test aws_cloudformation_stack		
run_test aws_ec2_regional_settings		
run_test aws_guardduty_threat_intel_set		
run_test aws_redshift_snapshot			
run_test aws_vpc_internet_gateway
run_test aws_cloudfront_cache_policy		
run_test aws_ec2_reserved_instance		
run_test aws_iam_access_advisor			
run_test aws_redshift_subnet_group		
run_test aws_vpc_nat_gateway
run_test aws_cloudfront_distribution		
run_test aws_ec2_ssl_policy			
run_test aws_iam_access_key			
run_test aws_region				
run_test aws_vpc_network_acl
run_test aws_cloudfront_origin_access_identity	
run_test aws_ec2_target_group			
run_test aws_iam_account_password_policy		
run_test aws_route53_record			
run_test aws_vpc_route
run_test aws_cloudfront_origin_request_policy	
run_test aws_ec2_transit_gateway			
run_test aws_iam_account_summary			
run_test aws_route53_resolver_endpoint		
run_test aws_vpc_route_table
run_test aws_cloudtrail_trail			
run_test aws_ec2_transit_gateway_route		
run_test aws_iam_group				
run_test aws_route53_resolver_rule		
run_test aws_vpc_security_group
run_test aws_cloudwatch_alarm			
run_test aws_ec2_transit_gateway_route_table	
run_test aws_iam_policy				
run_test aws_route53_zone			
run_test aws_vpc_subnet
run_test aws_cloudwatch_log_group		
run_test aws_ec2_transit_gateway_vpc_attachment	
run_test aws_iam_policy_simulator		
run_test aws_s3_access_point			
run_test aws_vpc_vpn_connection
run_test aws_cloudwatch_log_stream		
run_test aws_ecr_repository			
run_test aws_iam_role				
run_test aws_s3_bucket				
run_test aws_vpc_vpn_gateway
run_test aws_codebuild_project			
run_test aws_ecrpublic_repository		
run_test aws_iam_server_certificate		
run_test aws_sagemaker_endpoint_configuration	
run_test aws_waf_rate_based_rule
run_test aws_codebuild_source_credential		
run_test aws_ecs_cluster				
run_test aws_iam_user				
run_test aws_sagemaker_model			
run_test aws_waf_rule
run_test aws_codecommit_repository		
run_test aws_ecs_service				
run_test aws_inspector_assessment_target		
run_test aws_sagemaker_notebook_instance		
run_test aws_wafv2_ip_set
run_test aws_codepipeline_pipeline		
run_test aws_ecs_task_definition			
run_test aws_inspector_assessment_template	
run_test aws_secretsmanager_secret		
run_test aws_wafv2_regex_pattern_set
run_test aws_config_configuration_recorder	
run_test aws_efs_access_point			
run_test aws_kinesis_consumer			
run_test aws_securityhub_hub			
run_test aws_wafv2_rule_group
run_test aws_config_conformance_pack		
run_test aws_efs_file_system			
run_test aws_kinesis_firehose_delivery_stream	
run_test aws_securityhub_product			
run_test aws_wafv2_web_acl
run_test aws_config_rule				
run_test aws_efs_mount_target			
run_test aws_kinesis_stream			
run_test aws_securityhub_standards_subscription	
run_test aws_wellarchitected_workload
run_test aws_dax_cluster				
run_test aws_eks_addon				
run_test aws_kinesis_video_stream		
run_test aws_sns_topic

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt