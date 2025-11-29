-- Insert categories for each product using product code to look up product_id

-- Category 1: Clothing
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD001'), (SELECT id FROM categories WHERE code = 'CAT001')),
((SELECT id FROM products WHERE code = 'PROD004'), (SELECT id FROM categories WHERE code = 'CAT001')),
((SELECT id FROM products WHERE code = 'PROD007'), (SELECT id FROM categories WHERE code = 'CAT001'));

-- Category 2: Shoes
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD002'), (SELECT id FROM categories WHERE code = 'CAT002')),
((SELECT id FROM products WHERE code = 'PROD006'), (SELECT id FROM categories WHERE code = 'CAT002'));

-- Category 3: Accessories
INSERT INTO product_categories (product_id, category_id) VALUES-- Category 2: Shoes
((SELECT id FROM products WHERE code = 'PROD003'), (SELECT id FROM categories WHERE code = 'CAT003')),
((SELECT id FROM products WHERE code = 'PROD005'), (SELECT id FROM categories WHERE code = 'CAT003')),
((SELECT id FROM products WHERE code = 'PROD008'), (SELECT id FROM categories WHERE code = 'CAT003'));