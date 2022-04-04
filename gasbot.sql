
create database gasbot;

create user gasbot password 'grimble' ;

create table if not exists users (
    id serial  primary key,
    firstname CHARACTER VARYING(20),
    lastname CHARACTER VARYING(20),
    fullname CHARACTER VARYING(40),
    user_id  int unique,
    show_level bool,
    show_info  bool,
    gas_level numeric(8,4),
    version numeric(8,4),
    dateadded timestamptz,
    last_pinged timestamptz,
    level_wait int default 10
);

# alter table users add column last_pinged timestamptz default now();
#  alter table users add column level_wait int default 10;

grant select , insert , update on users to gasbot;

grant select , usage , update on users_id_seq to gasbot;

update users set show_level=false where show_level is null;
update users set show_info= false where show_info is null;


create table if not exists gas_levels (
    id serial not null,
    user_id int unique references users(user_id),
    gas_level numeric(8,4)
);

grant select , insert , update on gas_levels to gasbot;

grant select , usage , update on gas_levels_id_seq to gasbot;


//----
alter table users add column version numeric(8,4);

create table gasrates ( 
    id serial  primary key,
    fastest  NUMERIC(8,4), 
    fast  NUMERIC(8,4), 
    safelow  NUMERIC(8,4),
    average  NUMERIC(8,4),
    fastestwait  NUMERIC(8,4), 
    fastwait  NUMERIC(8,4), 
    safelowwait  NUMERIC(8,4),
    averagewait  NUMERIC(8,4),
    blocknum INT,   
    blocktime NUMERIC(8,4),
    dateadded TIMESTAMPTZ 
 );

 grant select , insert , update on gasrates to gasbot;

grant select , usage , update on gasrates_id_seq to gasbot;
