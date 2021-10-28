terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.66.0"
    }
  }
  backend "azurerm" {}
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

data "azurerm_resource_group" "group" {
  name = var.resource_group_name
}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                      = var.cluster_name
  location                  = data.azurerm_resource_group.group.location
  resource_group_name       = data.azurerm_resource_group.group.name
  dns_prefix                = "aks-create-action"
  automatic_channel_upgrade = "stable"
  sku_tier                  = "Paid"

  default_node_pool {
    name                = "app"
    vm_size             = "Standard_D8s_v3"
    enable_auto_scaling = true
    max_count           = 5
    min_count           = 2
    os_disk_type        = "Ephemeral"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Test"
  }
}


resource "azurerm_kubernetes_cluster_node_pool" "systempool" {
  name                  = "system"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.cluster.id
  vm_size               = "Standard_D8s_v3"
  enable_auto_scaling   = true
  max_count             = 5
  min_count             = 2
  mode                  = "System"
  os_disk_type          = "Ephemeral"
  os_disk_size_gb       = "200"
  node_taints = [
    "kubernetes.azure.com/CriticalAddonsOnly=true:NoSchedule"
  ]
}

resource "azurerm_container_registry" "acr" {
  count               = var.create_acr ? 1 : 0
  name                = var.cluster_name
  location            = data.azurerm_resource_group.group.location
  resource_group_name = data.azurerm_resource_group.group.name
  sku                 = "Basic"
  admin_enabled       = false
}

variable "create_acr" {
  type    = bool
  default = false
}

variable "resource_group_name" {
  type = string
}

variable "cluster_name" {
  type = string
}