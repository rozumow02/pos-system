INSERT INTO products (name, sku, barcode, price, stock, is_active)
VALUES
    ('USB-C Cable 1m', 'ELEC-001', '100000000001', 8.50, 30, TRUE),
    ('Wireless Mouse', 'ELEC-002', '100000000002', 18.99, 18, TRUE),
    ('Keyboard Mechanical', 'ELEC-003', '100000000003', 49.90, 10, TRUE),
    ('Phone Charger 20W', 'ELEC-004', '100000000004', 15.75, 6, TRUE),
    ('Bluetooth Speaker Mini', 'ELEC-005', '100000000005', 32.00, 4, TRUE)
ON CONFLICT (sku) DO NOTHING;

INSERT INTO inventory_movements (product_id, change_qty, reason)
SELECT id, stock, 'initial_stock'
FROM products
WHERE sku IN ('ELEC-001', 'ELEC-002', 'ELEC-003', 'ELEC-004', 'ELEC-005')
  AND NOT EXISTS (
      SELECT 1
      FROM inventory_movements im
      WHERE im.product_id = products.id
        AND im.reason = 'initial_stock'
  );
