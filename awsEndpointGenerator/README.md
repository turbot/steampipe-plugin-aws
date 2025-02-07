
# **AWS Endpoint Generator**

## **Overview**
The **AWS Endpoint Generator** is a Go-based tool that automatically fetches AWS service endpoint data from the **AWS SDK for Go v2** and generates Go files containing:

1. **AWS-supported service endpoints** across partitions and regions.
2. **Service ID constants** to facilitate programmatic access to service metadata.

This tool helps in determining **AWS service availability** across different regions and partitions.

---

## **How It Works**
The tool retrieves endpoint data from the following AWS SDK JSON file:  
🔗 **[AWS Endpoints JSON](https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json)**  

It then **parses the JSON**, extracts relevant data, and generates two Go files:
- **`awsSupportedEndpoints/endpoints_gen.go`** → Contains partition, region, and service endpoint definitions.
- **`awsSupportedEndpoints/service_id_gen.go`** → Contains service ID constants for all AWS services.

---

## **Installation**
### **1. Move to the dir steampipe-plugin-aws**
```sh
cd steampipe-plugin-aws
```

#### **2. Updated `main.go` Content with:**
```go
package main

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-aws/awsEndpointGenerator"
)

func main() {
	if err := awsEndpointGenerator.Generate(); err != nil {
		fmt.Printf("Error generating Service supported endpoint file: %v\n", err)
	}

	if err := awsEndpointGenerator.GenerateServiceID(); err != nil {
		fmt.Printf("Error generating Service IDs file: %v\n", err)
	}
}
```


### **3. Install Dependencies**
Ensure Go modules are up to date:
```sh
go mod tidy
```

---

## **Usage**
### **1. Generate AWS Service Endpoints**
To fetch endpoint data and generate the Go file:
```sh
go run main.go
```
This will:
- Fetch and parse the **AWS endpoints JSON**.
- Generate **Go files** containing AWS-supported endpoints.

### **2. Generate AWS Service ID Constants**
To create a Go file containing **service ID constants**:
```sh
go run main.go
```
This generates `service_id_gen.go` and `endpoints_gen.go`, which helps in mapping AWS service names and endpoint details for all the AWS service.

---

## **Generated Files**
The tool generates the following files inside `awsSupportedEndpoints/`:

### **1. `endpoints_gen.go`** (AWS Service Endpoints)
Contains structured **Go definitions** of AWS partitions, regions, and service endpoints.

Example:
```go
var AWSPartition = Partition{
	ID:          "aws",
	Name:        "AWS Standard",
	DNSSuffix:   "amazonaws.com",
	Regions: map[string]Region{
		"us-east-1": {ID: "us-east-1", Description: "US East (N. Virginia)"},
		"us-west-2": {ID: "us-west-2", Description: "US West (Oregon)"},
	},
}
```

### **2. `service_id_gen.go`** (Service ID Constants)
Defines **constants** for AWS services.

Example:
```go
const AWS_S3_SERVICE_ID = "s3"
const AWS_EC2_SERVICE_ID = "ec2"
const AWS_LAMBDA_SERVICE_ID = "lambda"
```