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

POST http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments
{
    "comment": {
        "id": "65b999c4f33023deae33606b",
        "author": "65b954c547c4f420dc911a6c",
        "post": "65b95f86e65c69d83a76c2e5",
        "body": "a comment",
        "votes": 10
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 201

POST http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments
{
    "comment": {
        "id": "65b99a2b3ccfffc3ef96db65",
        "author": "65b954c547c4f420dc911a6c",
        "post": "65b95f86e65c69d83a76c2e5",
        "body": "another comment",
        "votes": 10
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 201

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments
HTTP 200

GET http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments/65b999c4f33023deae33606b
HTTP 200

PUT http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments/65b999c4f33023deae33606b
{
    "comment": {
        "id": "65b999c4f33023deae33606b",
        "author": "65b954c547c4f420dc911a6c",
        "post": "65b95f86e65c69d83a76c2e5",
        "body": "edited comment",
        "votes": 10
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 202

PUT http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments/65b999c4f33023deae33606b
{
    "comment": {
        "id": "65b999c4f33023deae33606b",
        "author": "65b954c547c4f420dc911a6c",
        "post": "65b95f86e65c69d83a76c2e5",
        "body": "edited comment by administrator",
        "votes": 10
    },
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}
HTTP 202

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments/65b999c4f33023deae33606b
{
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 200

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5/comments/65b99a2b3ccfffc3ef96db65
{
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}
HTTP 200

DELETE http://localhost:8080/boards/65b95156097680ef41e8f930/posts/65b95f86e65c69d83a76c2e5
{
    "requester": {
        "name": "Administrator",
        "password": "passsword"
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
