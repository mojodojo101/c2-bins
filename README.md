# c2-bins


## Greetings, this is a little command and control server i am working on

##### alot of the files are subject to change and will improve over time

##### POC can be found [here](https://youtu.be/NPZT92Imsnc)

#### TODO LIST FOR ME:

* Add a way for clients to retrieve beacons
* Add beacon obfuscation
* Add multiple beacon templates 
* Add propper client authentication (might try to read up on firefly crypt)
* Maybe redo the structure of the domain logic (usecases importing another)
* Add my webclient to this repo and add a GOOD commandline client 

How to set this up :

Install docker and docker-compose 

Add a server.crt and server.key to c2server/cmd/c2api/ folder (x509)

Using certbot to get a valid signed cert:

sudo certbot standalone --http-01-port 8888

sudo certbot certonly --standalone

cp /etc/letsencrypt/live/<your-domain>/fullchain.pem c2server/cmd/c2api/server.crt

cp /etc/letsencrypt/live/<your-domain>/privkey.pem c2server/cmd/c2api/server.key 



Check the docker-compose.yml and change the ENV variables to fit your needs.


```
 And if you want to use this for anything that shouldnt be volatile,
 mount psql in the docker and compose file and do the same with the c2server
 mount the target file of the c2server somewhere 
```


```sh

	docker-compose up --build -d

```


