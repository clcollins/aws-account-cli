package cmd

import (
	"flag"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// rootCmd represents the base command when called without any subcommands

func init() {
	NewCmdRoot(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
}

func NewCmdRoot(streams genericclioptions.IOStreams) *cobra.Command {
	// defaultClusterName := os.Getenv("CLUSTER_NAME")
	// cmd.PersistentFlags().StringVarP(&rootCommand.clusterName, "name", "", defaultClusterName, "Name of cluster. Overrides KOPS_CLUSTER_NAME environment variable")

	rootCmd := &cobra.Command{
		Use:   "aws-account-cli",
		Short: "",
		Long:  ``,
		Run:   help,
	}

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// Reuse kubectl global flags to provide namespace, context and credential options
	kubeFlags := genericclioptions.NewConfigFlags(false)
	kubeFlags.AddFlags(rootCmd.PersistentFlags())

	f := cmdutil.NewFactory(kubeFlags)

	// create subcommands
	rootCmd.AddCommand(newCmdReset(f))

	return rootCmd
}

func help(cmd *cobra.Command, _ []string) {
	cmd.Help()
}
