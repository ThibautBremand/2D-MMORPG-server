# 2D-MMORPG-server
The server for 2D-MMORPG, written in Golang.
For the client, please see https://github.com/ThibautBremand/2D-MMORPG-client

This server handles communications from users, and interacts with the databases accordingly. It uses Redis and PostgreSQL.

#### Local installation

- Clone the repo
- Launch a local Redis instance (default port will be 6379)
    * Docker (recommended): *docker run --name some-redis -p 6379:6379 -d redis*
    * Ubuntu: *sudo service redis-server start*
- Launch a local PostgreSQL instance
    * Docker (recommended): [See Wiki page](https://github.com/ThibautBremand/2D-MMORPG-server/wiki/Configure-local-PostgreSQL-with-Docker)
    * Ubuntu: *sudo service postgresql start*
- Make sure you have data in your *Character* and *Gamemap* PostgreSQL tables.
  - Will be detailed later
- Create a .env file based on the .env.sample model, and fill the values.
  - **Note:** CLIENT_PATH value must correspond to the DEPLOY_PATH value of the client. That way, once you deploy the client (cf. https://github.com/ThibautBremand/2D-MMORPG-client), it'll be directly served by the server to the users.
- Build: *go build*
- Run: *go run main.go*

Please follow the steps detailed in the client's Readme file in order to correctly deploy the client.

## 2k20 Reborn!
This project is a new version of this one made in PHP with Symfony, in 2015:
https://github.com/ThibautBremand/WebApp_WebMMORPG-Server

The previous project is now deprecated.