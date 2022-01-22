Step 1: Create DB. Database model (use with Mysql):

create database info;
use info;

create table APIconf (
apikey varchar(100),
secretkey varchar(100),
memo varchar(100) primary key);

create table consults(
symbol varchar(50),
moment varchar(50),
price float, 
op varchar(20),
primary key(symbol, moment)
);


insert into APIconf (memo, apikey, secretkey) values ("nil", "nil", "nil");


Step 2: Replace 'user', 'pasword' and 'port' for your own database credentials and ports on database.


Step 3: Save changes and run


