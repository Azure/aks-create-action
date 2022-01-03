FROM mcr.microsoft.com/aks/github-actions/aks-create

COPY . /action

ENTRYPOINT ["/action/entrypoint.sh"]
