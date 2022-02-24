CREATE TABLE IF NOT EXISTS parking (
    _id UUID PRIMARY KEY,
    num INTEGER NOT NULL,
    inparking BOOLEAN,
    remark VARCHAR(60)
);

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(60),
    password VARCHAR(65),
    is_admin BOOLEAN
);

INSERT INTO users (username,password,is_admin) VALUES ('admin','$2a$10$WfAxvN0E59pH/PRB1gIzFO6JfRigvCLR.f6cIwbb4WI0dwYjntF/C',true)