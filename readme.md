# hyperloop-optimiser is a web service accessible via a REST API for computing the cheapest routes to build for a hyperloop network
This project was developped for the [2017 Prague Hackaton](https://praguehackaton.com).
Coded by [Alexandre A. Bizri](https://github.com/aabizri) in conjonction with Marin Godechot & Camille Marchetti (Team K-1000)
Licensed under UNLICENSE

## Building the server
In the directory:
`go build`

## Launching it
`./hyperloop-optimiser`
Flags:
- `-p` allows you to select the port to listen to (ex: `-p 8080`)
- `-logpath` is the directory in which to log requests (ex: `-logpath /tmp/hyperloop`)
- `-path` is the path to listen to (ex: `-path /task1/input`)
Calling without flags is equivalent to using `-p 8080 -logpath "/tmp/" -path "/"`

## Sample input
```json
{
  "citiesCount" : 4,
  "costOffers" : [
    {
      "from" : 0,
      "to" : 1,
      "price" : 6
    },{
      "from" : 1,
      "to" : 2,
      "price" : 10
    },{
      "from" : 2,
      "to" : 1,
      "price" : 10
    },{
      "from" : 1,
      "to" : 3,
      "price" : 12
    },{
      "from" : 3,
      "to" : 2,
      "price" : 8
    },{
      "from" : 3,
      "to" : 0,
      "price" : 1
    }
  ]
}
```

## Sample output
```json
{
  "feasible" : true,
  "totalCost" : 15,
  "depotId" : 3,
  "recommendedOffers" : [
    {
      "from" : 0,
      "to" : 1,
      "price" : 6
    },{
      "from" : 3,
      "to" : 0,
      "price" : 1
    },{
      "from" : 3,
      "to" : 2,
      "price" : 8
    }
  ]
}
```
