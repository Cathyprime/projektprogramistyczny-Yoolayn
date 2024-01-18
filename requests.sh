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
'

# Make a post
function new_post() {
    curl -X POST localhost:8080/posts        \
         -H 'Content-Type: application/json' \
         -d '@./requests/post_test.json'
}

# Get all users
function get_users() {
    curl -X GET localhost:8080/users
}

# Create new user
function new_user() {
    curl -X POST localhost:8080/users        \
         -H 'Content-Type: application/json' \
         -d '@./requests/user_test.json'
}

# --- --- --- ---
function jobs() {
    case "$1" in
        "") builtin jobs ;;
        "-P") builtin jobs -p | sed 's;\[[0-9]\+\]\s\s.\s\([0-9]\+\).*;\1;' ;;
        *) builtin jobs "$@" ;;
    esac
}

function exit_req() {
    kill $(jobs -P)
    docker compose down
    exit 0
}

trap "kill \$(jobs -P)" INT
trap exit_req EXIT
