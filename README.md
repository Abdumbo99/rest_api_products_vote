# Rest API Go - "Products Vote"

> 🌱 This projetct implements products "Tinder"

## 🚀 Endpoints

| Name                           | HTTP Method | Route               |
|--------------------------------|-------------|---------------------|
| List Products                  | GET         | /products           |
| List Votes                     | GET         | /votes              |
| Submit/Update vote             | POST        | /votes              |
| List votes of a product        | GET         | /votes/product/{id} |
| List votes of a session        | GET         | /votes/session/{id} |
| List average votes per product | GET         | /products/avgs      |

## 🗄️ Database design

| Column Name    | Datatype  | Primary Key |
|----------------|-----------|-------------|
| product_id     | TEXT      | ✅          |
| session_id     | UUID      | ✅          |
| rate           | INT       |             |

| Column Name    | Datatype  | Primary Key |
|----------------|-----------|-------------|
| product_id     | TEXT      | ✅          |
| product_name   | TEXT      |             |

## 📁 Project structure

```shell
foodji_assignment
├── cmd
│  ├── api
│     └── main.go
│
├── api
│  ├── models
│  │  ├── vote
│  │  │  ├── vote.go
│  │  │  ├── repository.go
│  │  ├── Product
│  │     └── product.go
│  │
│  │── middleware
│  │  ├── cors.go
│  │  │── logger.go
│  │  └── session_id.go
│  │
│  └── handler
│     ├── hanlder.go
│     │── handler_test.go
│     └── mock.go
│
├── .env
│
├── go.mod
├── go.sum
│
├── products.json
│
└── Dockerfile

```

## 🚀 Cloud Deployment

The code was deployed on <https://render.com> and Mongo's Atals, and can be accessed through the following URI <https://products-vote.onrender.com> (render sleeps after long time of no use so please keep in mind).

## 🚀 Local Deployment

1. set up the .env parameters to connect to your database, or keep them to connect to the remote database
2. build the docker image through `docker build .`
3. run the image `docker run -p 80:8080  {IMAGE_ID}`
4. Test it!

## 🚀 Calling the API

1. **Posting/updating a vote**: for posting/updating a vote all you have to do is calling the endpoint `https://products-vote.onrender.com/votes` with the data of the vote included in the following structure `'{"product_id":{id}, "rate":{int}}'`. In case the vote already exists it automatically updates it, without duplication.
2. **Listing Products**: for viewing all products in the system call the endpoint `https://products-vote.onrender.com/products`.
3. **Listing Votes**: for viewing all products in the system call the endpoint `https://products-vote.onrender.com/votes` while this orignially was not required, it is usefull for validation purposes to be able to see the votes, additionaly there are no other practical ways to view session ids (save checking the cookie's content).
4. **Listing votes of a specific product**: to list votes of a specific product call `https://products-vote.onrender.com/votes/product/{id}`. This, again, was not required, but come in handy for testing and validating the  system.
5. **Listing votes of a session**: to list votes of a specific session call `https://products-vote.onrender.com/votes/session/{id}`.
6. **Listing average votes per product**: to calculate the avg. vote/rate of each product call `https://products-vote.onrender.com/products/avgs`

## 🚀 Requests Examples

While the get calls can be performed easily through any means, browser, postman, etc. A list of curl requests are provided below:

    // View the available products
    curl --location -X  GET 'https://products-vote.onrender.com/products' -c cookies.txt --header 'Content-Type: text/plain'
    // View the votes so far
    curl --location -X GET 'https://products-vote.onrender.com/votes' -b cookies.txt --header 'Content-Type: text/plain'
    // Post your vote
    curl --location -X POST 'https://products-vote.onrender.com/votes' -b cookies.txt --header 'Content-Type: text/plain' --data '{"product_id":"3", "rate":10}'
    // Get all Votes for a specific product
    curl --location -X  GET 'https://products-vote.onrender.com/votes/product/1' -c cookies.txt --header 'Content-Type: text/plain'
    // Get all Votes for a specific session
    curl --location -X  GET 'https://products-vote.onrender.com/votes/session/b5b5c578-b561-4fef-9366-ee21e5d21e3a' -c cookies.txt --header 'Content-Type: text/plain' 
    // Fetch the avg. vote for each product
    curl --location -X  GET 'https://products-vote.onrender.com/products/avgs' -c cookies.txt --header 'Content-Type: text/plain'
    // Call the docker container
    curl -b cookies.txt -X GET http://localhost:80/votes/session/b5b5c578-b561-4fef-9366-ee21e5d21e3a -H "content-Type: application/json" 
