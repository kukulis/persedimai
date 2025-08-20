create or replace table albums
(
    ID     int primary key auto_increment,
    Title  varchar(64),
    Artist varchar(64),
    Price  decimal(10,4)
);


insert albums ( Title, Artist, Price )
       values
       ( 'Summer', 'John Coltrane', 100.99  ),
       ( 'Winter', 'John Coltrane', 50.99  ),
       ( 'Gabalas', 'Riguzas', 40.99  );