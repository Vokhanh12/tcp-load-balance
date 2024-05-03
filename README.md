# A TCP Server/Client/load-balance in Go #

This is the client-server load balancing process. Here, I create a server (load balancer) in between to distribute the load among the servers behind it and solve problems returning results to the client.


The problem-solving task:


![baitoan](https://github.com/Vokhanh12/tcp-load-balance-client-server/assets/36543564/c2bed836-e9b6-4307-8910-28e49e9f353d)


Modles:


![image](https://github.com/Vokhanh12/tcp-load-balance-client-server/assets/36543564/065c9b84-6340-4679-b6b5-bcbfd80ff5d5)


![image](https://github.com/Vokhanh12/tcp-load-balance-client-server/assets/36543564/2e204417-e37e-4b67-a10d-b383a79f5219)




# Usage #

Hey, look, it uses `flag`, how quaint!

Run the servers like so:

```bash
 ~/server-8000   
go run server.go -port 8000


 ~/server-8001   
go run server.go -port 8001


 ~/server-8002   
go run server.go -port 8002


 ~/server-8003   
go run server.go -port 8003
```


Run the load balance server like so:

```bash
~/load-balance
go run server.go -port 9999
```


And then connect to it with the client like so:


```bash
go run client.go -host localhost -port 9999
```



# License #

This is "do whatever you want with it"-ware. There is nothing here that is
particularly novel or valuable. Obviously this software comes with no warranty
of any kind. It might cause your computer to become self-aware and destroy
you. I take no responsibility for any outcomes, Skynet or otherwise.
