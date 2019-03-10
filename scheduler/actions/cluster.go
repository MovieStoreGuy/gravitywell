package actions

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider"
	awsprovider "github.com/AlexsJones/gravitywell/platform/provider/aws"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
)

func NewAmazonWebServicesConfig() (*awsprovider.AWSProvider,error){
	awsp := awsprovider.AWSProvider{}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "test-account"),
	})
	if err != nil {
		return nil,err
	}

	awsp.AWSClient = sess

	return &awsp,err
}
func AmazonWebServicesClusterProcessor(awsprovider *awsprovider.AWSProvider,
	commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) error {

	ctx := context.Background()
	cmc, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	create := func() {

		err := provider.Create(awsprovider,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName,
			cluster.Zones,
			int32(cluster.InitialNodeCount),
			cluster.InitialNodeType,
			cluster.Labels, cluster.NodePools)

		if err != nil {
			color.Red(err.Error())
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}

			}
		}
	}
	delete := func(){
		err := provider.Delete(awsprovider,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHooak {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}
			}
		}
	}
	switch commandFlag {
	case configuration.Create:
		create()
	case configuration.Apply:
		create()
	case configuration.Replace:
		delete()
		create()
	case configuration.Delete:
		delete()
	}
	return nil
}

func GoogleCloudPlatformClusterProcessor(commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) error {

	gcpProviderClient := &gcp.GCPProvider{}

	ctx := context.Background()
	cmc, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	create := func() {

		err := provider.Create(gcpProviderClient,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName,
			cluster.Zones,
			int32(cluster.InitialNodeCount),
			cluster.InitialNodeType,
			cluster.Labels, cluster.NodePools)

		if err != nil {
			color.Red(err.Error())
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}

			}
		}
	}
	delete := func() {
		err := provider.Delete(gcpProviderClient,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHooak {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}
			}
		}

	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		create()
	case configuration.Apply:
		create()
	case configuration.Replace:
		delete()
		create()
	case configuration.Delete:
		delete()
	}
	return nil
}
