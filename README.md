Client:

Game has two responsibilities:
1. Drawing to the screen
2. Sending updates to Wind

It shares the same state as servers, and only knows the most basic state of all

First step is to query if there is a server willing to join
If there is, try to connect to it
On first entry, send enter state to the server
Then only send messages on actions
Receive messages from the server throughout the life cycle of the game

Before the game starts is where HTTP calls may come into effect

Writes to a single server
The server handles any writes and replicates to the other servers
Reads from a pub-sub broadcast
