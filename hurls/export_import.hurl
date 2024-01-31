POST http://localhost:8080/export
{
    "name": "regular_user",
    "password": "password4"
}

POST http://localhost:8080/export
{
    "name": "Administrator",
    "password": "passsword"
}

POST http://localhost:8080/import
{
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    },
    "comments": [
        {
            "author": "65b98e43380afb815dcadb8b",
            "post": "65b9846c7ef4232d56209d10",
            "body": "a comment",
            "votes": 50
        }
    ],
    "users": [
        {
            "id": "65b98e43380afb815dcadb8b",
            "name": "regular_user3",
            "bio": "Dictator",
            "avatar": "base64encodedfile",
            "pronouns": "over/lord",
            "password": "password10",
            "email": "mail@mail.com"
        }
    ],
    "boards": [
        {
            "id": "65b95156097680ef41e8f932",
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
        {
            "id": "65b95156097680ef41e8f935",
            "name": "technology",
            "bio": "board about technology",
            "moderators": [
                "65b9521f08488450adcbd92d",
                "65b9521f08488450adcbd92e",
                "65b9521f08488450adcbd92f"
            ],
            "owner": "65b94ef156e6d7c59f478392",
            "rules": "rulez"
        }
    ],
    "posts": [
        {
            "id": "65b9846c7ef4232d56209d10",
            "title": "first post",
            "bodyType": 0,
            "bodyContent": "This is my first post",
            "votes": 501,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f935"
        },
        {
            "id": "65b9846c7ef4232d56209d0b",
            "title": "second post",
            "bodyType": 0,
            "bodyContent": "This is my second post",
            "votes": 1501,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f935"
        },
        {
            "id": "65b9846c7ef4232d56209d0c",
            "title": "third post",
            "bodyType": 0,
            "bodyContent": "This is my third post",
            "votes": 1001,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f935"
        },
        {
            "id": "65b9846c7ef4232d56209d0d",
            "title": "fourth post",
            "bodyType": 0,
            "bodyContent": "This is my fourth post",
            "votes": 1201,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f932"
        },
        {
            "id": "65b9846c7ef4232d56209d0e",
            "title": "fourth post",
            "bodyType": 0,
            "bodyContent": "This is my fourth post",
            "votes": -1201,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f932"
        },
        {
            "id": "65b9846c7ef4232d56209d0f",
            "title": "fourth post",
            "bodyType": 0,
            "bodyContent": "This is my fourth post",
            "votes": -2402,
            "author": "65b954c547c4f420dc911a6c",
            "board": "65b95156097680ef41e8f932"
        }
    ]
}
