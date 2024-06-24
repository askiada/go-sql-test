-- START_TEST
/*
ROW "coucou",true,5,3.14,"{'a': 'b'}"
ROW "coucou2",false,45,18,"{'m': 'n'}"
*/
-- END_TEST
SELECT * FROM table2

-- START_TEST
/*
ROW "coucou",true,5,3.14,"{'a': 'b'}"
ROW "coucou2",false,45,18,"{'m': 'n'}"
*/
-- END_TEST
SELECT * FROM table2
