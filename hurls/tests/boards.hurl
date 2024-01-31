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
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 403

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
HTTP 201

GET http://localhost:8080/boards
HTTP 200

GET http://localhost:8080/boards/65b95156097680ef41e8f929
HTTP 200

PUT http://localhost:8080/boards/65b95156097680ef41e8f929
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
        "rules": "new rules"
    },
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}
HTTP 202

PUT http://localhost:8080/boards/65b95156097680ef41e8f929
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
        "rules": "new rules"
    },
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 403

GET http://localhost:8080/boards/search?name=sci
HTTP 200

DELETE http://localhost:8080/boards/65b95156097680ef41e8f929
{
    "requester": {
        "name": "regular_user",
        "password": "password4"
    }
}
HTTP 403

DELETE http://localhost:8080/boards/65b95156097680ef41e8f929
{
    "requester": {
        "name": "Administrator",
        "password": "passsword"
    }
}
HTTP 200
