# example-go-graphql

Example Go using Graphql

Hi guys, I was making an example of Go GraphQL.
In this example, I was only use in memory data without database connection.

Make sure you already have go installed in your local system.

## Steps
Below are steps that you can follow before can test the code :
1. clone the repository

2. exec `go mod download`

3. after that simply run `go run main.go`

4. check the terminal and see there is this message : `Server is running on port 8080`

5. copy this url : `http://localhost:8080/api?query={users{id,name}}` and paste it into your browser

6. there will be response list users : 
```
data: {
    users: [
        {
            id: 1,
            name: "Dicky"
        },
        {
            id: 2,
            name: "Josephine"
        }
    ]
}
```

### Contact
Send a message to dw_authorized@yahoo.co.id
