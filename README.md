**Description**

This project provides a backend solution for an **online reading and audiobook platform**, enabling RESTful access to books, authors, and genres. It serves as a core component of a larger digital reading and streaming system.

**Tech Stack**

- **Backend**: Golang  
- **Routing**: httprouter
- **Database**: MongoDB  
- **Containerization**: Docker  

**Docker Usage**

To run the project using Docker, follow these steps:

1.Setup `.env` parameters as below :

```
MONGODB_URI=""# ATLAS URI or local MongoDB URI
DB_NAME=book # Database name
PORT=8000 # Port for the server
HOST=0.0.0.0 # Host for the server
```

2.Ensure that Docker is installed on your system.

3.Open a terminal or command prompt.

4.Navigate to the project directory.

5.Run the following command:

    sudo docker-compose build
    sudo docker-compose up -d

