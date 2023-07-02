### Dice-Game-Api

A Dice game api that you will enjoy. roll and win

Postman documentation https://documenter.getpostman.com/view/14213539/2s93sc4CnB



<br>

```GO

How to play.


Register player and login with your username and password.
Player funds account with 155 sats. player starts a game, which costs him 20 sats. 20 sat is deducted from the user’s wallet bringing the user’s balance down to 135 sats. A number generated at strat of game eg.  7 

Player starts a dice duo roll; on this first roll, the user is charged 5 sat bringing the user’s wallet balance down from 135 sat to 130 sat. Let’s assume the die roll result is 2. 

For the player to win, we can only hope that the player rolls 5 on the next roll. 

The player rolls the die again; this time he is not charged because we’ve already been charged for this session. Let us assume the user rolls 5.

The sum of the first and second die rolls (5+2=7) is 7, which means the player just won, and should thereby be awarded 10 dice, bringing the players’s wallet balance up from 130 sat to 140 sat. 

At this point, the dice roll duo session has been completed, but our game session is still active. Which means we can keep rolling. 

If the player’s dual roll does not win, well, he has lost his 5 sat and can retry. 

* start program
* register user       POST localhost:8000/api/v1/register 
* login             POST localhost:8000/api/v1/login
* fund wallet       GET localhost:8000/api/v1/credit
* start game        GET localhost:8000/api/v1/start
* roll dice          GET localhost:8000/api/v1/roll
* end game          GET localhost:8000/api/v1/end
* get transactions   GET localhost:8000/api/v1/transactions

```

## Requirements
*** 

* #### PostMan
* #### Postgresql
* #### Docker (Optional)
* #### Fiber Router


### Start Program

```GO
terminal> 
>> go mod tidy

>> make setup-docker // if you have docker running on your machine

for local database, set up postgres // set up postgres on your machine 

Port = 4300
USER=root 
PASSWORD=root 
DATABASE=dice_game

//on the terminal run 

>> make migrate-up
>> make run

```
