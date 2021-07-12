#!/bin/bash
for ((i = 0 ; i < $1 ; i++)); do
   uuid=$(uuidgen)
   curl --form "description=Some image" --form "main-image=@/tmp/mypict.jpeg;type=image/jpeg" --form "label-image=@/tmp/mypict.jpeg;type=image/jpeg" http://127.0.0.1:57401/users/a@a.com/wardrobes --trace -
done
