# download-terraform-poc
POC for a service to download and run Terraform binary. POC only for Terraform `Apply` and `Destroy`

## Get Started

POC uses example Terraform configuration file from [Terraform Getting Started Learn Module](https://learn.hashicorp.com/terraform/getting-started/install#quick-start-tutorial) which provisions NGINX server using Docker.

Requires: installing Terraform and Docker

To download dependencies: `make deps`

Edit Makefile variables that will ensure the right Terraform binary for your system is downloaded and where to place the binary.
- TF_PATH: path to move the Terraform binary to
- OS: your system's operating system
- ARCH: your system's architecture

To apply: `make apply`

To confirm NGINX and Docker resources were created:
- Visit NGINX server at `localhost:8000`
- Look at docker container `docker ps`

To destroy: `make destroy`

To modify terraform file: see `main.tf`
