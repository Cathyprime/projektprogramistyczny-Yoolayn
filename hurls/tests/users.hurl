POST http://localhost:8080/users
{
    "user": {
        "id": "65b944449980e20df0c2f3ef",
        "name": "example_user1",
        "bio": "example_bio1",
        "password": "THY END IS NOW",
        "avatar": "avatar string",
        "pronouns": "she/her",
        "email": "mail@email.com"
    }
}

POST http://localhost:8080/users
{
    "user": {
        "name": "user2",
        "bio": "example_bio1",
        "password": "THY END IS NOW",
        "avatar": "avatar string",
        "pronouns": "she/her",
        "email": "mail@email.com"
    }
}

GET http://localhost:8080/users/65b944449980e20df0c2f3ef
GET http://localhost:8080/users

PUT http://localhost:8080/users/65b944449980e20df0c2f3ef
{
    "user": {
        "name": "example_user1",
        "bio": "example_bio1",
        "password": "THY END IS NOW",
        "avatar": "avatar string",
        "pronouns": "she/her",
        "email": "mail@email.com"
    },
    "requester": {
        "name": "example_user1",
        "password": "THY END IS NOW"
    }
}

PUT http://localhost:8080/users/65b944449980e20df0c2f3ef
{
    "user": {
        "name": "example_user1",
        "bio": "example_bio1",
        "password": "THY END IS NOW",
        "avatar": "avatar string",
        "pronouns": "she/her",
        "email": "mail@email.com"
    },
    "requester": {
        "name": "example_user1",
        "password": "pad passwoerd"
    }
}

GET http://localhost:8080/users/search?name=example

DELETE http://localhost:8080/users/65b944449980e20df0c2f3ef
{
    "requester": {
        "name": "example_user1",
        "password": "THY END IS NOW"
    }
}
