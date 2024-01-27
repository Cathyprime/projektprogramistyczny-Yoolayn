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
new_user                              # create new user
get_users                             # get all users
add_then_get                          # find user of id
get_bad                               # get not existant user
add_then_update                       # update user
add_then_update_bad                   # update user non existant
add_then_update_noauth                # try to update with wrong credentials
create_delete_get                     # delete user
delete_bad                            # delete non existant
delete_noauth                         # delete without authorization
fill_dummy_users                      # fill db with 5 users
search_user_by_name user1             # find user by the name of name1
search_user_by_bio bio1               # find the user by the bio of bio1
search_user_by_both user1 bio2        # find the users that have either name=name1 or bio=bio2
new_board                             # create a new board
get_boards                            # get all boards
add_then_get_board                    # get board by id
add_then_update_board                 # update a board
add_then_delete                       # delete a board
search_board_by_name "board_example1" # search the board by name
new_post                              # create new post
get_post                              # get post by id
get_posts                             # get all posts in a board
update_post                           # update a post
delete_post                           # delete specified post
search_post "outragous title"         # search for a post with title
'

baseurl="localhost:8080"

# Get all users
function get_users() {
    curl -X GET $baseurl/users 2>/dev/null | jq
}

function bad_email() {
    jq '.user.email |= "badmail.xdddd"'     \
        < ./requests/user_test.json         \
        | curl -X POST $baseurl/users       \
        -H 'Content-Type: application/json' \
        -d '@-'                             \
        2>/dev/null                         \
        | jq
}

# Create new user
function new_user() {
    name=$(bin/generator --strip)
    body=$(jq ".user.name |= . + \"-$name\"" < requests/user_test.json)
    curl -X POST $baseurl/users              \
         -H 'Content-Type: application/json' \
         -d "$body"                          \
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

function add_then_get_board() {
    id=$(                                    \
        new_board 3                          \
        | jq '.id'                           \
        | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    url="$baseurl/boards/$id"

    curl -X GET "$url" \
        2>/dev/null    \
        | jq
}

function add_then_update_board() {
    board=$(add_then_get_board)
    echo "original board:"
    echo "$board" | jq
    echo "updated:"
    owner=$(echo "$board" | jq '.owner')
    update=$(jq                                                                           \
        --argjson board "$board"                                                          \
        '. |= (.id = $board.id | .moderators = $board.moderators | .owner = $board.owner)'\
        < ./requests/update_board.json | tr '\n' ' ')
    full="{\"board\": $update, \"requester\": $owner }"
    echo "$full" | jq
}

function add_then_delete() {
    board=$(add_then_get_board)
    echo "created board:"
    id=$(jq '.id' <<< "$board" | sed 's#"##g')
    echo "$board" | jq

    url="$baseurl/boards/$id"
    body="{\"requester\": $(jq '.owner' <<< "$board")}"

    echo "deleting:"
    curl -X DELETE "$url"                  \
        -H 'Content-Type: application/json'\
        --data-binary "$body"              \
        2>/dev/null                        \
        | jq

    echo "check"
    curl -X GET "$url" 2>/dev/null | jq
}

function new_post() {
    id=$(add_then_get_board | jq '.id')
    author=$(add_then_get | jq '.id')
    filler="{ \"author\": $author, \"board\": $id }"
    url="$baseurl/boards/${id//\"/}/posts"
    body=$(jq --argjson filler "$filler" '.post |= ( .author = $filler.author | .board = $filler.board )' < ./requests/post_test.json)

    curl -X POST "$url"                     \
        -H 'Content-Type: application/json' \
        --data-binary @- <<< "$body"        \
        2>/dev/null                         \
        | jq --argjson id "$id" '.board = $id'
}

function n_new_posts() {
    id=$1
    author=$(add_then_get | jq '.id')
    filler="{\"author\": $author, \"board\": \"$id\"}"
    url="$baseurl/boards/$id/posts"

    body=$(jq --argjson filler "$filler" '.post |= ( .author = $filler.author | .board = $filler.board )' < ./requests/post_test.json)

    for _ in $(seq 1 "$2"); do
        curl -X POST "$url" \
            -H 'Content-Type: application/json' \
            --data-binary @- <<< "$body"        \
            2>/dev/null
    done
}

function get_post() {
    post=$(new_post)
    url="$baseurl/boards/$(jq '.board' <<< "$post" | sed 's#\("\|\\\)##g')/posts/$(jq '.id' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')"

    curl -X GET "$url" 2>/dev/null | jq
}

function update_post() {
    post=$(get_post)
    url="$baseurl/boards/$(jq '.board' <<< "$post" | sed 's#\("\|\\\)##g')/posts/$(jq '.id' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')"

    body=$(jq --argjson post "$post" '.post |= (.id = $post.id | .author = $post.author | .board = $post.board) | .requester = $post.author' < ./requests/updated_post.json)

    curl -X PUT "$url"                      \
        -H 'Content-Type: application/json' \
        --data-binary @- <<< "$body"        \
        2>/dev/null                         \
        | jq
}

function get_posts() {
    post=$(new_post)
    id=$(jq '.board' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g' )
    url="$baseurl/boards/$id"

    n_new_posts "$id" 5

    curl -X GET "$url" 2>/dev/null | jq
}

function delete_post() {
    post=$(get_post)
    echo "new post:"
    echo "$post" | jq
    boardid=$(jq '.board' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')
    postid=$(jq '.id' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    url="$baseurl/boards/$boardid/posts/$postid"

    body="{\"requester\": $(jq '.author' <<< "$post" )}"
    echo "deleting:"
    curl -X DELETE "$url"                  \
        -H 'Content-Type: application/json'\
        --data-binary "$body"              \
        2>/dev/null                        \
        | jq

    echo "check"
    curl -X GET "$url" 2>/dev/null | jq
}

function search_board_by_name() {
    curl -X GET $baseurl/boards/search\?name="$1" \
        2>/dev/null                              \
        | jq
}

function search_post() {
    name=$1
    post=$(get_post)
    id=$(jq '.id' <<< "$post" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')
    post2=$(jq "del(.id) | .title |= \"$name\"" <<< "$post")
    post2=$(jq --argjson post "$post2" '.post |= $post' <<< '{ "post": "" }')
    board=$(jq '.board' <<< "$post2" | sed -E -e 's#ObjectID|"|\(|\)|\\##g')

    url="$baseurl/boards/$board/posts"
    curl -X POST "$url"                    \
        -H 'Content-Type: application/json'\
        --data-binary @- <<< "$post2"      \
        2>/dev/null                        \
        | jq

    url="$baseurl/boards/$board/posts/search?title=$name"
    curl -X GET "$url"\
        2>/dev/null   \
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
