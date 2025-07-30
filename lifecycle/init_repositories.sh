#!/bin/bash

YEAR_IDENTIFIER="ws2025"
IDENTIFIERS="a b c d e f g h i j k l"
MATCHED_GROUPS="aj bc di ef gl hk"

SERVICES="accounts transactions"
GROUP_SERVICES="broker frontend"

IFS=' ' read -r -a array <<< "$IDENTIFIERS"
IFS=' ' read -r -a array_groups <<< "$MATCHED_GROUPS"
IFS=' ' read -r -a array_services <<< "$SERVICES"
IFS=' ' read -r -a array_group_services <<< "$GROUP_SERVICES"

# function to initialize repository based on the identifier and service
initialize_repository() {
    repo_name="swm2-$YEAR_IDENTIFIER-group-$1-$2"

    echo "[$identifier][$service] checkout"
    git clone git@github.com:fhburgenland-bswe/$repo_name.git

    cd $repo_name

    echo "[$identifier][$service] create README.md"
    echo $repo_name > README.md

    echo "[$identifier][$service] push to remote"
    git add README.md
    git commit -am "feat: initial commit"
    git push

    echo "[$identifier][$service] finalize"
    cd ..
    echo ""
}

# initialize repositories for each individual group and service
for service in "${array_services[@]}"; do
  echo ""
  echo "----------------------------------------"
  echo "$service"
  echo "----------------------------------------"
  echo ""

  for identifier in "${array[@]}"; do
    echo "[$identifier][$service] initializing repository..."
    initialize_repository $identifier $service
  done
done

# initialize repositories for each matched group and service
for service in "${array_group_services[@]}"; do
  echo ""
  echo "----------------------------------------"
  echo "$service"
  echo "----------------------------------------"
  echo ""

  for identifier in "${array_groups[@]}"; do
    echo "[$identifier][$service] initializing repository..."
    initialize_repository $identifier $service
  done
done
