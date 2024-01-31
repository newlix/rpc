# RPC

Simple RPC style APIs with generated clients & servers.

## About

All RPC methods are invoked with the POST method, and the RPC method name is placed in the URL path. Input is passed as a JSON object in the body, following a JSON response for the output as shown here:

```sh
$ curl -d '{ "project_id": "ping_production" }' https://api.example.com/get_alerts
{
  "alerts": [...]
}
```

All inputs are objects, all outputs are objects, this improves future-proofing as additional fields can be added without breaking existing clients. This is similar to the approach AWS takes with their APIs.

## Commands

There are several commands provided for generating clients, servers, and documentation. Each of these commands accept a `-schema` flag defaulting to `schema.json`, see the `-h` help output for additional usage details.

### Clients

- `rpc-dotnet-client` generates .NET clients
- `rpc-ruby-client` generates Ruby clients
- `rpc-php-client` generates PHP clients
- `rpc-elm-client` generates Elm clients
- `rpc-go-client` generates Go clients
- `rpc-go-types` generates Go type definitions
- `rpc-ts-client` generates TypeScript clients

### Servers

- `rpc-go-server` generates Go servers

### Documentation

- `rpc-md-docs` generates markdown documentation

## Schemas

Currently the schemas are loosely a superset of [JSON Schema](https://json-schema.org/), however, this is a work in progress. See the [example schema](./examples/todo/schema.json).

## FAQ

<details>
  <summary>Why did you create this project?</summary>
  There are many great options when it comes to building APIs, but to me the most important aspect is simplicity, for myself and for the end user. Simple JSON in, and JSON out is appropriate for 99% of my API work, there's no need for the additional performance provided by alternative encoding schemes, and rarely a need for more complex features such as bi-directional streaming provided by gRPC.
</details>

<details>
  <summary>Should I use this in production?</summary>
  Only if you're confident that it supports everything you need, or you're comfortable with forking. I created this project for my work at Apex Software, it may not suit your needs.
</details>

<details>
  <summary>Why JSON schemas?</summary>
  I think concise schemas using a DSL are great, until they're a limiting factor. Personally I have no problem with JSON, and it's easy to expand upon when you introduce a new feature, such as inline examples for documentation.
</details>

<details>
  <summary>Why doesn't it follow the JSON-RPC spec?</summary>
  I would argue this spec is outdated, there is little reason to support batching at the request level, as HTTP/2 handles this for you.
</details>


---

[![GoDoc](https://godoc.org/github.com/newlix/rpc?status.svg)](https://godoc.org/github.com/newlix/rpc)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)


