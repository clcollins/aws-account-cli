package cmd

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"
)

// rootCmd represents the base command when called without any subcommands

func init() {
	NewCmdRoot(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
}

func NewCmdRoot(streams genericclioptions.IOStreams) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "aws-account-cli",
		Short: "AWS account cli",
		Long:  `AWS account command line utilities`,
		Run:   help,
	}

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// Reuse kubectl global flags to provide namespace, context and credential options
	kubeFlags := genericclioptions.NewConfigFlags(false)
	kubeFlags.AddFlags(rootCmd.PersistentFlags())

	// add sub commands
	rootCmd.AddCommand(newCmdReset(streams, kubeFlags))

	// add options command to list global flags
	templates.ActsAsRootCommand(rootCmd, []string{"options"})
	rootCmd.AddCommand(newCmdOptions(streams))

	return rootCmd
}

func help(cmd *cobra.Command, _ []string) {
	cmd.Help()
}
