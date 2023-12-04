INSERT INTO users(email, phone_number, hashed_password, first_name, last_name, created_at, updated_at)
VALUES ('test@test.test', '+1234567890', '$2a$10$dbN2LIdj.CRCignp0ePzuu5SpQJOzXR5fHh/HxKQqD.FgiONnc2Hu', 'TestUser', 'TestLastname', current_timestamp, current_timestamp);

INSERT INTO businesses(name, business_name, timezone, created_at, updated_at)
VALUES ('Business One', 'business1', 'America/Toronto', current_timestamp, current_timestamp);

INSERT INTO permissions(action) VALUES ('Add a booking');
INSERT INTO permissions(action) VALUES ('Cancel a booking');
INSERT INTO permissions(action) VALUES ('Modify business settings');
INSERT INTO permissions(action) VALUES ('View all bookings');

INSERT INTO roles(name) VALUES ('owner');

INSERT INTO role_permissions(role_id, permission_id, created_at) VALUES (1, 1, current_timestamp);
INSERT INTO role_permissions(role_id, permission_id, created_at) VALUES (1, 2, current_timestamp);
INSERT INTO role_permissions(role_id, permission_id, created_at) VALUES (1, 3, current_timestamp);
INSERT INTO role_permissions(role_id, permission_id, created_at) VALUES (1, 4, current_timestamp);

INSERT INTO business_users (business_id, user_id, role_id, first_name, last_name, created_at)
VALUES 
(1, 1, 1, 'TestUser/AtBiz', 'TestLastname/AtBiz', current_timestamp);

INSERT INTO business_config (business_id, key, value, created_at, updated_at)
VALUES 
(1, 'ui-config', 
'{
  "primaryColorLight": "#ffd27a",
  "primaryColorDark": "#ffaf14",
  "secondaryColorLight": "#e9edf0",
  "secondaryColorDark": "#ced4da",
  "secondaryColorDarker": "#b1bbc4",
  "secondaryColorDarkest": "#404040",
  "complementaryColorLight": "#eaf0f0",
  "complementaryColorDark": "#45645b",
  "disableColor": "#f3f3f3"
}',current_timestamp, current_timestamp);

INSERT INTO business_config (business_id, key, value, created_at, updated_at)
VALUES 
(1, 'booking-config', 
'{
  "maxFutureBookingDays": 5
}', current_timestamp, current_timestamp);

INSERT INTO business_config (business_id, key, value, created_at, updated_at)
VALUES 
(1, 'business-hours', 
'{
  "monday": { "start": "09:00", "end": "19:00" },
  "tuesday": { "start": "09:00", "end": "19:00" },
  "wednesday": { "start": "09:00", "end": "19:00" },
  "thursday": { "start": "09:00", "end": "19:00" },
  "friday": { "start": "09:00", "end": "19:00" },
  "saturday": { "start": "09:00", "end": "17:00" },
  "sunday": { "closed": true }
}',current_timestamp, current_timestamp);

INSERT INTO service_types (business_id, name, icon, description, position, created_at)
VALUES 
(1, 'Show Room', 'show_room_icon', 'A premium service to make your vehicle look as good as new.', 1, current_timestamp),
(1, 'Basic', 'basic_icon', 'Basic cleaning and maintenance, ideal for quick touch-ups.', 2, current_timestamp),
(1, 'Interior', 'interior_icon', 'Focused on cleaning and sanitizing the vehicle''s interior.', 3, current_timestamp),
(1, 'Exterior', 'exterior_icon', 'Focused on exterior wash and wax, to make your vehicle shine.', 4, current_timestamp);

INSERT INTO vehicle_types (business_id, name, icon, description, position, created_at)
VALUES 
(1, 'Sedan', 'sedan_icon', 'A small to medium-sized vehicle with comfortable seating for 4-5 passengers.', 1, current_timestamp),
(1, 'SUV', 'suv_icon', 'A medium to large-sized vehicle suitable for families, with optional all-wheel drive.', 2, current_timestamp),
(1, 'Large SUV / Truck', 'large_suv_truck_icon', 'A large vehicle with ample cargo space, often used for towing or off-road activities.', 3, current_timestamp),
(1, 'Motorcycle', 'motorcycle_icon', 'A two-wheeler suitable for individual riders or a couple, fuel-efficient and quick.', 4, current_timestamp);

INSERT INTO users (email, phone_number, hashed_password, first_name, last_name)
VALUES ('unknown@customer.com', '0000000000', 'hashed_password_for_unknown', 'Unknown', 'Customer');

INSERT INTO channels (name)
VALUES ('sms'), ('email');

INSERT INTO bookings (business_id, user_id, vehicle_type_id, service_type_id, datetime, cost, discount, deposit, bill_number, status, estimated_minutes)
VALUES 
(1, 1, 1, 1, '2023-10-02 10:00:00', 100, 10, 10, 1001, 'pending_payment', 180),
(1, 1, 2, 2, '2023-10-02 14:00:00', 120, 15, 10, 1002, 'pending_payment', 60),
(1, 2, 3, 3, '2023-10-03 09:30:00', 150, 20, 10, 1003, 'pending_payment', 90),
(1, 2, 4, 4, '2023-10-04 15:00:00', 130, 10, 10, 1004, 'active', 120),
(1, 2, 2, 1, '2023-10-05 13:00:00', 110, 5,  10, 1005, 'cancelled', 180),
(1, 1, 1, 2, '2023-10-06 12:00:00', 90, 10, 10, 1006, 'pending_payment', 60),
(1, 2, 2, 3, '2023-10-07 10:30:00', 130, 15, 10, 1007, 'active', 90),
(1, 1, 3, 1, '2023-10-08 09:00:00', 190, 10, 10, 1008, 'active', 180),
(1, 1, 4, 4, '2023-10-09 16:00:00', 90, 5, 10, 1009, 'cancelled', 120),
(1, 2, 3, 2, '2023-10-10 11:00:00', 140, 15, 10, 1010, 'pending_payment', 60);

INSERT INTO booking_logs (booking_id, user_id, state, details)
VALUES 
(1, 1, 'Booked', 'Owner made a booking, pending payment.'),
(2, 1, 'Booked', 'Owner made another booking, pending payment.'),
(3, 2, 'Booked', 'Unknown customer booking, pending payment.'),
(4, 2, 'Paid', 'Unknown customer made a successful payment.'),
(5, 2, 'Payment Failed', 'Payment failed for this booking.');

INSERT INTO transactions (booking_id, amount, gateway, status, notes)
VALUES 
(1, 90, 'Stripe', 'pending', 'Owner booking 1.'),
(2, 105, 'Stripe', 'pending', 'Owner booking 2.'),
(3, 130, 'Stripe', 'pending', 'Customer booking 3.'),
(4, 120, 'Stripe', 'successful', 'Customer booking 4.'),
(5, 105, 'Stripe', 'failed', 'Customer booking 5.');

INSERT INTO communications (channel_id, user_id, from_address, to_address, content)
VALUES 
(2, 1, 'system@yourdomain.com', 'test@test.test', 'Your booking is pending payment.'),
(1, 1, '+1111111111', '+1234567890', 'Thank you for your booking.'), -- SMS with correct format
(2, 2, 'system@yourdomain.com', 'unknown@customer.com', 'Your booking was successful.'),
(1, 2, '+1111111111', '0000000000', 'Your payment failed. Please retry.'), -- SMS with correct format
(2, 2, 'system@yourdomain.com', 'unknown@customer.com', 'Booking cancelled due to failed payment.');

INSERT INTO contacts (contactable_type, contactable_id, address, phone, email)
VALUES 
('user', 1, '123 Owner St., City', '+1234567890', 'test@test.test'),
('user', 2, 'Unknown Address', '0000000000', 'unknown@customer.com'),
('business', 1, '123 Business St., City', '+1111111111', 'business@yourdomain.com'); -- business entry


-- For Sedan
INSERT INTO service_details (business_id, vehicle_type_id, service_type_id, price, duration_minutes)
VALUES 
(1, 1, 1, 10000, 180),  -- Show Room service for Sedan ($100 in cents)
(1, 1, 2, 2000, 60),   -- Basic service for Sedan ($20 in cents)
(1, 1, 3, 4000, 90),   -- Interior service for Sedan ($40 in cents)
(1, 1, 4, 6000, 120);  -- Exterior service for Sedan ($60 in cents)

-- For SUV
INSERT INTO service_details (business_id, vehicle_type_id, service_type_id, price, duration_minutes)
VALUES 
(1, 2, 1, 15000, 180),  -- Show Room service for SUV ($150 in cents)
(1, 2, 2, 2500, 60),   -- Basic service for SUV ($25 in cents)
(1, 2, 3, 5000, 90),   -- Interior service for SUV ($50 in cents)
(1, 2, 4, 7000, 120);  -- Exterior service for SUV ($70 in cents)

-- For Large SUV / Truck
INSERT INTO service_details (business_id, vehicle_type_id, service_type_id, price, duration_minutes)
VALUES 
(1, 3, 1, 20000, 180),  -- Show Room service for Large SUV / Truck ($200 in cents)
(1, 3, 2, 3000, 60),   -- Basic service for Large SUV / Truck ($30 in cents)
(1, 3, 3, 5500, 90),   -- Interior service for Large SUV / Truck ($55 in cents)
(1, 3, 4, 8000, 120);  -- Exterior service for Large SUV / Truck ($80 in cents)

-- For Motorcycle
INSERT INTO service_details (business_id, vehicle_type_id, service_type_id, price, duration_minutes)
VALUES 
(1, 4, 1, 8000, 180),   -- Show Room service for Motorcycle ($80 in cents)
(1, 4, 2, 1500, 60),   -- Basic service for Motorcycle ($15 in cents)
(1, 4, 3, 2500, 90),   -- Interior service for Motorcycle ($25 in cents)
(1, 4, 4, 4000, 120);  -- Exterior service for Motorcycle ($40 in cents)
