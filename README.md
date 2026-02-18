# Hochschule Burgenland: BSWE - Assignment Infrastructure

[![Build status](https://img.shields.io/github/actions/workflow/status/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pipeline.yml?style=for-the-badge)](https://github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/actions/workflows/pipeline.yml)
[![License](https://img.shields.io/github/license/muhlba91/fh-burgenland-bswe-assignment-infrastructure?style=for-the-badge)](LICENSE.md)
[![](https://api.scorecard.dev/projects/github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/badge?style=for-the-badge)](https://scorecard.dev/viewer/?uri=github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure)

This repository contains the automation for the required assignment infrastructure for the Hochschule Burgenland "Softwaremanagent II" course using [Pulumi](http://pulumi.com).

---

## Requirements

- [Go](https://golang.org/dl/)
- [Pulumi](https://www.pulumi.com/docs/install/)

## Creating the Infrastructure

To create the services, a [Pulumi Stack](https://www.pulumi.com/docs/concepts/stack/) with the correct configuration needs to exists.

The stack can be deployed via:

```bash
pulumi up
```

## Destroying the Infrastructure

The entire infrastructure can be destroyed via:

```bash
pulumi destroy
```

**Attention**: you must set `ALLOW_REPOSITORY_DELETION="true"` as an environment variable to be able to delete repositories!

## Environment Variables

To successfully run, and configure the Pulumi plugins, you need to set a list of environment variables. Alternatively, refer to the used Pulumi provider's configuration documentation.

- `AWS_REGION`: the AWS region to use
- `AWS_ACCESS_KEY_ID`: the AWS secret key
- `AWS_SECRET_ACCESS_KEY`: the AWS secret access key
- `GITHUB_TOKEN`: the GitHub token with permissions to manage repositories
- `PULUMI_ACCESS_TOKEN`: the Pulumi access token
- `VAULT_ADDR`: the Vault address
- `VAULT_TOKEN`: the Vault token
- `HARBOR_URL`: the Harbor URL
- `HARBOR_USERNAME`: the Harbor username
- `HARBOR_PASSWORD`: the Harbor password
- `DEX_HARBOR_CLIENT_SECRET`: the client secret in Dex for Harbor

---

## Configuration

The following section describes the configuration which must be set in the Pulumi Stack.

***Attention:*** do use [Secrets Encryption](https://www.pulumi.com/docs/concepts/secrets/#:~:text=Pulumi%20never%20sends%20authentication%20secrets,“secrets”%20for%20extra%20protection.) provided by Pulumi for secret values!

### YAML Configuration

Repositories and teams are defined in YAML format in the file `assets/data_<stack>.yml` where `<stack>` is the name of the Pulumi stack.

```yaml
---
teams:
  - name: name
    deleteOnDestroy: false # deletes the team on stack destroy (optional, default: false)
    members:
      - user-a
      - user-b
repositories:
  - name: name
    service: service # the service name
    teams: # provides access to the repository
      - name: team-a
        role: developer
      - name: team-b
        role: viewer
    approvers: 1 # the number of required approvers (optional, default: 1)
    deleteOnDestroy: false # deletes the repository on stack destroy (optional, default: false)
    aws: false # creates and stores an AWS access key
    terraform: false # creates a Terraform backend bucket
    harbor: false # creates a Harbor project and robot accounts
    requiredChecks: [] # sets required checks for the repository
```

---

## Outputs

The following outputs are available after the stack has been created:

- `aws`: a map containing the repositories with AWS access and their role ARN for authentication
- `terraform`: a map containing the repositories with Terraform backend bucket information
- `teams`: a list of all created teams
- `repositories`: a list of all created repositories
- `harbor.projects`: a list of all created Harbor projects
- `harbor.robotAccounts`: a list of all created Harbor robot accounts

---

## Continuous Integration and Automations

- [GitHub Actions](https://docs.github.com/en/actions) are linting, and verifying the code.
- [Renovate Bot](https://github.com/renovatebot/renovate) is updating Go modules, and GitHub Actions.
