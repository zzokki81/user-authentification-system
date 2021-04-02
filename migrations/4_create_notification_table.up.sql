CREATE TABLE notifications(
  id                      serial primary key,
  recipient_email         varchar(255) not null,
  type                    varchar(50)
);