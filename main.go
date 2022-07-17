package main

import (
	"context"
	"log"

	"github.com/metio/terraform-provider-git/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/metio/git",
	}

	err := providerserver.Serve(context.Background(), provider.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
