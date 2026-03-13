DELETE FROM inventory_movements
WHERE product_id IN (
    SELECT id
    FROM products
    WHERE sku IN ('ELEC-001', 'ELEC-002', 'ELEC-003', 'ELEC-004', 'ELEC-005')
);

DELETE FROM products
WHERE sku IN ('ELEC-001', 'ELEC-002', 'ELEC-003', 'ELEC-004', 'ELEC-005');
