---
title: "Steampipe Table: aws_ec2_ssl_policy - Query AWS EC2 SSL Policies using SQL"
description: "Allows users to query AWS EC2 SSL Policies to retrieve detailed information about SSL policies used in AWS EC2 Load Balancers."
folder: "EC2"
---

# Table: aws_ec2_ssl_policy - Query AWS EC2 SSL Policies using SQL

The AWS EC2 SSL Policies are predefined security policies that determine the SSL/TLS protocol that an AWS EC2 instance uses when it's communicating with clients. These policies help to establish the ciphers and protocols that services like Elastic Load Balancing use when negotiating SSL/TLS connections. They can be customized to meet specific security requirements, ensuring secure and reliable client-to-server communications.

## Table Usage Guide

The `aws_ec2_ssl_policy` table in Steampipe provides you with information about SSL policies used in AWS Elastic Compute Cloud (EC2) Load Balancers. This table allows you as a developer or cloud architect to query SSL policy-specific details, including the policy name, the SSL protocols, and the cipher suite configurations. You can utilize this table to gather insights on the SSL policies, such as enabled SSL protocols, preferred cipher suites, and more. The schema outlines the various attributes of the SSL policy for you, including the policy name, the SSL protocols, the SSL ciphers, and the server order preference.

## Examples

### Basic info
Determine the areas in which your AWS EC2 instances are using certain SSL protocols. This can be beneficial for identifying potential security risks and ensuring that your instances are configured to use the most secure protocols.

```sql+postgres
select
  name,
  ssl_protocols
from
  aws_ec2_ssl_policy;
```

```sql+sqlite
select
  name,
  ssl_protocols
from
  aws_ec2_ssl_policy;
```


### List load balancer listeners that use an SSL policy with weak ciphers
Identify the load balancer listeners that are using an SSL policy with weak ciphers. This is beneficial for enhancing the security of your applications by pinpointing potential vulnerabilities.

```sql+postgres
select
  arn,
  ssl_policy
from
  aws_ec2_load_balancer_listener listener
join 
  aws_ec2_ssl_policy ssl_policy
on
  listener.ssl_policy = ssl_policy.Name
where
  ssl_policy.ciphers @> '[{"Name":"DES-CBC3-SHA"}]';
```

```sql+sqlite
select
  arn,
  ssl_policy
from
  aws_ec2_load_balancer_listener listener
join 
  aws_ec2_ssl_policy ssl_policy
on
  listener.ssl_policy = ssl_policy.Name
where
  json_extract(ssl_policy.ciphers, '$[*].Name') LIKE '%DES-CBC3-SHA%';
```