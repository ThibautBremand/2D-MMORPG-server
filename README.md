# 2D-MMORPG-server
The server for 2D-MMORPG, written in Golang.
For the client, please see https://github.com/ThibautBremand/2D-MMORPG-client

This server handles communications from users, and interacts with the databases accordingly. It uses Redis and PostgreSQL.

#### Local installation

- Clone the repo
- Launch a local Redis instance (default port will be 6379)
    * Ubuntu: *sudo service redis-server start*
- Launch a local PostgreSQL instance
    * Ubuntu: *sudo service postgresql start*
- Create a .env file based on the .env.sample model, and fill the values.
- Build: *go build*
- Run: *go run main.go*

The client needs to be deployed into the *client directory*, otherwise the server won't have anything to serve to the user! Please follow the steps detailed in the client's Readme file in order to correctly deploy the client.

## 2k20 Reborn!
This project is a new version of this one made in PHP with Symfony, in 2015:
https://github.com/ThibautBremand/WebApp_WebMMORPG-Server
The previous project is now deprecated.