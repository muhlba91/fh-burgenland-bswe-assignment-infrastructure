#!/bin/bash

YEAR_IDENTIFIER="ws2024"
IDENTIFIERS="a b c d e f g h i j k l"
MATCHED_GROUPS="aj bc di ef gl hk"

IFS=' ' read -r -a array <<< "$IDENTIFIERS"
IFS=' ' read -r -a array_groups <<< "$MATCHED_GROUPS"

echo ""
echo "----------------------------------------"
echo "accounts"
echo "----------------------------------------"
echo ""

for identifier in "${array[@]}"; do
    echo "[$identifier][accounts] checkout"
    git clone git@github.com:fhburgenland-bswe/swm2-$YEAR_IDENTIFIER-group-$identifier-accounts.git

    cd swm2-$YEAR_IDENTIFIER-group-$identifier-accounts

    echo "[$identifier][accounts] create README.md"
    echo "swm2-$YEAR_IDENTIFIER-group-$identifier-accounts" > README.md

    echo "[$identifier][accounts] push to remote"
    git add README.md
    git commit -am "feat: initial commit"
    git push

    echo "[$identifier][accounts] finalize"
    cd ..
    echo ""
done

echo ""
echo "----------------------------------------"
echo "transactions"
echo "----------------------------------------"
echo ""

for identifier in "${array[@]}"; do
    echo "[$identifier][transactions] checkout"
    git clone git@github.com:fhburgenland-bswe/swm2-$YEAR_IDENTIFIER-group-$identifier-transactions.git

    cd swm2-$YEAR_IDENTIFIER-group-$identifier-transactions

    echo "[$identifier][transactions] create README.md"
    echo "swm2-$YEAR_IDENTIFIER-group-$identifier-transactions" > README.md

    echo "[$identifier][transactions] push to remote"
    git add README.md
    git commit -am "feat: initial commit"
    git push

    echo "[$identifier][transactions] finalize"
    cd ..
    echo ""
done

echo ""
echo "----------------------------------------"
echo "broker"
echo "----------------------------------------"
echo ""

for identifier in "${array_groups[@]}"; do
    echo "[$identifier][broker] checkout"
    git clone git@github.com:fhburgenland-bswe/swm2-$YEAR_IDENTIFIER-group-$identifier-broker.git

    cd swm2-$YEAR_IDENTIFIER-group-$identifier-broker

    echo "[$identifier][broker] create README.md"
    echo "swm2-$YEAR_IDENTIFIER-group-$identifier-broker" > README.md

    echo "[$identifier][broker] push to remote"
    git add README.md
    git commit -am "feat: initial commit"
    git push

    echo "[$identifier][broker] finalize"
    cd ..
    echo ""
done

echo ""
echo "----------------------------------------"
echo "frontend"
echo "----------------------------------------"
echo ""

for identifier in "${array_groups[@]}"; do
    echo "[$identifier][frontend] checkout"
    git clone git@github.com:fhburgenland-bswe/swm2-$YEAR_IDENTIFIER-group-$identifier-frontend.git

    cd swm2-$YEAR_IDENTIFIER-group-$identifier-frontend

    echo "[$identifier][frontend] create README.md"
    echo "swm2-$YEAR_IDENTIFIER-group-$identifier-frontend" > README.md

    echo "[$identifier][frontend] push to remote"
    git add README.md
    git commit -am "feat: initial commit"
    git push

    echo "[$identifier][frontend] finalize"
    cd ..
    echo ""
done
