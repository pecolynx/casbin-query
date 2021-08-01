create table `pet` (
 `id` integer primary key autoincrement
,`version` int not null default 1
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`name` varchar(20) not null
,unique(`name`)
);
