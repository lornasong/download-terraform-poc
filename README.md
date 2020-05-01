# download-terraform-poc
POC for a service to download and run Terraform binary. POC only for Terraform `Apply` and `Destroy`

## Get Started

POC uses example Terraform configuration file from [Terraform Getting Started Learn Module](https://learn.hashicorp.com/terraform/getting-started/install#quick-start-tutorial) which provisions NGINX server using Docker.

Requires: installing Terraform and Docker

To download dependencies: `make deps`

Edit Makefile `TF_PATH` variable with the path that the Terraform binary should be moved to.

To apply: `make apply`

To confirm NGINX and Docker resources were created:
- Visit NGINX server at `localhost:8000`
- Look at docker container `docker ps`

To destroy: `make destroy`

To modify terraform file: see `main.tf`
