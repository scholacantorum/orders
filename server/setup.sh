curl -i -d'{"name":"Summer Sing","start":"2019-07-08T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Summer Sing","start":"2019-07-15T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Summer Sing","start":"2019-07-22T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Summer Sing","start":"2019-07-29T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Summer Sing","start":"2019-08-05T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Summer Sing","start":"2019-08-12T19:30:00-07:00"}' http://localhost:8100/api/event
curl -i -d'{"name":"Donation","stripeID":"donation"}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"donation"}' http://localhost:8100/api/product/1/sku
curl -i -d'{"name":"Ticket to Summer Sing, July 8, 2019","stripeID":"ticket-2019-07-08","ticketCount":1,"events":[1]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-07-08","salesEnd":"2019-07-08T16:30:00-07:00"}' http://localhost:8100/api/product/2/sku
curl -i -d'{"name":"Ticket to Summer Sing, July 15, 2019","stripeID":"ticket-2019-07-15","ticketCount":1,"events":[2]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-07-15","salesEnd":"2019-07-15T16:30:00-07:00"}' http://localhost:8100/api/product/3/sku
curl -i -d'{"name":"Ticket to Summer Sing, July 22, 2019","stripeID":"ticket-2019-07-22","ticketCount":1,"events":[3]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-07-22","salesEnd":"2019-07-22T16:30:00-07:00"}' http://localhost:8100/api/product/4/sku
curl -i -d'{"name":"Ticket to Summer Sing, July 29, 2019","stripeID":"ticket-2019-07-29","ticketCount":1,"events":[4]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-07-29","salesEnd":"2019-07-29T16:30:00-07:00"}' http://localhost:8100/api/product/5/sku
curl -i -d'{"name":"Ticket to Summer Sing, August 5, 2019","stripeID":"ticket-2019-08-05","ticketCount":1,"events":[5]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-08-05","salesEnd":"2019-08-05T16:30:00-07:00"}' http://localhost:8100/api/product/6/sku
curl -i -d'{"name":"Ticket to Summer Sing, August 12, 2019","stripeID":"ticket-2019-08-12","ticketCount":1,"events":[6]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"ticket-2019-08-12","salesEnd":"2019-08-12T16:30:00-07:00"}' http://localhost:8100/api/product/7/sku
curl -i -d'{"name":"2019 Summer Sings Flex Pass","stripeID":"summer-sings-2019","ticketCount":6,"events":[1,2,3,4,5,6]}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"summer-sings-2019","salesEnd":"2019-08-12T16:30:00-07:00"}' http://localhost:8100/api/product/8/sku
curl -i -d'{"name":"Concert Recording: 2018-11 For the Love of Bach","stripeID":"recording-2018-11-03"}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"recording-2018-11-03","salesEnd":"2019-08-01T00:00:00-07:00"}' http://localhost:8100/api/product/9/sku
curl -i -d'{"name":"Concert Recording: 2018-12 A John Rutter Christmas","stripeID":"recording-2018-12-16"}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"recording-2018-12-16","salesEnd":"2019-08-01T00:00:00-07:00"}' http://localhost:8100/api/product/10/sku
curl -i -d'{"name":"Concert Recording: 2019-03 Carmina Burana","stripeID":"recording-2019-03-16"}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"recording-2019-03-16","salesStart":"2019-06-01T00:00:00-07:00","salesEnd":"2019-08-01T00:00:00-07:00"}' http://localhost:8100/api/product/11/sku
curl -i -d'{"name":"Concert Recording: 2019-05 Ein deutsches Requiem","stripeID":"recording-2019-05-24"}' http://localhost:8100/api/product
curl -i -d'{"stripeID":"recording-2019-05-24","salesStart":"2019-06-01T00:00:00-07:00","salesEnd":"2019-08-01T00:00:00-07:00"}' http://localhost:8100/api/product/12/sku
