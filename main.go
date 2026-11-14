// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run gen/generator.go

//go:generate go run golang.org/x/tools/cmd/goimports -w internal/provider/

//go:generate terraform fmt -recursive ./examples/

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

//go:generate go run gen/doc_category.go

var (
	// Set by goreleaser.
	version string = "dev"
)

func main() {

	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/CiscoDevNet/crosswork",
	}

	if debug {
		log.Println("[DEBUG] Provider starting in debug mode...")
		opts.Debug = true
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	// Serve returns when Terraform stops the provider; release the Crosswork sessions
	provider.Shutdown()

	if err != nil {
		log.Fatal(err.Error())
	}
}
