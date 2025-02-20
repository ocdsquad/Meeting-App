CREATE TYPE language_enum AS ENUM ('English', 'Indonesia');

CREATE TYPE reservation_status_enum AS ENUM ('booked', 'paid', 'cancel');

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       is_active BOOLEAN NOT null default true,
                       language language_enum DEFAULT 'English',
                       is_admin BOOLEAN DEFAULT false,
                       avatar_url VARCHAR(255) DEFAULT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT NULL,
                       deleted_at TIMESTAMP DEFAULT NULL,
                       created_by VARCHAR(255) DEFAULT NULL,
                       updated_by VARCHAR(255) DEFAULT NULL,
                       deleted_by VARCHAR(255) DEFAULT NULL
);

CREATE TABLE category_snacks (
                                 id SERIAL PRIMARY KEY,
                                 name VARCHAR(255) NOT NULL,
                                 price DECIMAL(10, 2) NOT NULL,
                                 currency CHAR(3) NOT NULL,
                                 uom VARCHAR(255) NOT NULL,
                                 is_active BOOLEAN DEFAULT true,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT NULL,
                                 deleted_at TIMESTAMP DEFAULT NULL,
                                 created_by varchar(255) DEFAULT NULL,
                                 updated_by varchar(255) DEFAULT NULL,
                                 deleted_by varchar(255) DEFAULT NULL
);

CREATE TABLE capacities (
                            id SERIAL PRIMARY KEY,
                            value_minimum INT NOT NULL,
                            value_maximum INT NOT NULL,
                            uom VARCHAR(255) NOT NULL,
                            is_active BOOLEAN DEFAULT true,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT NULL,
                            deleted_at TIMESTAMP DEFAULT NULL,
                            created_by varchar(255) DEFAULT NULL,
                            updated_by varchar(255) DEFAULT NULL,
                            deleted_by varchar(255) DEFAULT NULL
);

CREATE TABLE room_types (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            is_active BOOLEAN DEFAULT true,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT NULL,
                            deleted_at TIMESTAMP DEFAULT NULL,
                            created_by varchar(255) DEFAULT NULL,
                            updated_by varchar(255) DEFAULT NULL,
                            deleted_by varchar(255) DEFAULT NULL
);

CREATE TABLE rooms (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       is_active BOOLEAN DEFAULT true,
                       description TEXT DEFAULT NULL,
                       price_hour DECIMAL(10, 2) DEFAULT 0,
                       room_type_id INT NOT NULL,
                       capacity_id INT NOT NULL,
                       capacity INT NOT NULL,
                       attachment_url VARCHAR(500) DEFAULT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT NULL,
                       deleted_at TIMESTAMP DEFAULT NULL,
                       created_by VARCHAR(255) DEFAULT NULL,
                       updated_by VARCHAR(255) DEFAULT NULL,
                       deleted_by VARCHAR(255) DEFAULT NULL,
                       CONSTRAINT fk_rooms_capacities FOREIGN KEY (capacity_id)
                           REFERENCES capacities(id)
                           ON DELETE CASCADE
                           ON UPDATE CASCADE,
                       CONSTRAINT fk_rooms_room_types FOREIGN KEY (room_type_id)
                           REFERENCES room_types(id)
                           ON DELETE CASCADE
                           ON UPDATE CASCADE
);


CREATE TABLE attachments (
                             id SERIAL PRIMARY KEY,
                             file_name VARCHAR(500) NOT NULL,
                             file_size int NOT NULL,
                             file_type VARCHAR(50) NOT NULL,
                             file_path VARCHAR(50) NOT NULL,
                             attachable_id INT NOT NULL,
                             attachable_type VARCHAR(100) NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT NULL,
                             deleted_at TIMESTAMP DEFAULT NULL,
                             created_by varchar(255) DEFAULT NULL,
                             updated_by varchar(255) DEFAULT NULL,
                             deleted_by varchar(255) DEFAULT NULL
);

CREATE TABLE reservation_rooms (
                                   id SERIAL PRIMARY KEY,
                                   room_id INT NOT NULL,
                                   user_id INT NOT NULL,
                                   status reservation_status_enum DEFAULT 'booked',
                                   category_snack_id INT DEFAULT NULL,
                                   name VARCHAR(255) NOT NULL,
                                   date DATE NOT NULL,
                                   start_time TIME NOT NULL,
                                   end_time TIME NOT NULL,
                                   phone VARCHAR(255) DEFAULT NULL,
                                   total_participant INT DEFAULT 0,
                                   organization VARCHAR(255) NOT NULL,
                                   notes TEXT DEFAULT NULL,
                                   total_duration INT DEFAULT 0,
                                   grand_total DECIMAL(10, 2) DEFAULT 0.00,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMP DEFAULT NULL,
                                   deleted_at TIMESTAMP DEFAULT NULL,
                                   created_by VARCHAR(255) DEFAULT NULL,
                                   updated_by VARCHAR(255) DEFAULT NULL,
                                   deleted_by VARCHAR(255) DEFAULT NULL,
                                   CONSTRAINT fk_reservation_rooms_rooms FOREIGN KEY (room_id)
                                       REFERENCES rooms(id)
                                       ON DELETE CASCADE
                                       ON UPDATE CASCADE,
                                   CONSTRAINT fk_reservation_rooms_category_snacks FOREIGN KEY (category_snack_id)
                                       REFERENCES category_snacks(id)
                                       ON DELETE CASCADE
                                       ON UPDATE CASCADE
);

CREATE INDEX idx_rooms_name ON rooms (name);
