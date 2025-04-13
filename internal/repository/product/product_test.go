package product_test

// import (
// 	"avito-pvz/internal/entity"
// 	"avito-pvz/internal/repository/product"
// 	"context"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/pashagolub/pgxmock/v2"
// 	"github.com/stretchr/testify/require"
// )

// func TestAddProduct(t *testing.T) {
// 	// 1. Создаём мок pgxpool.Pool
// 	mockPool, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mockPool.Close()

// 	// 2. Создаём репозиторий, передаём мок
// 	repo := product.NewRepository(&pgxpool.Pool{})

// 	// 3. Ожидаем SQL-запрос и возвращаем фейковые данные
// 	receivingID := "123"
// 	categoryID := 1

// 	mockPool.ExpectQuery("INSERT INTO products (.+) VALUES (.+) RETURNING id").
// 		WithArgs(receivingID, categoryID).
// 		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

// 	// 4. Вызываем метод репозитория
// 	result, err := repo.AddProduct(context.Background(), receivingID, categoryID)

// 	// 5. Проверяем:
// 	// - что запрос выполнился без ошибок
// 	// - что SQL-запрос соответствовал ожидаемому
// 	require.NoError(t, err)
// 	require.NoError(t, mockPool.ExpectationsWereMet()) // Проверяем, что все ожидания выполнены
// 	require.Equal(t, &entity.Product{ID: "1"}, result)
// }
