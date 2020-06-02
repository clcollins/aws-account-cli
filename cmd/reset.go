package cmd

import (
	"log"

	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// resetCmd represents the export command

type resetOptions struct {
	accountName string

	flags   *genericclioptions.ConfigFlags
	kubeCli client.Client

	genericclioptions.IOStreams
	logger *log.Logger
}

func newResetOptions(streams genericclioptions.IOStreams) *resetOptions {
	return &resetOptions{
		flags:     genericclioptions.NewConfigFlags(false),
		IOStreams: streams,
		logger:    log.New(streams.Out, "aws-account-cli", log.LstdFlags|log.Lshortfile),
	}
}

func (o *resetOptions) complete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmdutil.UsageErrorf(cmd, "")
	}
	o.accountName = args[0]

	var err error
	configLoader := o.flags.ToRawKubeConfigLoader()
	cfg, err := configLoader.ClientConfig()
	if err != nil {
		return err
	}

	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	o.kubeCli = cli
	return nil
}

func (o *resetOptions) run() {

}

func newCmdReset(streams genericclioptions.IOStreams) *cobra.Command {
	ops := newResetOptions(streams)
	resetCmd := &cobra.Command{
		Use:                   "reset [flags] <account name>",
		Short:                 "",
		Version:               version.Version,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(ops.complete(cmd, args))
			cmdutil.CheckErr(op)
		},
	}

	//resetCmd.PersistentFlags().String("account-name", "", "the name of the aws account cr")

	return resetCmd
}
