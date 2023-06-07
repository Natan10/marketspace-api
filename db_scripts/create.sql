create database workspace;

drop table if exists users; 
drop table if exists payment_methods; 
drop table if exists announcements; 

create table if not exists users (
	id serial primary key,
	username varchar(255) not null,
	email varchar(255) unique not null,
	phone varchar(20) not null,
	created_at timestamp not null default now()
);

create table if not exists payment_methods(
	id serial primary key,
	boleto bool default false,
	pix bool default false,
	cash bool default false,
	credit_card bool default false,
	bank_deposit bool default false
);

create table if not exists announcements(
	id serial primary key,
	title varchar(100) not null,
	description varchar(255),
	is_new bool default true,
	price decimal,
	is_exchangeable bool default false,
	is_active bool default true,
	users_id int not null,
	payment_methods_id int not null,
	
	foreign key(users_id) references users(id) on delete cascade,
	foreign key(payment_methods_id) references payment_methods(id)
);

