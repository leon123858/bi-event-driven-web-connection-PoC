all:

start: 
	terraform init
	terraform plan -var-file="terraform.tfvars"
	terraform apply -var-file="terraform.tfvars"

format:
	terraform fmt -recursive 

init:
	terraform init

plan:
	terraform plan

destroy:
	terraform destroy