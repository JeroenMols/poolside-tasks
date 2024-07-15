# Backend
> ⚠️ Restarting the server clears all data from the database

## Live reload
Live reloading the server is possible using [Air](https://github.com/air-verse/air).

- `go install github.com/air-verse/air@latest`
- `alias air='$(go env GOPATH)/bin/air'`
- `air`

## Curl
- `curl "http://localhost:8080/debug"`
- `curl -X POST "http://localhost:8080/users/register" -d '{"name":"jeroen"}'`
- `curl -X POST "http://localhost:8080/users/login" -d "{\"user_id\":\"$USER_ID\"}"`
- `curl -X POST "http://localhost:8080/todolists" -d '{}' -H "Authorization: $TOKEN"`
- `curl -X GET "http://localhost:8080/todolists/$LIST" -H "Authorization: $TOKEN"`
- `curl -X POST "http://localhost:8080/todos" -d "{\"todo_list_id\":\"$LIST\", \"description\":\"my first todo\"}" -H "Authorization: $TOKEN"`
- `curl -X PUT "http://localhost:8080/todos/$TODO" -d '{"status":"ongoing"}' -H "Authorization: $TOKEN"` 
