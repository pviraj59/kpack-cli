// Copyright 2020-2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package custombuilder

import (
	"sort"

	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	expv1alpha1 "github.com/pivotal/kpack/pkg/apis/experimental/v1alpha1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pivotal/build-service-cli/pkg/commands"
	"github.com/pivotal/build-service-cli/pkg/k8s"
)

func NewListCommand(clientSetProvider k8s.ClientSetProvider) *cobra.Command {
	var (
		namespace string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available custom builders",
		Long: `Prints a table of the most important information about the available custom builders in the provided namespace.

namespace defaults to the kubernetes current-context namespace.`,
		Example:      "kp cb list\nkp cb list -n my-namespace",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cs, err := clientSetProvider.GetClientSet(namespace)
			if err != nil {
				return err
			}

			builderList, err := cs.KpackClient.ExperimentalV1alpha1().CustomBuilders(cs.Namespace).List(metav1.ListOptions{})
			if err != nil {
				return err
			}

			if len(builderList.Items) == 0 {
				return errors.New("no builders found")
			} else {
				sort.Slice(builderList.Items, Sort(builderList.Items))
				return displayClusterBuildersTable(cmd, builderList)
			}
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace")

	return cmd
}

func displayClusterBuildersTable(cmd *cobra.Command, builderList *expv1alpha1.CustomBuilderList) error {
	writer, err := commands.NewTableWriter(cmd.OutOrStdout(), "Name", "Ready", "Stack", "Image")
	if err != nil {
		return err
	}

	for _, bldr := range builderList.Items {
		err := writer.AddRow(
			bldr.ObjectMeta.Name,
			getStatus(bldr),
			bldr.Status.Stack.ID,
			bldr.Status.LatestImage,
		)

		if err != nil {
			return err
		}
	}

	return writer.Write()
}

func Sort(builds []expv1alpha1.CustomBuilder) func(i int, j int) bool {
	return func(i, j int) bool {
		return builds[j].ObjectMeta.Name > builds[i].ObjectMeta.Name
	}
}

func getStatus(b expv1alpha1.CustomBuilder) string {
	cond := b.Status.GetCondition(corev1alpha1.ConditionReady)
	switch {
	case cond.IsTrue():
		return "true"
	case cond.IsFalse():
		return "false"
	default:
		return "unknown"
	}
}
