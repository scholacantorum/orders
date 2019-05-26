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
    "receipt": "<p>Thank you for your generous donation of ${{ dollars .Amount }} to Schola Cantorum, which supports our mission to bring choral music to the community through live performances and educational outreach programs. Schola Cantorum is a 501(c)(3) tax-exempt organization. Our tax ID number is 94-2597822.</p>",
    "skus": [{}]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "summer-sings-2019",
    "name": "2019 Summer Sing Ticket",
    "ticketName": "General Admission",
    "type": "ticket",
    "ticketCount": 1,
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} 2019 Summer Sings ticket{{ if gt .Quantity 1 }}s{{ end }}, for ${{ dollars .Amount }}.  {{ if gt .Quantity 1 }}Each{{ else }}Your{{ end }} ticket is good for entry to any one of the sings.  The sings are on Monday nights from July 8 through August 12, 2019 (see <a href=\"https://scholacantorum.org/summer-sings\">schedule</a>).  Each sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": ["2019-07-08", "2019-07-15", "2019-07-22", "2019-07-29", "2019-08-05", "2019-08-12"],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 1700
    }, {
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "quantity": 6,
        "price": 8500
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "summer-sings-2019-student",
    "name": "2019 Summer Sing Ticket (Student)",
    "ticketName": "Student",
    "type": "ticket",
    "ticketCount": 1,
    "ticketClass": "Student",
    "events": ["2019-07-08", "2019-07-15", "2019-07-22", "2019-07-29", "2019-08-05", "2019-08-12"],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 0
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2018-11-03",
    "name": "Concert Recording: 2018-11 For the Love of Bach",
    "type": "recording",
    "receipt": "<p>We confirm your purchase of an archival recording of \"For the Love of Bach\", November 2018, for {{ dollars .Amount }}. Your recording is available for download from the members web site.</p>",
    "skus": [{
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "membersOnly": true,
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2018-12-16",
    "name": "Concert Recording: 2018-12 A John Rutter Christmas",
    "type": "recording",
    "receipt": "<p>We confirm your purchase of an archival recording of \"A John Rutter Christmas\", December 2018, for {{ dollars .Amount }}. Your recording is available for download from the members web site.</p>",
    "skus": [{
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "membersOnly": true,
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2019-03-16",
    "name": "Concert Recording: 2019-03 Carmina Burana",
    "type": "recording",
    "receipt": "<p>We confirm your purchase of an archival recording of \"Carmina Burana\", March 2019, for {{ dollars .Amount }}. Once the recording is made available, you can download it from the members web site.</p>",
    "skus": [{
        "salesStart": "2019-06-01T00:00:00-07:00",
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "membersOnly": true,
        "price": 2000
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "recording-2019-05-24",
    "name": "Concert Recording: 2019-05 Ein deutsches Requiem",
    "type": "recording",
    "receipt": "<p>We confirm your purchase of an archival recording of \"Ein deutsches Requiem\", May 2019, for {{ dollars .Amount }}. Once the recording is made available, you can download it from the members web site.</p>",
    "skus": [{
        "salesStart": "2019-06-01T00:00:00-07:00",
        "salesEnd": "2019-08-01T00:00:00-07:00",
        "membersOnly": true,
        "price": 2000
    }]
}' http://localhost:8100/api/product
