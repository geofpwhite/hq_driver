Engine for various html games.

All game states implement a Game interface on the backend. That way games can be handled by the same thread despite being different games.

Whenever a user makes an action that would change the game state, a UserInput object is created that implements an Input interface. These input interfaces use a ChangeState function that call mutator methods on the user's current game.


User makes action --> Websocket message is sent --> user input object is created and sent to game loop thread via a channel--> game loop calls ChangeState, which calls one of the mutator functions on the game state --> new game state sent to output loop via a channel --> new game state is sent to all players registered for that game via their web socket
