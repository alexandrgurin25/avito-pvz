package product

import (
	"avito-pvz/internal/entity"
	product "avito-pvz/internal/repository/product/mocks"
	pvz "avito-pvz/internal/repository/pvz/mocks"
	reception "avito-pvz/internal/repository/reception/mocks"
	"context"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestProduct_Add(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	categoryIDInt := 1
	categoryIDString := "1"
	receivingID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"

	expectedProduct := &entity.Product{
		ID:          "1",
		ReceptionID: receivingID,
		Category:    categoryIDString,
	}

	expectedPvz := &entity.PVZ{
		UUID: pvzId,
	}

	mockProductRepo := product.NewMockRepository(ctrl)
	mockPvzRepo := pvz.NewMockRepository(ctrl)
	mockReceptionRepo := reception.NewMockRepository(ctrl)

	mockPvzRepo.EXPECT().
		GetPvzById(gomock.Any(), pvzId).
		Return(expectedPvz, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(gomock.Any(), pvzId).
		Return(&entity.Reception{ID: receivingID, PvzID: pvzId}, nil)

	mockProductRepo.EXPECT().
		GetIdCategoryByName(gomock.Any(), category).
		Return(categoryIDInt, nil)

	mockProductRepo.EXPECT().
		AddProduct(gomock.Any(), receivingID, categoryIDInt).
		Return(expectedProduct, nil)

	productService := New(
		mockProductRepo,
		mockPvzRepo,
		mockReceptionRepo,
	)

	result, err := productService.AddProduct(ctx, category, pvzId)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "1" || result.Category != category || result.ReceptionID != receivingID {
		t.Fatalf("unexpected product: %+v", result)
	}

}
