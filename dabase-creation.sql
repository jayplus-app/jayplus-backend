CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(250) NOT NULL,
  phone_number VARCHAR(250) NOT NULL,
  hashed_password VARCHAR(250) NOT NULL,
  first_name VARCHAR(250) NOT NULL,
  last_name VARCHAR(250) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE permissions (
  id BIGSERIAL PRIMARY KEY,
  action VARCHAR(250) NOT NULL
);

CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(250) NOT NULL
);

CREATE TABLE role_permissions (
  role_id BIGINT NOT NULL REFERENCES roles(id),
  permission_id BIGINT NOT NULL REFERENCES permissions(id),
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE businesses (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(250) NOT NULL,
  business_name VARCHAR(250) NOT NULL,
  timezone VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE business_users (
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  user_id BIGINT NOT NULL REFERENCES users(id),
  role_id BIGINT NOT NULL,
  first_name VARCHAR(250) NOT NULL,
  last_name VARCHAR(250) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE business_config (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  key VARCHAR(250) NOT NULL,
  value JSONB NOT NULL, 
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE vehicle_types (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  name VARCHAR(250) NOT NULL,
  icon VARCHAR(250) NOT NULL,
  description TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE service_types (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  name VARCHAR(250) NOT NULL,
  icon VARCHAR(250) NOT NULL,
  description TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE bookings (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  user_id BIGINT NOT NULL REFERENCES users(id),
  vehicle_type_id BIGINT NOT NULL REFERENCES vehicle_types(id),
  service_type_id BIGINT NOT NULL REFERENCES service_types(id),
  datetime TIMESTAMP NOT NULL,
  estimated_minutes INTEGER NOT NULL,
  cost INTEGER NOT NULL,
  discount INTEGER NOT NULL,
  deposit INTEGER NOT NULL,
  bill_number BIGINT NOT NULL,
  status VARCHAR(200) NOT NULL
);

CREATE TABLE booking_logs (
  id BIGSERIAL PRIMARY KEY,
  booking_id BIGINT NOT NULL REFERENCES bookings(id),
  user_id BIGINT NOT NULL REFERENCES users(id),
  state VARCHAR(250) NOT NULL,
  details TEXT NOT NULL
);

CREATE TABLE transactions (
  id BIGSERIAL PRIMARY KEY,
  booking_id BIGINT NOT NULL REFERENCES bookings(id),
  amount INTEGER NOT NULL,
  gateway VARCHAR(250) NOT NULL,
  status VARCHAR(250) NOT NULL,
  notes TEXT NOT NULL
);

CREATE TABLE channels (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(250) NOT NULL
);

CREATE TABLE communications (
  id BIGSERIAL PRIMARY KEY,
  channel_id BIGINT NOT NULL REFERENCES channels(id),
  user_id BIGINT NOT NULL REFERENCES users(id),
  from_address VARCHAR(250) NOT NULL,
  to_address VARCHAR(250) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE attachments (
  id BIGSERIAL PRIMARY KEY,
  communication_id BIGINT NOT NULL REFERENCES communications(id),
  url VARCHAR(250) NOT NULL
);

CREATE TABLE contacts (
  id BIGSERIAL PRIMARY KEY,
  contactable_type VARCHAR(250) NOT NULL,
  contactable_id BIGINT NOT NULL,
  address VARCHAR(250) NOT NULL,
  phone VARCHAR(200) NOT NULL,
  email VARCHAR(250) NOT NULL
);

CREATE TABLE service_details (
  business_id BIGINT NOT NULL REFERENCES businesses(id),
  vehicle_type_id BIGINT NOT NULL REFERENCES vehicle_types(id),
  service_type_id BIGINT NOT NULL REFERENCES service_types(id),
  price INTEGER NOT NULL,
  duration_minutes INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp,

  UNIQUE(business_id, vehicle_type_id, service_type_id)  -- This ensures that combinations are unique
);