PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
INSERT INTO product VALUES ('discount','Discount','','other','',0,'');
INSERT INTO product VALUES ('donation','Donation','Donation','donation','<p>Thank you for your generous donation of ${{ dollars .Price }} to Schola Cantorum, which supports our mission to bring choral music to the community through live performances and educational outreach programs. Schola Cantorum is a 501(c)(3) tax-exempt organization. Our tax ID number is 94-2597822.</p>',0,'');
INSERT INTO product VALUES ('gala-purchase','Gala Purchase','','auctionitem','',0,'');
INSERT INTO product VALUES ('quantity-discount','Quantity Discount','','other','',0,'');
INSERT INTO product VALUES ('recording-2011-11-05','Concert Recording: 2011-11','','recording','',0,'');
INSERT INTO product VALUES ('recording-2011-12-11','Concert Recording: 2011-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2012-03-10','Concert Recording: 2012-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2012-05-19','Concert Recording: 2012-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2012-11-03','Concert Recording: 2012-11','','recording','',0,'');
INSERT INTO product VALUES ('recording-2012-12-09','Concert Recording: 2012-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2013-03-22','Concert Recording: 2013-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2013-05-18','Concert Recording: 2013-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2013-11-02','Concert Recording: 2013-11','','recording','',0,'');
INSERT INTO product VALUES ('recording-2013-12-08','Concert Recording: 2013-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2014-03-15','Concert Recording: 2014-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2014-05-17','Concert Recording: 2014-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2014-11-02','Concert Recording: 2014-11','','recording','',0,'');
INSERT INTO product VALUES ('recording-2014-12-07','Concert Recording: 2014-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2015-03-14','Concert Recording: 2015-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2015-05-16','Concert Recording: 2015-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2015-10-23','Concert Recording: 2015-10','','recording','',0,'');
INSERT INTO product VALUES ('recording-2015-12-06','Concert Recording: 2015-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2016-03-20','Concert Recording: 2016-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2016-05-14','Concert Recording: 2016-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2016-10-15','Concert Recording: 2016-10','','recording','',0,'');
INSERT INTO product VALUES ('recording-2016-12-18','Concert Recording: 2016-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2017-03-18','Concert Recording: 2017-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2017-05-20','Concert Recording: 2017-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2017-10-28','Concert Recording: 2017-10','','recording','',0,'');
INSERT INTO product VALUES ('recording-2017-12-17','Concert Recording: 2017-12','','recording','',0,'');
INSERT INTO product VALUES ('recording-2018-03-18','Concert Recording: 2018-03','','recording','',0,'');
INSERT INTO product VALUES ('recording-2018-05-19','Concert Recording: 2018-05','','recording','',0,'');
INSERT INTO product VALUES ('recording-2018-11-03','Concert Recording: 2018-11 For the Love of Bach','2018-11','recording','<p>We confirm your purchase of an archival recording of "For the Love of Bach", November 2018, for {{ dollars .Price }}. Your recording is available for download from the members web site.</p>',0,'');
INSERT INTO product VALUES ('recording-2018-12-16','Concert Recording: 2018-12 A John Rutter Christmas','2018-12','recording','<p>We confirm your purchase of an archival recording of "A John Rutter Christmas", December 2018, for {{ dollars .Price }}. Your recording is available for download from the members web site.</p>',0,'');
INSERT INTO product VALUES ('recording-2019-03-16','Concert Recording: 2019-03 Carmina Burana','2019-03','recording','<p>We confirm your purchase of an archival recording of "Carmina Burana", March 2019, for {{ dollars .Price }}. Once the recording is made available, you can download it from the members web site.</p>',0,'');
INSERT INTO product VALUES ('recording-2019-05-24','Concert Recording: 2019-05 Ein deutsches Requiem','2019-05','recording','<p>We confirm your purchase of an archival recording of "Ein deutsches Requiem", May 2019, for {{ dollars .Price }}. Once the recording is made available, you can download it from the members web site.</p>',0,'');
INSERT INTO product VALUES ('recording-ceremony-carols','CD: A Ceremony of Carols','','other','',0,'');
INSERT INTO product VALUES ('recording-christmas-reflections','CD: Christmas Reflections','','other','',0,'');
INSERT INTO product VALUES ('recording-jester-hairston','CD: An Evening with Jester Hairston','','other','',0,'');
INSERT INTO product VALUES ('recording-shipping','CD Shipping and Handling','','other','',0,'');
INSERT INTO product VALUES ('score-2019-03-16','Carmina Burana Score','','sheetmusic','',0,'');
INSERT INTO product VALUES ('subscription-2010-11-ss','Season Subscription 2010-11','','ticket','',4,'Senior/Student');
INSERT INTO product VALUES ('subscription-2010-11','Season Subscription 2010-11','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2011-12','Season Subscription 2011-12','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2012-13','Season Subscription 2012-13','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2013-14','Season Subscription 2013-14','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2014-15','Season Subscription 2014-15','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2015-16','Season Subscription 2015-16','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2016-17','Season Subscription 2016-17','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2017-18','Season Subscription 2017-18','','ticket','',4,'');
INSERT INTO product VALUES ('subscription-2018-19','Season Subscription 2018-19','','ticket','',3,'');
INSERT INTO product VALUES ('summer-sings-2010','Summer Sings 2010','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2013','Summer Sings 2013','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2014','Summer Sings 2014','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2015','Summer Sings 2015','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2016','Summer Sings 2016','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2017','Summer Sings 2017','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2018','Summer Sings 2018','','ticket','',1,'');
INSERT INTO product VALUES ('summer-sings-2019-student','2019 Summer Sing Student Entry','Student','ticket','',1,'Student');
INSERT INTO product VALUES ('summer-sings-2019','Summer Sings Flex Pass','Flex Pass','ticket','<p>We confirm your purchase of {{ .Quantity }} Flex Pass{{ if gt .Quantity 1 }}es{{ end }} to the 2019 Summer Sings, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  {{ if gt .Quantity 1 }}Each{{ else }}Your{{ end }} pass is good for six entries to the sings: six people on one night, one person on six nights, or any mixture.  The sings are on Monday nights from July 8 through August 12, 2019 (see <a href="https://scholacantorum.org/summer-sings">schedule</a>).  Each sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',6,'');
INSERT INTO product VALUES ('summer-sings-shirt','Summer Sings T-Shirt','','other','',0,'');
INSERT INTO product VALUES ('ticket-2010-06-27','Event Ticket 2010-06-27','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-07-05','Event Ticket 2010-07-05','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-07-12','Event Ticket 2010-07-12','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-07-19','Event Ticket 2010-07-19','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-07-26','Event Ticket 2010-07-26','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-08-02','Event Ticket 2010-08-02','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-08-09','Event Ticket 2010-08-09','','ticket','',1      ,'');
INSERT INTO product VALUES ('ticket-2010-08-2-st','Event Ticket 2010-08-2','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-10-16-st','Event Ticket 2010-10-16','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-10-16','Event Ticket 2010-10-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-10-17-sr','Event Ticket 2010-10-17','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2010-10-17-ss','Event Ticket 2010-10-17','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2010-10-17-st','Event Ticket 2010-10-17','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-10-17','Event Ticket 2010-10-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-11-06','Event Ticket 2010-11-06','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-12-05-sr','Event Ticket 2010-12-05','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2010-12-05-st','Event Ticket 2010-12-05','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-12-05','Event Ticket 2010-12-05','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-12-11-sr','Event Ticket 2010-12-11','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2010-12-11-ss','Event Ticket 2010-12-11','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2010-12-11-st','Event Ticket 2010-12-11','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-12-11','Event Ticket 2010-12-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2010-12-12-sr','Event Ticket 2010-12-12','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2010-12-12-ss','Event Ticket 2010-12-12','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2010-12-12-st','Event Ticket 2010-12-12','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2010-12-12','Event Ticket 2010-12-12','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-03-12-sr','Event Ticket 2011-03-12','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-03-12-ss','Event Ticket 2011-03-12','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2011-03-12-st','Event Ticket 2011-03-12','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2011-03-12','Event Ticket 2011-03-12','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-03-13-sr','Event Ticket 2011-03-13','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-03-13-ss','Event Ticket 2011-03-13','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2011-03-13','Event Ticket 2011-03-13','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-04-02','Event Ticket 2011-04-02','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-04-10-sr','Event Ticket 2011-04-10','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-04-10-st','Event Ticket 2011-04-10','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2011-04-10','Event Ticket 2011-04-10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-05-21-sr','Event Ticket 2011-05-21','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-05-21-ss','Event Ticket 2011-05-21','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2011-05-21-st','Event Ticket 2011-05-21','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2011-05-21','Event Ticket 2011-05-21','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-05-22-ss','Event Ticket 2011-05-22','','ticket','',1,'Senior/Student');
INSERT INTO product VALUES ('ticket-2011-05-22','Event Ticket 2011-05-22','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-07-11-sr','Event Ticket 2011-07-11','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-07-11','Event Ticket 2011-07-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-07-18-sr','Event Ticket 2011-07-18','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-07-18','Event Ticket 2011-07-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-07-25-sr','Event Ticket 2011-07-25','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-07-25-st','Event Ticket 2011-07-25','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2011-07-25','Event Ticket 2011-07-25','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-08-01','Event Ticket 2011-08-01','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-08-08','Event Ticket 2011-08-08','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-08-1-sr','Event Ticket 2011-08-1','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-08-15-sr','Event Ticket 2011-08-15','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-08-15-st','Event Ticket 2011-08-15','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2011-08-15','Event Ticket 2011-08-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-08-8-sr','Event Ticket 2011-08-8','','ticket','',1,'Senior');
INSERT INTO product VALUES ('ticket-2011-10-15','Event Ticket 2011-10-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-11-05','Event Ticket 2011-11-05','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-11-06','Event Ticket 2011-11-06','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2011-12-11-c','Event Ticket 2011-12-11','','ticket','',1,'Child');
INSERT INTO product VALUES ('ticket-2011-12-11','Event Ticket 2011-12-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-03-10','Event Ticket 2012-03-10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-03-11','Event Ticket 2012-03-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-04-21','Event Ticket 2012-04-21','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-05-19','Event Ticket 2012-05-19','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-05-20','Event Ticket 2012-05-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-07-09','Event Ticket 2012-07-09','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-07-16','Event Ticket 2012-07-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-07-23','Event Ticket 2012-07-23','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-07-30','Event Ticket 2012-07-30','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-08-06','Event Ticket 2012-08-06','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-08-13','Event Ticket 2012-08-13','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-10-13','Event Ticket 2012-10-13','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-11-03-st','Event Ticket 2012-11-03','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2012-11-03','Event Ticket 2012-11-03','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-11-04','Event Ticket 2012-11-04','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-12-09-c','Event Ticket 2012-12-09','','ticket','',1,'Child');
INSERT INTO product VALUES ('ticket-2012-12-09-st','Event Ticket 2012-12-09','','ticket','',1,'Student');
INSERT INTO product VALUES ('ticket-2012-12-09','Event Ticket 2012-12-09','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2012-12-10','Event Ticket 2012-12-10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-03-22','Event Ticket 2013-03-22','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-04-20','Event Ticket 2013-04-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-05-18','Event Ticket 2013-05-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-05-19','Event Ticket 2013-05-19','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-07-08','Event Ticket 2013-07-08','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-07-15','Event Ticket 2013-07-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-07-22','Event Ticket 2013-07-22','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-07-29','Event Ticket 2013-07-29','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-08-05','Event Ticket 2013-08-05','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-08-12','Event Ticket 2013-08-12','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-11-02','Event Ticket 2013-11-02','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-11-03','Event Ticket 2013-11-03','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2013-12-08-c','Event Ticket 2013-12-08','','ticket','',1,'Child');
INSERT INTO product VALUES ('ticket-2013-12-08','Event Ticket 2013-12-08','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-03-15-y','Event Ticket 2014-03-15','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2014-03-15','Event Ticket 2014-03-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-03-16-y','Event Ticket 2014-03-16','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2014-03-16','Event Ticket 2014-03-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-05-03','Event Ticket 2014-05-03','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-05-17-y','Event Ticket 2014-05-17','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2014-05-17','Event Ticket 2014-05-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-07-07','Event Ticket 2014-07-07','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-07-14','Event Ticket 2014-07-14','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-07-21','Event Ticket 2014-07-21','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-07-28','Event Ticket 2014-07-28','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-08-04','Event Ticket 2014-08-04','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-08-11','Event Ticket 2014-08-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-11-02-y','Event Ticket 2014-11-02','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2014-11-02','Event Ticket 2014-11-02','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2014-12-07','Event Ticket 2014-12-07','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-03-14','Event Ticket 2015-03-14','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-03-15','Event Ticket 2015-03-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-04-25','Event Ticket 2015-04-25','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-05-16-y','Event Ticket 2015-05-16','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2015-05-16','Event Ticket 2015-05-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-05-17-y','Event Ticket 2015-05-17','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2015-05-17','Event Ticket 2015-05-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-07-13','Event Ticket 2015-07-13','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-07-20','Event Ticket 2015-07-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-07-27','Event Ticket 2015-07-27','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-08-03','Event Ticket 2015-08-03','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-08-10','Event Ticket 2015-08-10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-08-17','Event Ticket 2015-08-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-10-23-y','Event Ticket 2015-10-23','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2015-10-23','Event Ticket 2015-10-23','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-10-25-y','Event Ticket 2015-10-25','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2015-10-25','Event Ticket 2015-10-25','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-11-22','Event Ticket 2015-11-22','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2015-12-06','Event Ticket 2015-12-06','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-02-14','Event Ticket 2016-02-14','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-03-20','Event Ticket 2016-03-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-04-16','Event Ticket 2016-04-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-05-14','Event Ticket 2016-05-14','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-05-15','Event Ticket 2016-05-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-07-11','Event Ticket 2016-07-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-07-18','Event Ticket 2016-07-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-07-25','Event Ticket 2016-07-25','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-08-01','Event Ticket 2016-08-01','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-08-08','Event Ticket 2016-08-08','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-08-15','Event Ticket 2016-08-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-10-15-y','Event Ticket 2016-10-15','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2016-10-15','Event Ticket 2016-10-15','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-10-16','Event Ticket 2016-10-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-11-20','Event Ticket 2016-11-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-12-03','Event Ticket 2016-12-03','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2016-12-18','Event Ticket 2016-12-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-02-12','Event Ticket 2017-02-12','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-03-18-y','Event Ticket 2017-03-18','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2017-03-18','Event Ticket 2017-03-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-04-01','Event Ticket 2017-04-01','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-05-20','Event Ticket 2017-05-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-05-21-y','Event Ticket 2017-05-21','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2017-05-21','Event Ticket 2017-05-21','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-07-10','Event Ticket 2017-07-10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-07-17','Event Ticket 2017-07-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-07-24','Event Ticket 2017-07-24','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-07-31','Event Ticket 2017-07-31','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-08-07','Event Ticket 2017-08-07','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-08-14','Event Ticket 2017-08-14','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-10-28-y','Event Ticket 2017-10-28','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2017-10-28','Event Ticket 2017-10-28','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-10-29','Event Ticket 2017-10-29','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-11-19','Event Ticket 2017-11-19','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2017-12-17','Event Ticket 2017-12-17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-02-11','Event Ticket 2018-02-11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-03-18-y','Event Ticket 2018-03-18','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2018-03-18','Event Ticket 2018-03-18','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-04-28','Event Ticket 2018-04-28','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-05-19','Event Ticket 2018-05-19','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-05-20-y','Event Ticket 2018-05-20','','ticket','',1,'Youth');
INSERT INTO product VALUES ('ticket-2018-05-20','Event Ticket 2018-05-20','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-06-02','Event Ticket 2018-06-02','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-07-09','Event Ticket 2018-07-09','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-07-16','Event Ticket 2018-07-16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-07-23','Event Ticket 2018-07-23','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-07-30','Event Ticket 2018-07-30','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-08-06','Event Ticket 2018-08-06','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-08-13','Event Ticket 2018-08-13','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-11-03','For the Love of Bach, November 3','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2018-11-11','And the Winner Is, November 11','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2019-02-10','You''re Just in Love, February 10','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2019-03-16','Carmina Burana, March 16','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2019-03-17','Carmina Burana, March 17','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2019-04-26','¡Fiesta! April 26','','ticket','',1,'');
INSERT INTO product VALUES ('ticket-2019-07-08','July 8 Summer Sing','July 8','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 8, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES ('ticket-2019-07-15','July 15 Summer Sing','July 15','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 15, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES ('ticket-2019-07-22','July 22 Summer Sing','July 22','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sing on Monday, July 22, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES ('ticket-2019-07-29','July 29 Summer Sing','July 29','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, July 29, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES ('ticket-2019-08-05','August 5 Summer Sing','August 5','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 5, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO product VALUES ('ticket-2019-08-12','August 12 Summer Sing','August 12','ticket','<p>We confirm your purchase of {{ .Quantity }} ticket{{ if gt .Quantity 1 }}s{{ end }} to the Summer Sings on Monday, August 12, for ${{ dollars .Price }}{{ if gt .Quantity 1 }} each{{ end }}.  The sing starts at 7:30pm at Los Altos United Methodist Church, 655 Magdalena Avenue, Los Altos (see <a href="https://www.google.com/maps/place/Los+Altos+United+Methodist+Church/@37.3604399,-122.1163995,14z/data=!4m13!1m7!3m6!1s0x808fb13b09db205b:0x3cb6a0075024dc76!2s655+Magdalena+Ave,+Los+Altos,+CA+94024!3b1!8m2!3d37.3604399!4d-122.09889!3m4!1s0x808fb13baf46a387:0xcfbef6958c3a62d!8m2!3d37.3604399!4d-122.09889">map</a>). Please bring this email (printed or on your phone) for admission.</p>',1,'');
INSERT INTO sku VALUES ('donation','','','',0,0);
INSERT INTO sku VALUES ('ticket-2019-07-08','','','2019-07-08 22:00:00',0,1700);
INSERT INTO sku VALUES ('ticket-2019-07-15','','','2019-07-15 22:00:00',0,1700);
INSERT INTO sku VALUES ('ticket-2019-07-22','','','2019-07-22 22:00:00',0,1700);
INSERT INTO sku VALUES ('ticket-2019-07-29','','','2019-07-29 22:00:00',0,1700);
INSERT INTO sku VALUES ('ticket-2019-08-05','','','2019-08-05 22:00:00',0,1700);
INSERT INTO sku VALUES ('ticket-2019-08-12','','','2019-08-12 22:00:00',0,1700);
INSERT INTO sku VALUES ('summer-sings-2019','','','2019-08-12 22:00:00',0,8500);
INSERT INTO sku VALUES ('summer-sings-2019-student','','','2019-08-12 22:00:00',4,0);
INSERT INTO sku VALUES ('recording-2018-11-03','','','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES ('recording-2018-12-16','','','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES ('recording-2019-03-16','','2019-06-01 00:00:00','2019-08-01 00:00:00',1,2000);
INSERT INTO sku VALUES ('recording-2019-05-24','','2019-06-01 00:00:00','2019-08-01 00:00:00',1,2000);
INSERT INTO event VALUES ('2019-07-08',NULL,'Summer Sing','2019-07-08 19:30:00',0);
INSERT INTO event VALUES ('2019-07-15',NULL,'Summer Sing','2019-07-15 19:30:00',0);
INSERT INTO event VALUES ('2019-07-22',NULL,'Summer Sing','2019-07-22 19:30:00',0);
INSERT INTO event VALUES ('2019-07-29',NULL,'Summer Sing','2019-07-29 19:30:00',0);
INSERT INTO event VALUES ('2019-08-05',NULL,'Summer Sing','2019-08-05 19:30:00',0);
INSERT INTO event VALUES ('2019-08-12',NULL,'Summer Sing','2019-08-12 19:30:00',0);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-07-08',0);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-07-15',30);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-07-22',30);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-07-29',30);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-08-05',30);
INSERT INTO product_event VALUES ('ticket-2019-07-08','2019-08-12',30);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-07-08',30);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-07-15',0);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-07-22',30);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-07-29',30);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-08-05',30);
INSERT INTO product_event VALUES ('ticket-2019-07-15','2019-08-12',30);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-07-08',30);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-07-15',30);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-07-22',0);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-07-29',30);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-08-05',30);
INSERT INTO product_event VALUES ('ticket-2019-07-22','2019-08-12',30);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-07-08',30);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-07-15',30);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-07-22',30);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-07-29',0);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-08-05',30);
INSERT INTO product_event VALUES ('ticket-2019-07-29','2019-08-12',30);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-07-08',30);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-07-15',30);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-07-22',30);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-07-29',30);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-08-05',0);
INSERT INTO product_event VALUES ('ticket-2019-08-05','2019-08-12',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-07-08',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-07-15',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-07-22',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-07-29',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-08-05',30);
INSERT INTO product_event VALUES ('ticket-2019-08-12','2019-08-12',0);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-07-08',20);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-07-15',20);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-07-22',20);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-07-29',20);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-08-05',20);
INSERT INTO product_event VALUES ('summer-sings-2019','2019-08-12',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-07-08',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-07-15',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-07-22',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-07-29',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-08-05',20);
INSERT INTO product_event VALUES ('summer-sings-2019-student','2019-08-12',20);
COMMIT;
