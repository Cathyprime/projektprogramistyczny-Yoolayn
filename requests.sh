: '
# commands
docker compose up -d
make server
./bin/server &
jobs
kill $(jobs -P)
docker compose down
'

: '
# requests
new_post
get_users
new_user
add_then_get
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
