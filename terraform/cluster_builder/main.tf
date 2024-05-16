terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    version = "1.15.3" }
  }
}
provider "mongodbatlas" {
  public_key  = var.mongodb_atlas.public_key
  private_key = var.mongodb_atlas.private_key
}
resource "mongodbatlas_cluster" "josh-allan" {
  project_id = var.mongodb_atlas.project_id
  name       = var.mongodb_atlas.cluster_name

  //Provider Settings "block"
  provider_name               = var.mongodb_atlas.cloud_provider_name
  provider_instance_size_name = var.mongodb_atlas.instance_size
  provider_region_name        = var.mongodb_atlas.region
}

resource "mongodbatlas_database_user" "josh-allan" {
  username           = var.mongodb_atlas.username
  password           = var.mongodb_atlas.password
  project_id         = var.mongodb_atlas.project_id
  auth_database_name = "admin"

  roles {
    role_name     = var.mongodb_atlas.role
    database_name = "admin"
  }
}
output "connection_strings" {
  value = mongodbatlas_cluster.josh-allan.connection_strings.0.standard_srv
}
