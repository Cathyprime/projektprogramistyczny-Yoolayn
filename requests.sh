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
new_user                       # create new user
get_users                      # get all users
add_then_get                   # find user of id
get_bad                        # get not existant user
add_then_update                # update user
add_then_update_bad            # update user non existant
add_then_update_noauth         # try to update with wrong credentials
create_delete_get              # delete user
delete_bad                     # delete non existant
delete_noauth                  # delete without authorization
fill_dummy_users               # fill db with 5 users
search_user_by_name user1      # find user by the name of name1
search_user_by_bio bio1        # find the user by the bio of bio1
search_user_by_both user1 bio2 # find the users that have either name=name1 or bio=bio2
new_board                      # create a new board
'

baseurl="localhost:8080"

# Make a post
function new_post() {
    curl -X POST $baseurl/posts              \
         -H 'Content-Type: application/json' \
         -d '@./requests/post_test.json'     \
         2>/dev/null                         \
         | jq
}

# Get all users
function get_users() {
    curl -X GET $baseurl/users 2>/dev/null | jq
}

# Create new user
function new_user() {
    curl -X POST $baseurl/users              \
         -H 'Content-Type: application/json' \
         -d '@./requests/user_test.json'     \
         2>/dev/null                         \
         | jq
}

function add_then_get() {
    id=$(                                    \
        new_user                             \
        | jq '.id'                           \
        | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    url="$baseurl/users/$id"

    curl -X GET "$url" \
        2>/dev/null    \
        | jq
}

function add_then_update() {
    local id
    id=$(                                    \
        new_user                             \
        | jq '.id'                           \
        | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    echo "id created: $id"

    local url="$baseurl/users/$id"

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
    local url="$baseurl/users/$id"
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
    local url="$baseurl/users/65ad7cfd7a8fbdc2f7829aa6"
    local json
    json=$(jq                                         \
        -s                                            \
        '.[0] * .[1]'                                 \
        requests/update_user.json                     \
        <(echo "{\"requester\": {\"id\": \"$id\"}}")) \

    curl -X PUT "$url"                      \
        -H 'Content-Type: application/json' \
        -d "$json"                          \
        2>/dev/null                         \
        | jq
}

function get_bad() {
    curl -X GET $baseurl/users/65ad3d81421791df197f89eb \
        2>/dev/null                                     \
        | jq
}

function create_delete_get() {
    local id
    id=$(                                    \
        new_user                             \
        | jq '.id'                           \
        | sed -E -e 's#ObjectID|"|\(|\)|\\##g')
    
    echo "created user of id: $id"

    local url="$baseurl/users/$id"

    echo "deleted:"
    curl -X DELETE "$url"                                   \
        -H 'Content-Type: application/json'                 \
        -d "{\"requester\": {\"id\": \"$id\"}}" 2>/dev/null \
        | jq

    echo "checking if deleted:"
    curl -X GET "$url" 2>/dev/null | jq
}

function delete_bad() {
    curl -X DELETE $baseurl/users/65ad3d81421791df197f89eb           \
        -H 'Content-Type: application/json'                          \
        -d "{\"requester\": {\"id\": \"65ad3d81421791df197f89eb\"}}" \
        2>/dev/null                                                  \
        | jq
}

function delete_noauth() {
    curl -X DELETE $baseurl/users/65ad3d81421791df197f89eb           \
        -H 'Content-Type: application/json'                          \
        -d "{\"requester\": {\"id\": \"65ad7cfd7a8fbdc2f7829aa6\"}}" \
        2>/dev/null                                                  \
        | jq
}

function search_user_by_name() {
    curl -X GET $baseurl/users/search\?name="$1" \
        2>/dev/null                              \
        | jq
}

function search_user_by_bio() {
    curl -X GET $baseurl/users/search\?bio="$1" \
        2>/dev/null                             \
        | jq
}

function search_user_by_both() {
    curl -X GET $baseurl/users/search\?name="$1"\&bio="$2" \
        2>/dev/null                                        \
        | jq
}

function get_boards() {
    curl -X GET $baseurl/boards \
        2>/dev/null             \
        | jq
}

function new_board() {
    ids=()

    for _ in $(seq 1 "${1:-5}"); do
        id=$(add_then_get | jq '.id')
        ids+=("$id")
    done

    json=$(echo "{ \"moderators\": [$(echo "${ids[@]:1}" | tr ' ' ',')], \"owner\": ${ids[1]} }" | jq)
    rdy=$(jq --argjson json "$json" '.board |= (.moderators = $json.moderators | .owner = $json.owner)' < ./requests/board_test.json)

    echo "$rdy" | jq
    echo "$rdy"                                 \
        | curl -X POST "$baseurl/boards"        \
            -H 'Content-Type: application/json' \
            --data-binary @-                    \
            2>/dev/null                         \
        | jq
}

function fill_dummy_users() {
    files=$(ls ./requests/test_user?.json)
    while IFS= read -r file; do
        curl -X POST "$baseurl/users"           \
            -H 'Content-Type: application/json' \
            --data-binary @-                    \
            2>/dev/null < "$file"               \
            | jq
    done <<< "$files"
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
