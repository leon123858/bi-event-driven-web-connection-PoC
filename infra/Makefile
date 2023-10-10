all:

start: 
	terraform init
	terraform plan -var-file="terraform.tfvars"
	terraform apply -var-file="terraform.tfvars"
	# trigger all build imeadiately
	bash ./triggerCD.sh

format:
	terraform fmt -recursive 

init:
	terraform init

plan:
	terraform plan

destroy:
	terraform destroy