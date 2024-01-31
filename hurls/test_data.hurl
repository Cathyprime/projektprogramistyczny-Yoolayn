POST http://localhost:8080/boards
{
    "board": {
        "id": "65b95156097680ef41e8f929",
        "name": "science",
        "bio": "board about science",
        "moderators": [
            "65b9521f08488450adcbd92d",
            "65b9521f08488450adcbd92e",
            "65b9521f08488450adcbd92f"
        ],
        "owner": "65b94ef156e6d7c59f478392",
        "rules": "rulez"
    },
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}

POST http://localhost:8080/boards
{
    "board": {
        "id": "65b95156097680ef41e8f932",
        "name": "technology",
        "bio": "board about technology",
        "moderators": [
            "65b9521f08488450adcbd92d",
            "65b9521f08488450adcbd92e",
            "65b9521f08488450adcbd92f"
        ],
        "owner": "65b94ef156e6d7c59f478392",
        "rules": "rulez"
    },
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f932/posts
{
    "post": {
        "title": "first post",
        "bodytype": 0,
        "bodycontent": "This is my first post",
        "votes": 501,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f932"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f932/posts
{
    "post": {
        "title": "second post",
        "bodytype": 0,
        "bodycontent": "This is my second post",
        "votes": 1501,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f932"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f932/posts
{
    "post": {
        "title": "third post",
        "bodytype": 0,
        "bodycontent": "This is my third post",
        "votes": 1001,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f932"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f929/posts
{
    "post": {
        "title": "fourth post",
        "bodytype": 0,
        "bodycontent": "This is my fourth post",
        "votes": 1201,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f929"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f929/posts
{
    "post": {
        "title": "fourth post",
        "bodytype": 0,
        "bodycontent": "This is my fourth post",
        "votes": -1201,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f929"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

POST http://localhost:8080/boards/65b95156097680ef41e8f929/posts
{
    "post": {
        "title": "fourth post",
        "bodytype": 0,
        "bodycontent": "This is my fourth post",
        "votes": -2402,
        "author": "65b954c547c4f420dc911a6c",
        "board": "65b95156097680ef41e8f929"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}

GET http://localhost:8080
