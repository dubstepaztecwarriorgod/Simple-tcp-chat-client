# Simple chat client

In order to use the chat client
```
git clone https://github.com/dubstepaztecwarriorgod/Simple-tcp-chat-client.git
cd server
cargo run <SERVER_ADDRESS>
```

Then to run the client
```
cd .. (move out of the server dir if you're still in it)
cd client
go run client.go <SERVER_ADDRESS>
```

If either of the server address arguments are not passed then they will default to `127.0.0.1:8080`

There's also some commands you can use!
|  Commands         |  Functionality
|-------------------|--------------
|  /help            | Displays the help menu
|  /addr            | Displays the address of the server
|  /quit            | Quits the apllication cleanly
|  /limit           | Displays the character limit for messages
|  /file_Send FILE  | Sends the file you pass in to the server

Happy chatting!
![alt txt](https://cdn.discordapp.com/attachments/1048381919362035803/1237097351223578735/Screen_Shot_2024-05-06_at_10.42.40_AM.png?ex=663a6830&is=663916b0&hm=311aac36c0343b34fdadabb258b6d228f39d93f0bb5dd6b57c03a1ab0f077a61& "client image")





