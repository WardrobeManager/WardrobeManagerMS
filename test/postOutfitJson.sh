#!/bin/bash
for ((i = 0 ; i < $1 ; i++)); do
    curl --header "Content-Type: application/json"   --request POST   -d @$2  http://localhost:57401/users/a@a.com/outfits
done
