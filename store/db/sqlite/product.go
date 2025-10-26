package sqlite

import (
	"github.com/thetnaingtn/dirty-hand/store"
	"context"
)

func (d *DB) CreateProduct(ctx context.Context, p *store.Product) (*store.Product, error) {
	res, err := d.db.ExecContext(ctx, `INSERT INTO products (name, description, price, cover, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.Cover, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (d *DB) UpdateProduct(ctx context.Context, p *store.Product) (*store.Product, error) {
	_, err := d.db.ExecContext(ctx, `UPDATE products SET name=?, description=?, price=?, cover=?, updated_at=? WHERE id=?`,
		p.Name, p.Description, p.Price, p.Cover, p.UpdatedAt, p.ID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *DB) ListProducts(ctx context.Context) ([]*store.Product, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT id, name, description, price, cover, created_at, updated_at FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*store.Product
	for rows.Next() {
		var p store.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Cover, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (d *DB) DeleteProduct(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	return err
}
