package main
// This file is only ran when users enable USE_PULUMI
import (
	"fmt"
	"github.com/pulumi/pulumi-azure/sdk/v4/go/azure/containerservice"
	"github.com/pulumi/pulumi-azure/sdk/v4/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

const STORAGE_ACCOUNT_NAME = "akscreatesa"
const CLUSTER_NAME = "akscreatecluster"
const RESOURCE_GROUP_NAME = "akscreaterg"
const CONTAINER_NAME = "akscreatecontainer"


func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		location := getLocation()
		resourceGroup := getResourceGroup()

		// Create storage account
		storageAccount, err := createStorageAccount(ctx, pulumi.String(location), pulumi.String(resourceGroup))
		if err != nil {
			fmt.Print("error happened during storage account creation")
			return err
		}

		// Create storage container
		err = createStorageContainer(ctx, storageAccount.Name)
		if err != nil {
			fmt.Print("error happened during storage container creation")
			return err
		}


		// Create K8's Cluster
		k8sCluster, err := containerservice.NewKubernetesCluster(ctx, getClusterName(), &containerservice.KubernetesClusterArgs{
			Location:          pulumi.String(location),
			ResourceGroupName: pulumi.String(resourceGroup),
			DnsPrefix:         pulumi.String(getDnsPrefix()),
			DefaultNodePool: &containerservice.KubernetesClusterDefaultNodePoolArgs{
				Name:      pulumi.String("default"),
				NodeCount: pulumi.Int(1),
				VmSize:    pulumi.String("Standard_D2_v2"),
			},
			Identity: &containerservice.KubernetesClusterIdentityArgs{
				Type: pulumi.String("SystemAssigned"),
			},
			Tags: pulumi.StringMap{
				"Environment": pulumi.String("Dev"),
			},
		})
		if err != nil {
			return err
		}

		// Create ACR if needed
		if isCreateACR() {
			err = createACR(ctx, pulumi.String(location), pulumi.String(resourceGroup), k8sCluster)
			if err != nil {
				fmt.Print("error happened during ACR creation")
				return err
			}
		}



		ctx.Export("clientCertificate", k8sCluster.KubeConfigs.ApplyT(func(kubeConfigs []containerservice.KubernetesClusterKubeConfig) (*string, error) {
			return kubeConfigs[0].ClientCertificate, nil
		}).(pulumi.StringPtrOutput))
		ctx.Export("kubeConfig", k8sCluster.KubeConfigRaw)
		return nil
	})
}



// Helper functions

func getClusterName() string {
	clusterName := os.Getenv("CLUSTER_NAME")

	if clusterName == "" {
		return CLUSTER_NAME
	}
	return clusterName
}

func getLocation() string {
	location := os.Getenv("REGION")

	if location == "" {
		return "East US"
	}
	return location
}

func isCreateACR() bool {
	isACR := os.Getenv("CREATE_ACR")
	return isACR == "true"
}


func getResourceGroup() string {
	resourceGroup := os.Getenv("RESOURCE_GROUP_NAME")

	if resourceGroup == "" {
		return RESOURCE_GROUP_NAME
	}
	return resourceGroup
}

func getDnsPrefix() string {
	return "akscreate"
}

func createStorageAccount(ctx *pulumi.Context, location pulumi.StringInput, resourceGroup pulumi.StringInput ) (*storage.Account, error) {
	account, err := storage.NewAccount(ctx, STORAGE_ACCOUNT_NAME, &storage.AccountArgs{
		ResourceGroupName:      resourceGroup,
		Location:               location,
		AccountTier:            pulumi.String("Standard"),
		AccountReplicationType: pulumi.String("GRS"),
		Tags: pulumi.StringMap{
			"environment": pulumi.String("Dev"),
		},
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}


func createStorageContainer(ctx *pulumi.Context, accountName pulumi.StringInput) error {
	_, err := storage.NewContainer(ctx, CONTAINER_NAME, &storage.ContainerArgs{
		StorageAccountName:  accountName,
		ContainerAccessType: pulumi.String("private"),
	})
	if err != nil {
		return err
	}
	return nil
}


func createACR(ctx *pulumi.Context, location pulumi.StringInput, resourceGroup pulumi.StringInput, cluster *containerservice.KubernetesCluster) error {
	// Attach the principle id to the newly create ACR below
	registryIdentityArgs := containerservice.RegistryIdentityArgs{
		PrincipalId: cluster.Identity.PrincipalId(),
		TenantId:    cluster.Identity.TenantId(),
	}

	_, err := containerservice.NewRegistry(ctx, "akscreateacr", &containerservice.RegistryArgs{
		ResourceGroupName: resourceGroup,
		Location:          location,
		Sku:               pulumi.String("Premium"),
		AdminEnabled:      pulumi.Bool(false),
		Identity: registryIdentityArgs,
		Georeplications: containerservice.RegistryGeoreplicationArray{
			&containerservice.RegistryGeoreplicationArgs{
				Location:              pulumi.String("westeurope"),
				ZoneRedundancyEnabled: pulumi.Bool(true),
				Tags:                  nil,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

