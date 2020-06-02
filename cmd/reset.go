package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	awsv1alpha1 "github.com/openshift/aws-account-operator/pkg/apis/aws/v1alpha1"
	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AWS_ACCOUNT_NAMESPACE = "aws-account-operator"
)

func newCmdReset(streams genericclioptions.IOStreams, flags *genericclioptions.ConfigFlags) *cobra.Command {
	ops := newResetOptions(streams, flags)
	resetCmd := &cobra.Command{
		Use:                   "reset [flags] <account name>",
		Short:                 "",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(ops.complete(cmd, args))
			cmdutil.CheckErr(ops.run())
		},
	}

	return resetCmd
}

// resetCmd represents the export command
type resetOptions struct {
	accountName string

	flags   *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	kubeCli client.Client

	logger *log.Logger
}

func newResetOptions(streams genericclioptions.IOStreams, flags *genericclioptions.ConfigFlags) *resetOptions {
	return &resetOptions{
		flags:     flags,
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

func (o *resetOptions) run() error {
	key := types.NamespacedName{
		Namespace: AWS_ACCOUNT_NAMESPACE,
		Name:      o.accountName,
	}

	ctx := context.TODO()

	//cleanup secrets
	var secrets v1.SecretList
	if err := o.kubeCli.List(ctx, &secrets, &client.ListOptions{
		Namespace: AWS_ACCOUNT_NAMESPACE,
	}); err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		if secret.OwnerReferences != nil &&
			secret.OwnerReferences[0].Name == o.accountName {
			fmt.Println("Deleting secret", secret.Name)
			if err := o.kubeCli.Delete(ctx, &secret, &client.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	var account awsv1alpha1.Account
	if err := o.kubeCli.Get(ctx, key, &account); err != nil {
		return err
	}

	// reset fields in spec
	account.Spec.ClaimLink = ""
	account.Spec.ClaimLinkNamespace = ""
	account.Spec.IAMUserSecret = ""
	if err := o.kubeCli.Update(ctx, &account, &client.UpdateOptions{}); err != nil {
		return err
	}

	// reset fields in status
	var mergePatch []byte
	mergePatch, _ = json.Marshal(map[string]interface{}{
		"status": map[string]interface{}{
			"rotateCredentials":        false,
			"rotateConsoleCredentials": false,
			"claimed":                  false,
			"state":                    "",
		},
	})

	return o.kubeCli.Status().Patch(ctx, &account, client.RawPatch(types.MergePatchType, mergePatch))
}
