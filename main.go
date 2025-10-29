// Copyright (c) 2025 - Stacklet, Inc.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/stacklet/terraform-provider-stacklet/internal/provider"
)

// Set by goreleaser during compilation.
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with debugger support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/stacklet/stacklet",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
