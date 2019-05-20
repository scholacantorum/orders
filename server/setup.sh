curl -i -d'{
    "id": "2019-07-08",
    "name": "Summer Sing",
    "start": "2019-07-08T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "2019-07-15",
    "name": "Summer Sing",
    "start": "2019-07-15T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "2019-07-22",
    "name": "Summer Sing",
    "start": "2019-07-22T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "2019-07-29",
    "name": "Summer Sing",
    "start": "2019-07-29T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "2019-08-05",
    "name": "Summer Sing",
    "start": "2019-08-05T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "2019-08-12",
    "name": "Summer Sing",
    "start": "2019-08-12T19:30:00-07:00"
}' http://localhost:8100/api/event
curl -i -d'{
    "id": "donation",
    "name": "Donation",
    "type": "donation",
    "skus": [{}]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-08",
    "name": "Ticket to Summer Sing, July 8, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-07-08"}],
    "skus": [{
        "salesEnd": "2019-07-08T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-15",
    "name": "Ticket to Summer Sing, July 15, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-07-15"}],
    "skus": [{
        "salesEnd": "2019-07-15T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-22",
    "name": "Ticket to Summer Sing, July 22, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-07-22"}],
    "skus": [{
        "salesEnd": "2019-07-22T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-29",
    "name": "Ticket to Summer Sing, July 29, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-07-29"}],
    "skus": [{
        "salesEnd": "2019-07-29T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-08-05",
    "name": "Ticket to Summer Sing, August 5, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-08-05"}],
    "skus": [{
        "salesEnd": "2019-08-05T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-08-12",
    "name": "Ticket to Summer Sing, August 12, 2019",
    "type": "ticket",
    "ticketCount": 1,
    "events": [{"id": "2019-08-12"}],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "summer-sings-2019",
    "name": "2019 Summer Sings Flex Pass",
    "type": "flexpass",
    "ticketCount": 6,
    "events": [
        {"id": "2019-07-08"},
        {"id": "2019-07-15"},
        {"id": "2019-07-22"},
        {"id": "2019-07-29"},
        {"id": "2019-08-05"},
        {"id": "2019-08-12"}
    ],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 8500
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2018-11-03",
    "name": "Concert Recording: 2018-11 For the Love of Bach",
    "type": "recording",
    "skus": [{
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2018-12-16",
    "name": "Concert Recording: 2018-12 A John Rutter Christmas",
    "type": "recording",
    "skus": [{
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2019-03-16",
    "name": "Concert Recording: 2019-03 Carmina Burana",
    "type": "recording",
    "skus": [{
        "salesStart": "2019-06-01T00:00:00-07:00",
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2019-05-24",
    "name": "Concert Recording: 2019-05 Ein deutsches Requiem",
    "type": "recording",
    "skus": [{
        "salesStart": "2019-06-01T00:00:00-07:00",
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "price": 2000
    }]
}' http://localhost:8100/api/product
