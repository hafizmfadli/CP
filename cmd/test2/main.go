package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type SelectedProduct struct {
	ID int
	quantity int
}

type Cart struct {
	UserID int
	SelectedProducts []SelectedProduct
}

func main(){
	// db, err := openDB("root:debezium@tcp(localhost:3306)/inventory?parseTime=true")
	// if err != nil {
	//  fmt.Println(err)
	//  os.Exit(1)
	// }

	// carts := []Cart {
	// 	{
	// 		UserID: 1001,
	// 		SelectedProducts: []SelectedProduct {
	// 			{
	// 				ID: 101,
	// 				quantity: 1,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		UserID: 1002,
	// 		SelectedProducts: []SelectedProduct {
	// 			{
	// 				ID: 101,
	// 				quantity: 1,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		UserID: 1003,
	// 		SelectedProducts: []SelectedProduct {
	// 			{
	// 				ID: 101,
	// 				quantity: 1,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		UserID: 1004,
	// 		SelectedProducts: []SelectedProduct {
	// 			{
	// 				ID: 101,
	// 				quantity: 1,
	// 			},
	// 		},
	// 	},
	// }
		

}

func purchase(ctx context.Context, db *sql.DB, cart Cart) ([]int64, error) {
	
	// Create a helper function for preparing failure results.
	fail := func(err error) (error) {
		return fmt.Errorf("CreateOrder: %v", err)
	}

	//  Sort data (deadlock prevention)

	 // Get a Tx for making transaction requests.
	 tx, err := db.BeginTx(ctx, nil)
	 if err != nil {
			 return nil, fail(err)
	 }
	 // Defer a rollback in case anything fails.
	 defer tx.Rollback()
	
	var orderIDs []int64
	for _, selectedProduct := range cart.SelectedProducts {
		// Confirm that stock is enough for the order.
		var enough bool
		if err = tx.QueryRowContext(ctx, "SELECT quantity >= ? FROM products_on_hand poh  WHERE poh.product_id = ?", 
			selectedProduct.quantity, selectedProduct.ID).Scan(&enough); err != nil {
			if err == sql.ErrNoRows {
				return nil, fail(fmt.Errorf("no such product"))
			}
			return nil, fail(err)
		}

		if !enough {
			return nil, fail(fmt.Errorf("not enough inventory"))
		}

		// Update the product inventory to remove the quantity in the order.
		_, err = tx.ExecContext(ctx, "UPDATE inventory.products_on_hand SET quantity = quantity - ? WHERE product_id = ?",
        selectedProduct.quantity, selectedProduct.ID)
    if err != nil {
        return nil, fail(err)
    }

		// Create a new row in the order table.
		result, err := tx.ExecContext(ctx, `INSERT INTO inventory.orders (order_date,purchaser,quantity,product_id) VALUES (?,?,?,?)`,
        time.Now(), cart.UserID, selectedProduct.quantity, selectedProduct.ID)
    if err != nil {
        return nil, fail(err)
    }

		// Get the ID of the order item just created.
		orderID, err := result.LastInsertId()
		if err != nil {
				return nil, fail(err)
		}

		orderIDs = append(orderIDs, orderID)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return nil, fail(err)
	}

	return orderIDs, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}