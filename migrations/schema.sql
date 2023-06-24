-- Active: 1687509075746@@127.0.0.1@5432@bookings
select  * from rooms;
select * from room_restrictions;
select * from reservations;

select r.id, r.room_name from rooms r
	where r.id not in 
    (select room_id from room_restrictions rr 
    where make_date(2023,06,1) <= rr.end_date 
    and make_date(2023,06,1) >= rr.start_date);


select * from room_restrictions;

insert into room_restrictions (start_date,end_date,room_id,reservation_id,restriction_id,created_at,updated_at)
VALUES ('2023-06-01','2023-06-01','2',NULL,'1',make_date(2023,6,24),make_date(2023,6,24));

select id,room_name,created_at,updated_at from rooms where id = 1;

SELECT * FROM rooms;

SELECT CURRENT_DATE;

INSERT INTO rooms (room_name,created_at,updated_at) VALUES 
('General''s Quarters', CURRENT_DATE, CURRENT_DATE);

DELETE FROM rooms WHERE id = 4;



SELECT * from restrictions;