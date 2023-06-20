create database marketspace;

\connect marketspace;

drop table if exists users;
drop table if exists announcements; 
drop table if exists payment_methods;

create table if not exists users (
	id serial primary key,
	username varchar(255) not null,
	email varchar(255) unique not null,
	phone varchar(20) not null,
	password varchar(20) not null,
	photo varchar,
	created_at timestamp not null default now()
);

create table if not exists announcements(
	id serial primary key,
	title varchar(100) not null,
	description varchar(255),
	is_new bool default true,
	price decimal,
	is_exchangeable bool default false,
	is_active bool default true,
	images varchar(255)[] default '{}',
	user_id int not null,
	foreign key(user_id) references users(id) on delete cascade
);
 

create table if not exists payment_methods(
	id serial primary key,
	boleto bool,
	pix bool,
	cash bool,
	credit_card bool,
	bank_deposit bool,
	announcement_id int not null,
	foreign key(announcement_id) references announcements(id) on delete cascade
);

-- INSERT DATA

insert into users 
values 
(1, 'user1','user1@email.com', '8598640532', 'user1senha'),
(2, 'user2','user2@email.com', '8591234566', 'user2senha'),
(3, 'user3','user3@email.com', '8591523689', 'user3senha'),
(4, 'user4','user4@email.com', '8599536996', 'user4senha');


insert into announcements
(id, title, description, is_new, price, is_exchangeable, is_active,images, user_id)
values
(1, 'produto1','description for product1',true,12.05,false,true,'{}',1),
(2, 'produto2','description for product2',false,20.00,true,false,'{"https://doodleipsum.com/700/hand-drawn?i=76435a42d4a61505808419da3822a1e5", "https://doodleipsum.com/700/hand-drawn?i=65a7e087a460aafc7f60acc493c7b406"}',1),
(3, 'produto3','description for product3',true,120.00,false,true,'{}',1),
(4, 'produto4','description for product4',false,80.05,true,true,'{"https://doodleipsum.com/700/hand-drawn?i=65a7e087a460aafc7f60acc493c7b406"}',1),
(5, 'produto5','description for product5',true,35.05,false,true,'{"https://doodleipsum.com/700/outline?i=db7184e46fbfced62ecca99f9e1a86a2"}',1),
(6, 'produto12','description for product12',true,125.00,false,true,'{}',2),
(7, 'produto22','description for product22',false,29.05,true,true,'{}',2),
(8, 'produto23','description for product23',true,500.05,false,true,'{"https://doodleipsum.com/700/outline?i=2c350be916b8b173cd3026cbcdea1acb"}',2);


insert into payment_methods
values
(1,false,true,true,false, true,1),
(2,false,true,true,false, true,2),
(3,true,true,true,true, true,3),
(4,false,true,true,false, true,4),
(5,false,true,true,true, true,5),
(6,false,true,true,false, false,6),
(7,false,true,true,false, false,7),
(8,false,true,true,false, false,8);