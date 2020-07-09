// Copyright 2020-2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/pivotal/kpack/pkg/logs"
	"github.com/spf13/cobra"

	"github.com/pivotal/build-service-cli/pkg/buildpackage"
	"github.com/pivotal/build-service-cli/pkg/commands"
	buildcmds "github.com/pivotal/build-service-cli/pkg/commands/build"
	buildercmds "github.com/pivotal/build-service-cli/pkg/commands/custombuilder"
	clusterbuildercmds "github.com/pivotal/build-service-cli/pkg/commands/customclusterbuilder"
	imgcmds "github.com/pivotal/build-service-cli/pkg/commands/image"
	secretcmds "github.com/pivotal/build-service-cli/pkg/commands/secret"
	stackcmds "github.com/pivotal/build-service-cli/pkg/commands/stack"
	storecmds "github.com/pivotal/build-service-cli/pkg/commands/store"
	"github.com/pivotal/build-service-cli/pkg/image"
	"github.com/pivotal/build-service-cli/pkg/k8s"
	"github.com/pivotal/build-service-cli/pkg/secret"
	"github.com/pivotal/build-service-cli/pkg/source"
	"github.com/pivotal/build-service-cli/pkg/stack"
	"github.com/pivotal/build-service-cli/pkg/store"
)

var (
	Version   = "dev"
	CommitSHA = ""
)

func main() {
	var clientSetProvider k8s.DefaultClientSetProvider

	rootCmd := &cobra.Command{
		Use: "kp",
		Long: `kp controls the kpack installation on Kubernetes.

kpack extends Kubernetes and utilizes unprivileged kubernetes primitives to provide 
builds of OCI images as a platform implementation of Cloud Native Buildpacks (CNB).
Learn more about kpack @ https://github.com/pivotal/kpack`,
	}
	rootCmd.AddCommand(
		getVersionCommand(),
		getImageCommand(clientSetProvider),
		getBuildCommand(clientSetProvider),
		getSecretCommand(clientSetProvider),
		getClusterBuilderCommand(clientSetProvider),
		getBuilderCommand(clientSetProvider),
		getStackCommand(clientSetProvider),
		getStoreCommand(clientSetProvider),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display kp version",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), Version+" "+CommitSHA)
		},
	}
	return versionCmd
}

func getImageCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	sourceUploader := &source.Uploader{}

	imageFactory := &image.Factory{
		SourceUploader: sourceUploader,
	}

	imagePatchFactory := &image.PatchFactory{
		SourceUploader: sourceUploader,
	}

	imageRootCmd := &cobra.Command{
		Use:   "image",
		Short: "Image commands",
	}
	newImageWaiter := func(clientSet k8s.ClientSet) imgcmds.ImageWaiter {
		return logs.NewImageWaiter(clientSet.KpackClient, logs.NewBuildLogsClient(clientSet.K8sClient))
	}
	imageRootCmd.AddCommand(
		imgcmds.NewCreateCommand(clientSetProvider, imageFactory, newImageWaiter),
		imgcmds.NewPatchCommand(clientSetProvider, imagePatchFactory, newImageWaiter),
		imgcmds.NewListCommand(clientSetProvider),
		imgcmds.NewDeleteCommand(clientSetProvider),
		imgcmds.NewTriggerCommand(clientSetProvider),
		imgcmds.NewStatusCommand(clientSetProvider),
	)
	return imageRootCmd
}

func getBuildCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	buildRootCmd := &cobra.Command{
		Use:   "build",
		Short: "Build Commands",
	}
	buildRootCmd.AddCommand(
		buildcmds.NewListCommand(clientSetProvider),
		buildcmds.NewStatusCommand(clientSetProvider),
		buildcmds.NewLogsCommand(clientSetProvider),
	)
	return buildRootCmd
}

func getSecretCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	credentialFetcher := &commands.CredentialFetcher{}
	secretFactory := &secret.Factory{
		CredentialFetcher: credentialFetcher,
	}

	secretRootCmd := &cobra.Command{
		Use:   "secret",
		Short: "Secret Commands",
	}
	secretRootCmd.AddCommand(
		secretcmds.NewCreateCommand(clientSetProvider, secretFactory),
		secretcmds.NewDeleteCommand(clientSetProvider),
		secretcmds.NewListCommand(clientSetProvider),
	)
	return secretRootCmd
}

func getClusterBuilderCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	clusterBuilderRootCmd := &cobra.Command{
		Use:     "custom-cluster-builder",
		Short:   "Custom Cluster Builder Commands",
		Aliases: []string{"ccb"},
	}
	clusterBuilderRootCmd.AddCommand(
		clusterbuildercmds.NewCreateCommand(clientSetProvider),
		clusterbuildercmds.NewPatchCommand(clientSetProvider),
		clusterbuildercmds.NewListCommand(clientSetProvider),
		clusterbuildercmds.NewStatusCommand(clientSetProvider),
		clusterbuildercmds.NewDeleteCommand(clientSetProvider),
	)
	return clusterBuilderRootCmd
}

func getBuilderCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	builderRootCmd := &cobra.Command{
		Use:     "custom-builder",
		Short:   "Custom Builder Commands",
		Aliases: []string{"cb"},
	}
	builderRootCmd.AddCommand(
		buildercmds.NewCreateCommand(clientSetProvider),
		buildercmds.NewPatchCommand(clientSetProvider),
		buildercmds.NewListCommand(clientSetProvider),
		buildercmds.NewDeleteCommand(clientSetProvider),
		buildercmds.NewStatusCommand(clientSetProvider),
	)
	return builderRootCmd
}

func getStackCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	stackFactory := &stack.Factory{
		Fetcher:   &image.Fetcher{},
		Relocator: &image.Relocator{},
	}

	stackRootCmd := &cobra.Command{
		Use:   "stack",
		Short: "Stack Commands",
	}
	stackRootCmd.AddCommand(
		stackcmds.NewCreateCommand(clientSetProvider, stackFactory),
		stackcmds.NewListCommand(clientSetProvider),
		stackcmds.NewStatusCommand(clientSetProvider),
		stackcmds.NewUpdateCommand(clientSetProvider, &image.Fetcher{}, &image.Relocator{}),
		stackcmds.NewDeleteCommand(clientSetProvider),
	)
	return stackRootCmd
}

func getStoreCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	bpUploader := &buildpackage.Uploader{
		Fetcher:   &image.Fetcher{},
		Relocator: &image.Relocator{},
	}

	factory := &store.Factory{Uploader: bpUploader}

	storeRootCommand := &cobra.Command{
		Use:   "store",
		Short: "Store Commands",
	}
	storeRootCommand.AddCommand(
		storecmds.NewCreateCommand(clientSetProvider, factory),
		storecmds.NewDeleteCommand(clientSetProvider, commands.NewConfirmationProvider()),
		storecmds.NewAddCommand(clientSetProvider, factory),
		storecmds.NewStatusCommand(clientSetProvider),
		storecmds.NewRemoveCommand(clientSetProvider),
		storecmds.NewListCommand(clientSetProvider),
	)
	return storeRootCommand
}
