A (very) WIP rewrite of my `ozbargain_hunter` python module, but in Go (with some additional functionality).

# Project Overview
The basic idea of this project is split into (currently) 3 components:
- `cli`
    - The cli tool.
- `localstack`
    - Testing locally deployed infrastructure
- `terraform`
    - Infrastructure as Code handling the cloud deployments

A `.env` will need to be present at the root of `go_fish` as this _presently_ handles secrets. I'm currently noodling on options that don't involve handling secrets this way.

## CLI

The current hierarchy of the project is as follows:
```
cli
   |-- cmd
   |   |-- main.go
   |-- go.mod
   |-- go.sum
   |-- internal
   |   |-- config
   |   |   |-- load_config.go
   |   |-- db
   |   |   |-- datastore.go
   |   |   |-- mongo.go
   |   |   |-- mongo_datastore.go
   |   |-- parser
   |   |   |-- parser.go
   |   |-- tasks
   |   |   |-- add-term
   |   |   |   |-- addTerm.go
   |   |   |-- list-terms
   |   |   |   |-- listTerm.go
   |   |   |-- scraper
   |   |   |   |-- scraper.go
   |   |-- util
   |   |   |-- shared.go
```

### High-level overview of the modules:
- `load_config.go` houses the configuration struct that sources and passes the values stored in the `.env` at the project root.
- `parser.go` is the handler for the parser logic, as well as normalising the case of the search terms.
- `shared.go` contains shared logic (such as the feed URLs, and the search terms that are being passed into the feed parser).
- `mongo.go` handles the db lookup logic, as well as the document insertion for matching entries.
- `addTerm.go` is the handler for adding search terms to the database.
- `listTerm.go` is the handler for listing the current search terms in the database.
- `scraper.go` is the handler for initialising the scraper. 
- `main.go` the entrypoint for the program (as well as the combination of the above.). This also incorporates a webhook to a discord server to output the results.

### Running the package
The expected environment variables will all be stored in the config struct in `load_config.go` - ensure that the `.env` file is populated with these values.

`go get` to pull in all external modules.

Configure a discord webhook and store it under an env var.

`go run cmd/main.go` if you do not wish to compile the binary, otherwise `go build` && `./go_fish`

`gofish insert` takes a string argument and inserts it into the `search_terms` collection.

`gofish list` will return all current entries in the `search_terms` collection.

`gofish scrape` will initialise the scraper.

## Localstack

```
localstack
   |-- docker-compose.yml
```

This is a simple WIP implementation of running AWS services on premises for testing - the docker-compose file creates the entrypoint, whereas the endpoints are instantiated via Terraform.

## Terraform
```

terraform
   |-- cluster_builder
   |   |-- .terraform.lock.hcl
   |   |-- main.tf
   |   |-- variables.tf
   |-- localstack
   |   |-- .terraform.lock.hcl
   |   |-- main.tf
   |   |-- variables.tf

```

There are two individual resources under this module. 

`cluster_builder` houses the logic to create a MongoDB Atlas shared-tier instance (this will require the secrets to be handled in a `terraform.tfvars` file targeting an M0 unless you wish to pay for the cluster) - this resource has been written with a Shared Tier cluster in mind. 

## TODO

I would also like to adopt some more modularisation to the code. I intend to abstract the database logic to an agnostic interface, therefore this is very susceptible to breaking changes - use at own risk, I cannot be certain of the behaviour. 
