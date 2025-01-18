# OPA Policies - Terraform Cloud

:construction: **This repository is a work in progress, and is used by Damien to carry out various actions in his personal projects!** :construction:

![Go](https://img.shields.io/badge/Go-1.20-blue) ![Terraform](https://img.shields.io/badge/Terraform-1.5.0-purple) ![OPA](https://img.shields.io/badge/OPA-0.45.0-green)

## Overview

This repository serves as the omnibus OPA (Open Policy Agent) policy combiner for managing and enforcing compliance, security, and best practices in Terraform-based infrastructure as code (IaC) projects. The repository contains a collection of Rego policies tailored to address specific use cases, including:

- **AWS Resources**: Policies to enforce tagging, secure configurations, and instance restrictions.
- **DigitalOcean Droplets**: Policies to validate droplet sizes, enforce tagging, and disallow public IPs.
- **GCP Resources**: Policies to enforce secure configurations, resource naming conventions, and proper IAM roles.

## Features

- **Centralized Policy Management**: Combines multiple Rego policies into a single framework for Terraform plan validation.
- **Customizable Rules**: Easily adapt rules for specific use cases or environments.
- **Compatibility**: Works seamlessly with Terraform plan JSON outputs.
- **Flexibility**: Supports both advisory and enforcement levels for policies.

## Usage

### Prerequisites

1. Install [Open Policy Agent (OPA)](https://www.openpolicyagent.org/docs/latest/).
2. Generate a Terraform plan JSON:

   ```bash
   terraform plan -out=tfplan.binary
   terraform show -json tfplan.binary > tfplan.json
   ```

3. Clone this repository:

   ```bash
   git clone https://github.com/damienjburks/opa-policies.git
   cd opa-policies
   ```

### Running Policies

Evaluate a specific Terraform plan JSON against the policies:

```bash
opa eval -i tfplan.json -d . "data.aws.terraform.deny"
```

### Example Output

If the Terraform plan violates any policies, you will receive a detailed list of issues:

```json
[
  "S3 bucket 'aws_s3_bucket.my_bucket' must have versioning enabled.",
  "EC2 instance 'aws_instance.my_instance' must use instance type 't2.micro'."
]
```

### Policy Structure

- **AWS Policies**: Located in `aws.rego`
- **DigitalOcean Policies**: Located in `digitalocean_droplet.rego`
- **GCP Policies**: Located in `gcp.rego`

### Testing Policies

Use the following command to test individual policies:

```bash
opa eval -i tfplan.json -d <policy-file> "data.<package>.deny"
```

### Debugging Policies

Use OPA's trace feature to debug:

```bash
opa eval -i tfplan.json -d <policy-file> --log-level debug "data.<package>.deny"
```
