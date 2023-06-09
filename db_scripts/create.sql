create database workspace;

drop table if exists users 
drop table if exists announcements 
drop table if exists payment_methods

create table if not exists users (
	id serial primary key,
	username varchar(255) not null,
	email varchar(255) unique not null,
	phone varchar(20) not null,
	created_at timestamp not null default now()
)

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
)
 

create table if not exists payment_methods(
	id serial primary key,
	boleto bool,
	pix bool,
	cash bool,
	credit_card bool,
	bank_deposit bool,
	announcement_id int not null,
	foreign key(announcement_id) references announcements(id) on delete cascade
)

