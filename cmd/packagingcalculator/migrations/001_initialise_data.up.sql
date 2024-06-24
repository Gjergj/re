CREATE TABLE IF NOT EXISTS products (
    id int auto_increment,
    `date` datetime,
    pack_sizes JSON,
    active int,
    PRIMARY KEY(id)
);

insert into products (id, `date`, pack_sizes, active ) values (1, now(), '[250,500,1000,2000,5000]', 1);