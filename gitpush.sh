#!/bin/bash

git add .
git commit -m "cspmWorker: host&port"
git push origin main



az ad sp create-for-rbac --name "github-actions-deploy" --role contributor --scopes /subscriptions/<subscription-id>/resourceGroups/<resource-group> --sdk-auth
az ad sp create-for-rbac --name "github-actions-deploy" --role contributor --scopes /subscriptions/6c47a347-6277-47ea-bc05-aa9f2afc5ccd/resourceGroups/rg-soldout --sdk-auth