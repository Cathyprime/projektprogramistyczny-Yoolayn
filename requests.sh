: '
# commands
docker compose up -d
make server
./bin/server &
jobs
kill $(jobs -P)
docker compose down
source requests.sh
'

: '
# requests
new_user               # create new user
get_users              # get all users
add_then_get           # find user of id
get_bad                # get not existant user
add_then_update        # update user
add_then_update_bad    # update user non existant
add_then_update_noauth # try to update with wrong credentials
create_delete_get      # delete user
delete_bad             # delete non existant
delete_noauth          # delete without authorization
'

# Make a post
function new_post() {
    curl -X POST localhost:8080/posts        \
         -H 'Content-Type: application/json' \
         -d '@./requests/post_test.json'     \
         2>/dev/null                         \
    | jq
}

# Get all users
function get_users() {
    curl -X GET localhost:8080/users 2>/dev/null | jq
}

# Create new user
function new_user() {
    curl -X POST localhost:8080/users        \
         -H 'Content-Type: application/json' \
         -d '@./requests/user_test.json'     \
         2>/dev/null                         \
    | jq
}

function add_then_get() {
    id=$(                                \
    new_user                             \
    | jq '.id'                           \
    | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    url="localhost:8080/users/$id"

    curl -X GET "$url" \
    2>/dev/null        \
    | jq
}

function add_then_update() {
    local id
    id=$(                                \
    new_user                             \
    | jq '.id'                           \
    | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    echo "id created: $id"

    local url="localhost:8080/users/$id"

    echo "created:"
    curl -X GET "$url" 2>/dev/null | jq

    local json
    json=$(jq -s '.[0] * .[1]' requests/update_user.json <(echo "{\"requester\": {\"id\": \"$id\"}}"))

    echo "update response: "
    curl -X PUT "$url"                      \
        -H 'Content-Type: application/json' \
        -d "$json"                          \
        2>/dev/null                         \
    | jq

    echo "updated to:"
    curl -X GET "$url" 2>/dev/null | jq
}

function add_then_update_bad() {
    local id="65ad3d81421791df197f89eb"
    local url="localhost:8080/users/$id"
    local json
    json=$(jq -s '.[0] * .[1]' requests/update_user.json <(echo "{\"requester\": {\"id\": \"$id\"}}"))

    curl -X PUT "$url"                      \
        -H 'Content-Type: application/json' \
        -d "$json"                          \
        2>/dev/null                         \
    | jq
}

function add_then_update_noauth() {
    local id="65ad3d81421791df197f89eb"
    local url="localhost:8080/users/65ad7cfd7a8fbdc2f7829aa6"
    local json
    json=$(jq -s '.[0] * .[1]' requests/update_user.json <(echo "{\"requester\": {\"id\": \"$id\"}}"))

    curl -X PUT "$url"                      \
        -H 'Content-Type: application/json' \
        -d "$json"                          \
        2>/dev/null                         \
    | jq
}

function get_bad() {
    curl -X GET localhost:8080/users/65ad3d81421791df197f89eb \
    2>/dev/null                                               \
    | jq
}

function create_delete_get() {
    local id
    id=$(                                \
    new_user                             \
    | jq '.id'                           \
    | sed -E -e 's#ObjectID|"|\(|\)|\\##g')
    
    echo "created user of id: $id"

    local url="localhost:8080/users/$id"

    echo "deleted:"
    curl -X DELETE "$url"                                   \
        -H 'Content-Type: application/json'                 \
        -d "{\"requester\": {\"id\": \"$id\"}}" 2>/dev/null \
    | jq

    echo "checking if deleted:"
    curl -X GET "$url" 2>/dev/null | jq
}

function delete_bad() {
    curl -X DELETE localhost:8080/users/65ad3d81421791df197f89eb    \
        -H 'Content-Type: application/json'                         \
        -d "{\"requester\": {\"id\": \"65ad3d81421791df197f89eb\"}}"\
    2>/dev/null                                                     \
    | jq
}

function delete_noauth() {
    curl -X DELETE localhost:8080/users/65ad3d81421791df197f89eb    \
        -H 'Content-Type: application/json'                         \
        -d "{\"requester\": {\"id\": \"65ad7cfd7a8fbdc2f7829aa6\"}}"\
    2>/dev/null                                                     \
    | jq
}

function jobs() {
    case "$1" in
        "") builtin jobs ;;
        "-P") builtin jobs -p | sed 's;\[[0-9]\+\]\s\s.\s\([0-9]\+\).*;\1;' ;;
        *) builtin jobs "$@" ;;
    esac
}

function exit_req() {
    kill "$(jobs -P)"
    docker compose down
    exit 0
}

trap "kill \$(jobs -P)" INT
# trap exit_req EXIT
