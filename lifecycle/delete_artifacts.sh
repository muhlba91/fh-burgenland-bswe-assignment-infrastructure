#!/bin/bash

YEAR_IDENTIFIER="ws2024"
IDENTIFIERS="a b c d e f g h i j k l"
MATCHED_GROUPS="aj bc di ef gl hk"

SERVICES="accounts transactions"
GROUP_SERVICES="broker frontend"

IFS=' ' read -r -a array <<< "$IDENTIFIERS"
IFS=' ' read -r -a array_groups <<< "$MATCHED_GROUPS"
IFS=' ' read -r -a array_services <<< "$SERVICES"
IFS=' ' read -r -a array_group_services <<< "$GROUP_SERVICES"

# function to delete artifacts based on the identifier and service
delete_artifacts() {
  repo_name="swm2-$YEAR_IDENTIFIER-group-$1-$2"

  page=1
  while true; do
    artifact_exists=$(gh api repos/fhburgenland-bswe/$repo_name/actions/artifacts?per_page=100\&page=$page | jq -r '.artifacts[]')
    artifacts=$(gh api repos/fhburgenland-bswe/$repo_name/actions/artifacts?per_page=100\&page=$page | jq -r '.artifacts[].id')
    if [[ -z "$artifact_exists" ]]; then
      echo "[$identifier][$service] no more artifacts found"
      break
    fi

    for artifact_id in $artifacts; do
      artifact_name=$(gh api repos/fhburgenland-bswe/$repo_name/actions/artifacts/$artifact_id | jq -r '.name')
      echo "[$identifier][$service] deleting artifact $artifact_name ($artifact_id)..."
      gh api repos/fhburgenland-bswe/$repo_name/actions/artifacts/$artifact_id -X DELETE >& /dev/null
    done

    page=$((page+1))
  done

  echo "[$identifier][$service] artifacts deleted"
  echo ""
}

# delete artifacts for each individual group and service
for service in "${array_services[@]}"; do
  echo ""
  echo "----------------------------------------"
  echo "$service"
  echo "----------------------------------------"
  echo ""

  for identifier in "${array[@]}"; do
    echo "[$identifier][$service] deleting artifacts..."
    delete_artifacts $identifier $service
  done
done

# delete artifacts for each matched group and service
for service in "${array_group_services[@]}"; do
  echo ""
  echo "----------------------------------------"
  echo "$service"
  echo "----------------------------------------"
  echo ""

  for identifier in "${array_groups[@]}"; do
    echo "[$identifier][$service] checking for artifacts..."
    delete_artifacts $identifier $service
  done
done
