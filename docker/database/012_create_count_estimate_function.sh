#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  -- Create count_estimate function for performance monitoring
  CREATE OR REPLACE FUNCTION count_estimate(query text) RETURNS integer AS \$func\$
  DECLARE
    rec   record;
    rows  integer;
  BEGIN
    FOR rec IN EXECUTE 'EXPLAIN ' || query LOOP
      rows := substring(rec."QUERY PLAN" FROM ' rows=([[:digit:]]+)');
      EXIT WHEN rows IS NOT NULL;
    END LOOP;
    RETURN rows;
  END;
  \$func\$ LANGUAGE plpgsql VOLATILE STRICT;

  -- Grant execute permission to mallbots_user
  GRANT EXECUTE ON FUNCTION count_estimate(text) TO mallbots_user;

  -- Add comment for documentation
  COMMENT ON FUNCTION count_estimate(text) IS 'Estimates the number of rows that would be returned by a query without executing it';
EOSQL


# -- Estimate rows for a simple query
# SELECT count_estimate('SELECT * FROM customers.customers');

# -- Estimate rows for a complex query
# SELECT count_estimate('SELECT * FROM orders.orders WHERE status = ''pending''');

# -- Use in performance monitoring
# SELECT
#   'customers' as table_name,
#   count_estimate('SELECT * FROM customers.customers') as estimated_rows
# UNION ALL
# SELECT
#   'orders' as table_name,
#   count_estimate('SELECT * FROM ordering.orders') as estimated_rows;
