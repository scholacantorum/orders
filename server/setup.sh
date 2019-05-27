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
    "id": "ticket-2019-07-08",
    "name": "Ticket to July 8 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 8, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 0},
        {"event": "2019-07-15", "priority": 30},
        {"event": "2019-07-22", "priority": 30},
        {"event": "2019-07-29", "priority": 30},
        {"event": "2019-08-05", "priority": 30},
        {"event": "2019-08-12", "priority": 30}
    ],
    "skus": [{
        "salesEnd": "2019-07-08T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-15",
    "name": "Ticket to July 15 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 15, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 30},
        {"event": "2019-07-15", "priority": 0},
        {"event": "2019-07-22", "priority": 30},
        {"event": "2019-07-29", "priority": 30},
        {"event": "2019-08-05", "priority": 30},
        {"event": "2019-08-12", "priority": 30}
    ],
    "skus": [{
        "salesEnd": "2019-07-15T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-22",
    "name": "Ticket to July 22 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 22, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 30},
        {"event": "2019-07-15", "priority": 30},
        {"event": "2019-07-22", "priority": 0},
        {"event": "2019-07-29", "priority": 30},
        {"event": "2019-08-05", "priority": 30},
        {"event": "2019-08-12", "priority": 30}
    ],
    "skus": [{
        "salesEnd": "2019-07-22T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-07-29",
    "name": "Ticket to July 29 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, July 29, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 30},
        {"event": "2019-07-15", "priority": 30},
        {"event": "2019-07-22", "priority": 30},
        {"event": "2019-07-29", "priority": 0},
        {"event": "2019-08-05", "priority": 30},
        {"event": "2019-08-12", "priority": 30}
    ],
    "skus": [{
        "salesEnd": "2019-07-29T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-08-05",
    "name": "Ticket to August 5 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 5, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 30},
        {"event": "2019-07-15", "priority": 30},
        {"event": "2019-07-22", "priority": 30},
        {"event": "2019-07-29", "priority": 30},
        {"event": "2019-08-05", "priority": 0},
        {"event": "2019-08-12", "priority": 30}
    ],
    "skus": [{
        "salesEnd": "2019-08-05T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "ticket-2019-08-12",
    "name": "Ticket to August 12 Summer Sing",
    "type": "ticket",
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 12, for ${{ dollars .Amount }}.  The sings starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 30},
        {"event": "2019-07-15", "priority": 30},
        {"event": "2019-07-22", "priority": 30},
        {"event": "2019-07-29", "priority": 30},
        {"event": "2019-08-05", "priority": 30},
        {"event": "2019-08-12", "priority": 0}
    ],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 1700
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "summer-sings-2019",
    "name": "2019 Summer Sings Flex Pass",
    "type": "ticket",
    "ticketCount": 6,
    "receipt": "<p>We confirm your purchase of {{ .Quantity }} Flex Pass{{ if gt .Quantity 1 }}es{{ end }} to the 2019 Summer Sings, for ${{ dollars .Amount }}.  {{ if gt .Quantity 1 }}Each{{ else }}Your{{ end }} pass is good for six entries to the sings: six people on one night, one person on six nights, or any mixture.  The sings are on Monday nights from July 8 through August 12, 2019 (see <a href=\"https://scholacantorum.org/summer-sings\">schedule</a>).  Each sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href=\"https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889\">map</a>). Please bring this email (printed or on your phone) for admission.</p>",
    "events": [
        {"event": "2019-07-08", "priority": 20},
        {"event": "2019-07-15", "priority": 20},
        {"event": "2019-07-22", "priority": 20},
        {"event": "2019-07-29", "priority": 20},
        {"event": "2019-08-05", "priority": 20},
        {"event": "2019-08-12", "priority": 20}
    ],
    "skus": [{
        "salesEnd": "2019-08-12T16:30:00-07:00",
        "price": 8500
    }]
}' http://localhost:8100/api/product
curl -i -d'{
    "id": "summer-sings-2019-student",
    "name": "2019 Summer Sing Student Entry",
    "type": "ticket",
    "ticketCount": 1,
    "ticketClass": "Student",
    "events": [
        {"event": "2019-07-08", "priority": 20},
        {"event": "2019-07-15", "priority": 20},
        {"event": "2019-07-22", "priority": 20},
        {"event": "2019-07-29", "priority": 20},
        {"event": "2019-08-05", "priority": 20},
        {"event": "2019-08-12", "priority": 20}
    ],
    "skus": [{
        "salesEnd": "2019-01-01T00:00:00-07:00",
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
