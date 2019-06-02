PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
INSERT INTO product VALUES('donation','Donation','Donation','donation','<p>Thank you for your generous donation of ${{ dollars .Price }} to Schola Cantorum, which supports our mission to bring choral music to the community through live performances and educational outreach programs. Schola Cantorum is a 501(c)(3) tax-exempt organization. Our tax ID number is 94-2597822.</p>',0,'');
INSERT INTO product VALUES('ticket-2019-07-08','July 8 Summer Sing','July 8','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 8, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('ticket-2019-07-15','July 15 Summer Sing','July 15','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 15, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('ticket-2019-07-22','July 22 Summer Sing','July 22','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 22, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('ticket-2019-07-29','July 29 Summer Sing','July 29','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, July 29, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('ticket-2019-08-05','August 5 Summer Sing','August 5','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 5, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('ticket-2019-08-12','August 12 Summer Sing','August 12','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 12, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES('summer-sings-2019','Summer Sings Flex Pass','Flex Pass','ticket','<p>We confirm your purchase of {{ .Quantity }} Flex Pass{{ if gt .Quantity 1 }}es{{ end }} to the 2019 Summer Sings, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  {{ if gt .Quantity 1 }}Each{{ else }}Your{{ end }} pass is good for six entries to the sings: six people on one night, one person on six nights, or any mixture.  The sings are on Monday nights from July 8 through August 12, 2019 (see <a href="https://scholacantorum.org/summer-sings">schedule</a>).  Each sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',6,'');
INSERT INTO product VALUES('summer-sings-2019-student','2019 Summer Sing Student Entry','Student','ticket','',1,'Student');
INSERT INTO product VALUES('recording-2018-11-03','Concert Recording: 2018-11 For the Love of Bach','2018-11','recording','<p>We confirm your purchase of an archival recording of "For the Love of Bach", November 2018, for {{ dollars .Price }}. Your recording is available for download from the members web site.</p>',0,'');
INSERT INTO product VALUES('recording-2018-12-16','Concert Recording: 2018-12 A John Rutter Christmas','2018-12','recording','<p>We confirm your purchase of an archival recording of "A John Rutter Christmas", December 2018, for {{ dollars .Price }}. Your recording is available for download from the members web site.</p>',0,'');
INSERT INTO product VALUES('recording-2019-03-16','Concert Recording: 2019-03 Carmina Burana','2019-03','recording','<p>We confirm your purchase of an archival recording of "Carmina Burana", March 2019, for {{ dollars .Price }}. Once the recording is made available, you can download it from the members web site.</p>',0,'');
INSERT INTO product VALUES('recording-2019-05-24','Concert Recording: 2019-05 Ein deutsches Requiem','2019-05','recording','<p>We confirm your purchase of an archival recording of "Ein deutsches Requiem", May 2019, for {{ dollars .Price }}. Once the recording is made available, you can download it from the members web site.</p>',0,'');
INSERT INTO sku VALUES('donation','','','',0,0);
INSERT INTO sku VALUES('ticket-2019-07-08','','2019-07-01 00:00:00','2019-07-08 16:30:00',0,1700);
INSERT INTO sku VALUES('ticket-2019-07-15','','','2019-07-15 16:30:00',0,1700);
INSERT INTO sku VALUES('ticket-2019-07-22','','','2019-07-22 16:30:00',0,1700);
INSERT INTO sku VALUES('ticket-2019-07-29','','','2019-07-29 16:30:00',0,1700);
INSERT INTO sku VALUES('ticket-2019-08-05','','','2019-08-05 16:30:00',0,1700);
INSERT INTO sku VALUES('ticket-2019-08-12','','','2019-08-12 16:30:00',0,1700);
INSERT INTO sku VALUES('summer-sings-2019','','','2019-08-12 16:30:00',0,8500);
INSERT INTO sku VALUES('summer-sings-2019-student','','','2018-12-31 23:00:00',0,0);
INSERT INTO sku VALUES('recording-2018-11-03','','','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES('recording-2018-12-16','','','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES('recording-2019-03-16','','2019-06-01 00:00:00','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES('recording-2019-05-24','','2019-06-01 00:00:00','2019-08-01 00:00:00',1,2000);
INSERT INTO event VALUES('2019-07-08',NULL,'Summer Sing','2019-07-08 19:30:00',0);
INSERT INTO event VALUES('2019-07-15',NULL,'Summer Sing','2019-07-15 19:30:00',0);
INSERT INTO event VALUES('2019-07-22',NULL,'Summer Sing','2019-07-22 19:30:00',0);
INSERT INTO event VALUES('2019-07-29',NULL,'Summer Sing','2019-07-29 19:30:00',0);
INSERT INTO event VALUES('2019-08-05',NULL,'Summer Sing','2019-08-05 19:30:00',0);
INSERT INTO event VALUES('2019-08-12',NULL,'Summer Sing','2019-08-12 19:30:00',0);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-07-08',0);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-07-15',30);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-07-22',30);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-07-29',30);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-08-05',30);
INSERT INTO product_event VALUES('ticket-2019-07-08','2019-08-12',30);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-07-08',30);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-07-15',0);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-07-22',30);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-07-29',30);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-08-05',30);
INSERT INTO product_event VALUES('ticket-2019-07-15','2019-08-12',30);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-07-08',30);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-07-15',30);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-07-22',0);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-07-29',30);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-08-05',30);
INSERT INTO product_event VALUES('ticket-2019-07-22','2019-08-12',30);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-07-08',30);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-07-15',30);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-07-22',30);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-07-29',0);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-08-05',30);
INSERT INTO product_event VALUES('ticket-2019-07-29','2019-08-12',30);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-07-08',30);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-07-15',30);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-07-22',30);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-07-29',30);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-08-05',0);
INSERT INTO product_event VALUES('ticket-2019-08-05','2019-08-12',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-07-08',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-07-15',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-07-22',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-07-29',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-08-05',30);
INSERT INTO product_event VALUES('ticket-2019-08-12','2019-08-12',0);
INSERT INTO product_event VALUES('summer-sings-2019','2019-07-08',20);
INSERT INTO product_event VALUES('summer-sings-2019','2019-07-15',20);
INSERT INTO product_event VALUES('summer-sings-2019','2019-07-22',20);
INSERT INTO product_event VALUES('summer-sings-2019','2019-07-29',20);
INSERT INTO product_event VALUES('summer-sings-2019','2019-08-05',20);
INSERT INTO product_event VALUES('summer-sings-2019','2019-08-12',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-07-08',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-07-15',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-07-22',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-07-29',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-08-05',20);
INSERT INTO product_event VALUES('summer-sings-2019-student','2019-08-12',20);
COMMIT;