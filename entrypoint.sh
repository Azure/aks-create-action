#!/bin/sh -l

# uncomment when image is in MCR . /root/.bashrc

export ARM_CLIENT_ID=$INPUT_ARM_CLIENT_ID
export ARM_CLIENT_SECRET=$INPUT_ARM_CLIENT_SECRET
export ARM_SUBSCRIPTION_ID=$INPUT_ARM_SUBSCRIPTION_ID
export ARM_TENANT_ID=$INPUT_ARM_TENANT_ID
export STORAGE_ACCOUNT_NAME=$INPUT_STORAGE_ACCOUNT_NAME
export STORAGE_CONTAINER_NAME=$INPUT_STORAGE_CONTAINER_NAME
export STORAGE_ACCESS_KEY=$INPUT_STORAGE_ACCESS_KEY
export USE_PULUMI=$INPUT_USE_PULUMI

export TF_VAR_resource_group_name=$INPUT_RESOURCE_GROUP_NAME
export TF_VAR_cluster_name=$INPUT_CLUSTER_NAME
export TF_VAR_create_acr=$INPUT_CREATE_ACR
export TF_IN_AUTOMATION=true

## Use TF based on cluster size variable
cd /action/$INPUT_CLUSTER_SIZE

echo "*******************"
echo "Running init"
echo "*******************"


# Using Pulumi
if [ $INPUT_USE_PULUMI = "true" ]; then
 echo "Using Pulumi"
 pulumi stack select dev --create
 pulumi config set azure:clientId ${ARM_CLIENT_ID}
 pulumi config set azure:clientSecret ${ARM_CLIENT_SECRET} --secret
 pulumi config set azure:tenantId ${ARM_TENANT_ID}
 pulumi config set azure:subscriptionId ${ARM_SUBSCRIPTION_ID}

  if [ $INPUT_ACTION_TYPE = "destroy" ]; then
      echo "*******************"
      echo "Running destroy"
      echo "*******************"
      pulumi destroy -force
  else
      echo "*******************"
      echo "Running apply"
      echo "*******************"
      pulumi up --yes
  fi

  exit 1
fi

# Using Terraform -  Default
echo "Using Terraform"
terraform init -backend-config="resource_group_name=${TF_VAR_resource_group_name}" \
-backend-config="storage_account_name=${STORAGE_ACCOUNT_NAME}" \
-backend-config="container_name=${STORAGE_CONTAINER_NAME}" \
-backend-config="key=${TF_VAR_cluster_name}.tfstate" \
-backend-config="access_key=${STORAGE_ACCESS_KEY}"

if [ $INPUT_ACTION_TYPE = "destroy" ]; then
    echo "*******************"
    echo "Running destroy"
    echo "*******************"
    terraform destroy -force
else
    echo "*******************"
    echo "Running plan"
    echo "*******************"

    terraform plan -no-color -out=tfplan -input=false

    echo "*******************"
    echo "Running apply"
    echo "*******************"

    terraform apply tfplan
fi



