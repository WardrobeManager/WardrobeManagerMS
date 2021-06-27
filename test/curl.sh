#!/bin/bash
for ((i = 0 ; i < $1 ; i++)); do
   uuid=$(uuidgen)
   curl --form "description=Some image" --form "main-image=@/tmp/mypict.jpeg;type=image/jpeg" --form "label-image=@/tmp/mypict.jpeg;type=image/jpeg" http://localhost:57401/wardrobe/foo/$uuid --trace -
done
