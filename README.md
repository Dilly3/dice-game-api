# Dice-Game-Api

### A Dice game api that you will enjoy. roll and win

## Requirements
*** 

* ### PostMan
* ### Postgresql
* ### Docker (Optional)
* ### Fiber Router




## EndPoints

* ###  POST localhost:8000/api/v1/register
* ###  POST localhost:8000/api/v1/login
* ###  POST localhost:8000/api/v1/credit
* ###  GET localhost:8000/api/v1/roll
* ###  GET localhost:8000/api/v1/balance
* ###  GET localhost:8000/api/v1/start
* ###  GET localhost:8000/api/v1/end
* ###  GET localhost:8000/api/v1/logout
* ###  GET localhost:8000/api/v1/transactions
* ###  GET localhost:8000/api/v1/session

* ### GET localhost:8000/{code}


<br>

## POST 
* ### localhost:8000/api/register

>> Request Body 

```GO
     {
       {
    "firstname" : "meghan",
   "lastname" : "good",
    "username" : "good90",
    "password" : "12345",
    "confirm_password" : "12345"
}
    }
```

>> Response 

```GO
{
    {
    "id": 69,
    "firstname": "meghan",
    "lastname": "good",
    "username": "good90",
    "created_at": "2023-06-09T19:22:31.654395Z"
}
}
```

<br>


## POST 
* ### localhost:8000/api/login

>> Request Body 

```GO
     
       
   {
    "username" : "good90",
   "password" : "12345"
}

    
```

>> Response 

```GO
{
    "message": "login successful",
    "data": null,
    "errors": null,
    "status": 200,
    "timestamp": "2023-06-09 21:36:00"
}
```

 

```

Run Program

```GO
terminal> 
>> RUN go mod tidy
>> RUN make setup-docker // if you have docker running on your machine

for local database, set up postgres

Port = 4300
USER=root 
PASSWORD=root 
DATABASE=dice_game

>> RUN make setup-air
>> RUN migrate-up


>> make run
```