# General purpose targets

.PHONY: help
help: ## Display available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker compose targets.

.PHONY: docker-compose-up
docker-compose-up: ## Run all services locally using docker compose
	docker-compose -f docker-compose.yaml up -d --build

.PHONY: docker-compose-down
docker-compose-down: ## Run infra services locally using docker compose
	docker-compose -f docker-compose.yaml down

# Terraform targets.

.PHONY: infra
infra: kind-cluster-create argocd-services-create ## Initialize infrastructure

.PHONY: infra-destroy
infra-destroy: kind-cluster-destroy ## Destroy infrastructure

.PHONY: kind-cluster-create
kind-cluster-create: ## Create a kind cluster
	terraform -chdir=infra/terraform/cluster init
	terraform -chdir=infra/terraform/cluster plan -out=tfplan_cluster
	terraform -chdir=infra/terraform/cluster apply tfplan_cluster

.PHONY: kind-cluster-destroy
kind-cluster-destroy: ## Destroy a kind cluster
	terraform -chdir=infra/terraform/cluster init
	terraform -chdir=infra/terraform/cluster plan -destroy -out=tfplan_cluster
	terraform -chdir=infra/terraform/cluster apply tfplan_cluster

.PHONY: argocd-services-create
argocd-services-create: ## Create argocd services
	terraform -chdir=infra/terraform/services init
	terraform -chdir=infra/terraform/services plan -out=tfplan_services
	terraform -chdir=infra/terraform/services apply tfplan_services

.PHONY: argocd-services-destroy
argocd-services-destroy: ## Destroy argocd services
	terraform -chdir=infra/terraform/services init
	terraform -chdir=infra/terraform/services plan -destroy -out=tfplan_services
	terraform -chdir=infra/terraform/services apply tfplan_services