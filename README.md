# AKS Cluster Creation action

This action creates an Azure Kubernetes Service Cluster using Terraform

## Setup

Making use of this action requires an Azure Service Principal and a resource group containing a storage account to store the terraform state.

These can be created using the setup.sh script in this repo

```
./setup.sh -c <<cluster name> -g <<resource group name>> -s <<subscription id>> -r <<region>>
```

The output from this command should look like this and matches the variables that need to be passed to the action

```
CLUSTER_NAME: testCluster
RESOURCE_GROUP_NAME: newGroup
STORAGE_ACCOUNT_NAME: newgroup27941
STORAGE_CONTAINER_NAME: testclustertstate
STORAGE_ACCESS_KEY: ******
ARM_CLIENT_ID: ******
ARM_CLIENT_SECRET: ******
ARM_SUBSCRIPTION_ID: ******
ARM_TENANT_ID: ******
```


## Inputs

* `CLUSTER_NAME` ***required***
* `RESOURCE_GROUP_NAME` ***required***
* `STORAGE_ACCOUNT_NAME` ***required***
* `STORAGE_CONTAINER_NAME` ***required***
* `STORAGE_ACCESS_KEY` ***required***
* `ARM_CLIENT_ID` ***required***
* `ARM_CLIENT_SECRET` ***required***
* `ARM_SUBSCRIPTION_ID` ***required***
* `ARM_TENANT_ID` ***required***
* `CLUSTER_SIZE` ***optional*** - dev (default) or test
* `ACTION_TYPE` ***optional*** - create (default) or delete 
* `CREATE_ACR` ***optional*** - true or false (default)


## Example usage
```
uses: actions/aks_create_action@v1
with:
  CLUSTER_NAME: testCluster
  RESOURCE_GROUP_NAME: newGroup
  STORAGE_ACCOUNT_NAME: newgroup27941
  STORAGE_CONTAINER_NAME: testclustertstate
  STORAGE_ACCESS_KEY: ******
  ARM_CLIENT_ID: ******
  ARM_CLIENT_SECRET: ******
  ARM_SUBSCRIPTION_ID: ******
  ARM_TENANT_ID: ******
  ACTION_TYPE: create # optional
  CLUSTER_SIZE: dev # optional
  CREATE_ACR: false # optional
```

Full deployment workflow showing this action in use - https://github.com/gambtho/go_echo

## References

* https://docs.microsoft.com/en-us/azure/aks/kubernetes-action
* https://wahlnetwork.com/2020/05/12/continuous-integration-with-github-actions-and-terraform/
* https://github.com/Azure/actions-workflow-samples/tree/master/Kubernetes

## Contributing
This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the Microsoft Open Source Code of Conduct. For more information see the Code of Conduct FAQ or contact opencode@microsoft.com with any additional questions or comments.

## Trademarks
This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow Microsoft's Trademark & Brand Guidelines. Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.