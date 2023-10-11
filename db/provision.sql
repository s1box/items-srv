# CREATE DATABASE items;

CREATE TABLE IF NOT EXISTS items (
   id   INT          NOT NULL AUTO_INCREMENT,
   name VARCHAR(255) NOT NULL,
   num  FLOAT,
   PRIMARY KEY (id),
   CONSTRAINT unique_item_name UNIQUE (name)
);

INSERT INTO items ( id, name, num ) VALUES
  ( 1, 'name1', 1.1 ),
  ( 2, 'name2', 2.2 ),
  ( 3, 'name3', 3.3 ),
  ( 4, 'name4', 4.4 ),
  ( 5, 'name5', 5.5 );