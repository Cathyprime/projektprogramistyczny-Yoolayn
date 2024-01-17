# Make a post
curl -X POST localhost:8080/post \
    -H 'Content-Type: application/json' \
    -d '@./post_test.json'

# Get all users
curl -X GET localhost:8080/users
