// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package cmd

import (
	"errors"
	"fmt"

	"github.com/Azure/radius/pkg/cli"
	"github.com/Azure/radius/pkg/cli/bicep"
	"github.com/Azure/radius/pkg/cli/deploy"
	"github.com/Azure/radius/pkg/cli/output"
	"github.com/Azure/radius/pkg/version"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [app.bicep]",
	Short: "Deploy a RAD application",
	Long: `Deploy a RAD application

The deploy command compiles a .bicep file and deploys it to your default environment (unless otherwise specified).
	
You can combine Radius types as as well as other types that are available in Bicep such as Azure resources. See
the Radius documentation for information about describing your application and resources with Bicep.

You can specify parameters using the '--parameter' flag ('-p' for short). Parameters can be passed as:

- A file containing multiple parameters using the ARM JSON parameter format (see below)
- A file containing a single value in JSON format
- A key-value-pair passed in the command line

When passing multiple parameters in a single file, use the format described here:

	https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/parameter-files

You can specify parameters using multiple sources. Parameters can be overridden based on the 
order the are provided. Parameters appearing later in the argument list will override those defined earlier.
`,
	Example: `
# deploy a template (basic)

rad deploy myapp.bicep


# deploy to a specific environment

rad deploy myapp.bicep --environment production


# specify a string parameter

rad deploy myapp.bicep --parameters version=latest


# specify a non-string parameter using a JSON file

rad deploy myapp.bicep --parameters configuration=@myfile.json


# specify many parameters using an ARM JSON parameter file

rad deploy myapp.bicep --parameters @myfile.json


# specify parameters from multiple sources

rad deploy myapp.bicep --parameters @myfile.json --parameters version=latest
`,
	RunE: runDeploy,
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringP("environment", "e", "", "The environment name")
	deployCmd.Flags().StringArrayP("parameters", "p", []string{}, "Specify parameters for the deployment")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New(".bicep file is required")
	}

	filePath := args[0]
	err := deploy.ValidateBicepFile(filePath)
	if err != nil {
		return err
	}

	parameterArgs, err := cmd.Flags().GetStringArray("parameters")
	if err != nil {
		return err
	}

	parser := cli.ParameterParser{FileSystem: cli.OSFileSystem{}}
	parameters, err := parser.Parse(parameterArgs)
	if err != nil {
		return err
	}

	config := ConfigFromContext(cmd.Context())
	env, err := cli.RequireEnvironment(cmd, config)
	if err != nil {
		return err
	}

	ok, err := bicep.IsBicepInstalled()
	if err != nil {
		return fmt.Errorf("failed to find rad-bicep: %w", err)
	}

	if !ok {
		output.LogInfo(fmt.Sprintf("Downloading Bicep for channel %s...", version.Channel()))
		err = bicep.DownloadBicep()
		if err != nil {
			return fmt.Errorf("failed to download rad-bicep: %w", err)
		}
	}

	err = deploy.ValidateBicepFile(filePath)
	if err != nil {
		return err
	}

	step := output.BeginStep("Building %s...", filePath)
	template, err := bicep.Build(filePath)
	if err != nil {
		return err
	}
	output.CompleteStep(step)

	var progressText string
	status := env.GetStatusLink()
	if status == "" {
		progressText = fmt.Sprintf(
			"Deploying Application into environment '%v'...\n\n"+
				"Deployment In Progress...", env.GetName())
	} else {
		progressText = fmt.Sprintf(
			"Deploying Application into environment '%v'...\n\n"+
				"Meanwhile, you can view the environment '%v' at:\n%v\n\n"+
				"Deployment In Progress...", env.GetName(), env.GetName(), status)
	}

	_, err = deploy.DeployWithProgress(cmd.Context(), deploy.Options{
		Environment:    env,
		Template:       template,
		Parameters:     parameters,
		ProgressText:   progressText,
		CompletionText: "Deployment Complete",
	})
	if err != nil {
		return err
	}

	return nil
}