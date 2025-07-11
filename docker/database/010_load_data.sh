#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  INSERT INTO customers.customers (id, name, sms_number, enabled) VALUES
    ('f0e2d41a-a485-4008-b578-747732ae1089', 'Buyer #1', '555-1212', true),
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'John Smith', '555-1001', true),
    ('b2c3d4e5-f6g7-8901-bcde-f23456789012', 'Sarah Johnson', '555-1002', true),
    ('c3d4e5f6-g7h8-9012-cdef-345678901234', 'Michael Brown', '555-1003', true),
    ('d4e5f6g7-h8i9-0123-defg-456789012345', 'Emily Davis', '555-1004', true),
    ('e5f6g7h8-i9j0-1234-efgh-567890123456', 'David Wilson', '555-1005', true),
    ('f6g7h8i9-j0k1-2345-fghi-678901234567', 'Lisa Anderson', '555-1006', true),
    ('g7h8i9j0-k1l2-3456-ghij-789012345678', 'Robert Taylor', '555-1007', true),
    ('h8i9j0k1-l2m3-4567-hijk-890123456789', 'Jennifer Martinez', '555-1008', true),
    ('i9j0k1l2-m3n4-5678-ijkl-901234567890', 'Christopher Garcia', '555-1009', true),
    ('j0k1l2m3-n4o5-6789-jklm-012345678901', 'Amanda Rodriguez', '555-1010', true),
    ('k1l2m3n4-o5p6-7890-klmn-123456789012', 'James Lopez', '555-1011', true),
    ('l2m3n4o5-p6q7-8901-lmno-234567890123', 'Michelle Gonzalez', '555-1012', true),
    ('m3n4o5p6-q7r8-9012-mnop-345678901234', 'Daniel Perez', '555-1013', true),
    ('n4o5p6q7-r8s9-0123-nopq-456789012345', 'Jessica Torres', '555-1014', true),
    ('o5p6q7r8-s9t0-1234-opqr-567890123456', 'Matthew Flores', '555-1015', true),
    ('p6q7r8s9-t0u1-2345-pqrs-678901234567', 'Nicole Butler', '555-1016', true),
    ('q7r8s9t0-u1v2-3456-qrst-789012345678', 'Andrew Simmons', '555-1017', true),
    ('r8s9t0u1-v2w3-4567-rstu-890123456789', 'Stephanie Foster', '555-1018', true),
    ('s9t0u1v2-w3x4-5678-stuv-901234567890', 'Kevin Gonzales', '555-1019', true),
    ('t0u1v2w3-x4y5-6789-tuvw-012345678901', 'Rachel Collins', '555-1020', true);
EOSQL
