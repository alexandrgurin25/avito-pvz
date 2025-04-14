package product

// func TestAddProduct(t *testing.T) {

// 	mockPool, err :=
// 	repo := NewRepository(mockPool)

// 	receivingID := "123"
// 	categoryID := 1
// 	expectedProduct := &entity.Product{
// 		ID:          "1",
// 		ReceptionID: receivingID,
// 		Category:    "1",
// 		DateTime:    time.Now(),
// 	}

// 	// Настройка ожидания для вставки продукта
// 	mockPool.EXPECT().
// 		QueryRow(gomock.Any(), "INSERT INTO products (receiving_id, category_id) VALUES ($1, $2) RETURNING id, receiving_id, category_id, added_at", receivingID, categoryID).
// 		Return(pgxmock.NewRow("id", "receiving_id", "category_id", "added_at").AddRow(expectedProduct.ID, expectedProduct.ReceptionID, expectedProduct.Category, expectedProduct.DateTime))

// 	product, err := repo.AddProduct(context.Background(), receivingID, categoryID)

// 	require.NoError(t, err)
// 	require.Equal(t, expectedProduct, product)
// }

// func TestGetLastProductByReceigingId(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewRepository(mock)

// 	receptionID := "123"
// 	expectedProduct := &entity.Product{
// 		ID:          "1",
// 		ReceptionID: receptionID,
// 		Category:    "1",
// 		DateTime:    time.Now(),
// 	}

// 	mock.ExpectQuery("SELECT id, receiving_id, category_id, added_at FROM products").
// 		WithArgs(receptionID).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "receiving_id", "category_id", "added_at"}).
// 			AddRow(expectedProduct.ID, expectedProduct.ReceptionID, expectedProduct.Category, expectedProduct.DateTime))

// 	product, err := repo.GetLastProductByReceigingId(context.Background(), receptionID)

// 	require.NoError(t, err)
// 	require.Equal(t, expectedProduct, product)
// 	require.NoError(t, mock.ExpectationsWereMet())
// }
