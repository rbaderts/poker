##Model

A poker Table has a fixed number of seats.     Each of the seats can be occupied by a user or empty.   
When a User joins a table they must provide a minimum amount of funds to gamble with, an amount
specified by the table, this is the players "table stack".  The table stack grows and shrinks as 
the player makes bets and wins hands.   When the player exits the table the players table stack is
transferred to the users main account.  A User may add to their table stack in between poker games.   

A poker table continually deals new Games as long as 2 or more seats are occupied and not sitting out.

  

A Game represents a single game of poker.    A Game starts with 2 or more Players (represented
by occupied seats at the table).   Upon game completion the Pot(s) are distributed to the
winning Player(s).    

A Game's workflow consists of requesting a bet from each player, a player has a fixed time to 
responds to the request, after which a default action is taken (check or fold).    



##Game  Control flow
The poker game engine consists of driver function which manages the control flow of the game.

For betting rounds the game driver suspends control handing it over to a betting round driver ,
until the betting round completes and returns control.

###Game Actions
The control flow of a game  can be fully described by a sequential list of Game actions.
These are the actions triggered by the game engine or by users.

The current state of a game changes only in response to executing game actions


#####Events
Events describes some state change to some part of the app state (table, game, player, etc)
oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo
Events exists to communicate state changes to attached clients.   They currently are
delivered via websockets with the aid of Centrifugo

#####Messages 
Messages are used to send info about user actions, which 
result int he construction of a  game engine action


