FROM gambtho/azurecli_terraform:latest

COPY . /action

ENTRYPOINT ["/action/entrypoint.sh"]
