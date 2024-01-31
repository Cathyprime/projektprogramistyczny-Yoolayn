POST http://localhost:8080/boards
{
    "board": {
        "id": "65b95156097680ef41e8f930",
        "name": "science",
        "bio": "board about science",
        "moderators": [
            "65b9521f08488450adcbd92d",
            "65b9521f08488450adcbd92e",
            "65b9521f08488450adcbd92f"
        ],
        "owner": "65b954c547c4f420dc911a6c",
        "rules": "rulez"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 201

POST http://localhost:8080/boards/65b95156097680ef41e8f930/posts
{
    "post": {
        "id": "65b95f86e65c69d83a76c2e5",
        "title": "first post",
        "bodytype": 0,
        "bodycontent": "This is my first post",
        "votes": 501,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f930"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 201

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5
HTTP 200

POST http://localhost:8080/boards/65b95156097680ef41e8f930/posts
{
    "post": {
        "id": "65b96090a21f9d310d726750",
        "title": "second post",
        "bodytype": 0,
        "bodycontent": "This is my second post",
        "votes": 250,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f930"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 201

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts
HTTP 200

PUT http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5
{
    "post": {
        "id": "65b95f86e65c69d83a76c2e5",
        "title": "first post updated",
        "bodytype": 0,
        "bodycontent": "This is an edited post again",
        "votes": 501,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f930"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 202

PUT http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5
{
    "post": {
        "id": "65b95f86e65c69d83a76c2e5",
        "title": "Post 1",
        "bodytype": 0,
        "bodycontent": "This is an edited post",
        "votes": 501,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f930"
    },
    "requester": {
        "name": "regular_user2",
        "password": "password5"
    }
}
HTTP 403

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts
HTTP 200

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts/search?title=post

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5
{
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 200

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b96090a21f9d310d726750
{
    "requester": {
        "name": "regular_user2",
        "password": "password5"
    }
}
HTTP 403

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b96090a21f9d310d726750
{
    "requester": {
        "name": "Mod1",
        "password": "password1"
    }
}
HTTP 200

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930
{
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}
HTTP 200
