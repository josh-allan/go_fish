A (very) WIP rewrite of my `ozbargain_hunter` python module, but in Go (with some additional functionality).

# Project Overview
The basic idea of this project is split into (currently) 3 components:
- `go_fish`
    - The parsing tool.
- `localstack`
    - Testing locally deployed infrastructure
- `terraform`
    - Infrastructure as Code handling the cloud deployments

A `.env` will need to be present at the root of `go_fish` as this _presently_ handles secrets. I'm currently noodling on options that don't involve handling secrets this way.

## Go Fish

The current heirarchy of the project is as follows:
```
|-- ./go_fish
|   |--.env
|   |-- ./go_fish/config
|   |   `-- ./go_fish/config/load_config.go
|   |-- ./go_fish/go.mod
|   |-- ./go_fish/parser
|   |   `-- ./go_fish/parser/parser.go
|   |-- ./go_fish/util
|   |   `-- ./go_fish/util/shared.go
|   |-- ./go_fish/db
|   |   `-- ./go_fish/db/mongo.go
|   |-- ./go_fish/go.sum
|   `-- ./go_fish/main.go
```

### High-level overview of the modules:
- `load_config.go` houses the configuration struct that sources and passes the values tored in the `.env` at the project root.
- `parser.go` is the handler for the parser logic, as well as normalising the case of the search terms.
- `shared.go` contains shared logic (such as the feed URLs, and the search terms that are being passed into the feed parser).
- `mongo.go` handles the db lookup logic, as well as the document insertion for matching entries.
- `main.go` the entrypoint for the program (as well as the combination of the above.). This also incorporates a webhook to a discord server to output the results.

### Running the package
The expected environment variables will all be stored in the config struct in `load_config.go` - ensure that the `.env` file is populated with these values.

`go get` to pull in all external modules.

Configure a discord webhook and store it under an env var.

`go run main.go` if you do not wish to compile the binary, otherwise `go build` && `./go_fish`

## Localstack

```
|-- ./localstack
|   `-- ./localstack/docker-compose.yml
```

This is a simple WIP implementation of running AWS services on premises for testing - the docker-compose file creates the entrypoint, whereas the endpoints are instantiated via Terraform.

## Terraform
```

|-- ./terraform
|   |-- ./terraform/cluster_builder
|   |   |-- ./terraform/cluster_builder/main.tf
|   |   `-- ./terraform/cluster_builder/variables.tf
|   `-- ./terraform/localstack
|       |-- ./terraform/localstack/main.tf
|       `-- ./terraform/localstack/variables.tf

```

There are two individual resources under this module. 

`cluster_builder` houses all of the logic to create a MongoDB Atlas shared-tier instance (this will require the secrets to be handled in a `terraform.tfvars` file targeting an M0 unless you wish to pay for the cluster) - this resource has been written with a Shared Tier cluster in mind. 

## TODO
With the above in mind - why is this still incomplete? 

I have plans to abstract logic out of the `main.go` package and ensure this remains as only an entrypoint to the program, as well as I would like to move towards a more vendor-agnostic approach by holding things in agnostic interfaces which vendor software (such as MongoDB) can access without actually resulting in breaking changes inside the logic itself.

I would also like to adopt some more modularisation to the code, such as ensuring that the search terms are read from a database, instead of strings in the code. With this in mind, this is very susceptible to breaking changes - so use at own risk, I cannot be certain of the behaviour. 
