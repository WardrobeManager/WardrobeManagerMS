  335  curl http://localhost:8080/user
  336  curl http://localhost:8080/User
  370  history | grep curl
  373  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"xyz","password":"xyz"}' \
  380  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"bengy","password":"xyz", "first-name":"benja", "last-name":"min" }'   http://localhost:8080/User
  381  curl --header "Content-Type: application/json"   --request POST   --data '{"user":"bengy","password":"xyz", "first-name":"benja", "last-name":"min" }'   http://localhost:8080/User
  382  curl --header "Content-Type: application/json"   --request POST   GET   http://localhost:8080/User/bengy
  383  curl --request GET   http://localhost:8080/User/bengy
  399  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"bengy","password":"xyz", "first-name":"benja", "last-name":"min" }'   http://localhost:8080/User
  400   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  404  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"bengy","password":"xyz", "first-name":"benja", "last-name":"min" }'   http://localhost:8080/User
  405   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  406   curl --header "Content-Type: application/json"   --request PUT   --data '{"username":"bengy","password":"xyz", "first-name":"benjamin", "last-name":"button" }'   http://localhost:8080/User
  407   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  408  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:8080/User
  409   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  410   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/jenny
  411   curl --header "Content-Type: application/json" --request DELETE   http://localhost:8080/User/bengy
  412   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/jen
  413   curl --header "Content-Type: application/json" --request DELETE   http://localhost:8080/User/bengy
  414   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  431  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:8080/User
  473  history | grep curl
  475  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:8080/User
  478   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/bengy
  479   curl --header "Content-Type: application/json" --request GET   http://localhost:57400/User/bengy
  480  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  481  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  482  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  483  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  484  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  485  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny",rama":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  486  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"rama","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  487  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"jenny",rama":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  488  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"rama","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  489   curl --header "Content-Type: application/json" --request GET   http://localhost:8080/User/rama
  490   curl --header "Content-Type: application/json" --request GET   http://localhost:57400/User/rama
  491  curl --header "Content-Type: application/json"   --request POST   --data '{"username":"rama","password":"xyz", "first-name":"jenny", "last-name":"mccarthy" }'   http://localhost:57400/User
  535  history | grep curl
  726  history | grep curl
  727  history | grep curl > README.md
