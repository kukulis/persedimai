create table travels (
    travel_id int not null primary key auto_increment,
    travel_departure_point int not null,
    travel_arrival_point int not null,
    travel_departure_time datetime,
    travel_arrival_time datetime,
    travel_price decimal(10,2)
);