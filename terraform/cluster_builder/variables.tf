variable "mongodb_atlas" {
  description = "MongoDB Atlas Access"
  type = object({
    active              = bool
    public_key          = string
    private_key         = string
    project_id          = string
    cluster_name        = string
    collection          = string
    database            = string
    region              = string
    username            = string
    password            = string
    role                = string
    cloud_provider_name = string
    instance_size       = string
  })
}

