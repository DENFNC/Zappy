package category_tests

import (
	"testing"

	proto "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	suite "github.com/DENFNC/Zappy/catalog_service/tests/soute"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*************  ✨ Windsurf Command 🌟  *************/
func TestCreateCategory_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	name := gofakeit.Name()
	parentId := ""
	// parentId := gofakeit.UUID()

	// Сначала зарегистрируем нового клиента
	respClient, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
		Name:     name,
		ParentId: &parentId,
	})
	// Это вспомогательный запрос, поэтому делаем лишь минимальные проверки
	require.NoError(t, err)
	assert.NotEmpty(t, respClient.GetCategoryId())
}

// func TestCreateCategory_Duplicate(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	name := gofakeit.Name()
// 	parentId := ""
// 	// parentId := gofakeit.UUID()

// 	// Сначала зарегистрируем нового клиента
// 	respClient, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 		Name:     name,
// 		ParentId: &parentId,
// 	})
// 	// Это вспомогательный запрос, поэтому делаем лишь минимальные проверки
// 	require.NoError(t, err)
// 	assert.NotEmpty(t, respClient.GetCategoryId())

// 	// Затем зарегистрируем его снова
// 	respClient, err = st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 		Name:     name,
// 		ParentId: &parentId,
// 	})

// 	require.Error(t, err)
// 	assert.Empty(t, respClient.GetCategoryId())
// 	assert.Contains(t, err.Error(), "unique violation")

// 	// require.Error(t, err)
// 	// assert.Contains(t, err.Error(), "unique violation")
// }

// func TestCreateCategory_FailCases(t *testing.T) {
// 	ctx, st := suite.New(t)

// }

// func TestCreateCategory_WithParent(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	// 1. Создаём родительскую категорию
// 	parentName := "parent_" + gofakeit.Word()
// 	parentResp, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 		Name:     parentName,
// 		ParentId: nil, // nil = корневая категория
// 	})
// 	require.NoError(t, err)
// 	parentID := parentResp.GetCategoryId().GetId()

// 	// 2. Создаём дочернюю категорию
// 	childName := "child_" + gofakeit.Word()
// 	childResp, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 		Name:     childName,
// 		ParentId: &parentID, // передаём ID родителя

// 	})
// 	require.NoError(t, err)
// 	assert.NotEmpty(t, childResp.GetCategoryId())
// }

// func TestCreateCategory_FailCases(t *testing.T) {
// 	ctx, st := suite.New(t)

// 	// Создаём валидную родительскую категорию
// 	parentResp, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 		Name:     "parent_" + gofakeit.Word(),
// 		ParentId: nil,
// 	})
// 	require.NoError(t, err)
// 	validParentID := parentResp.GetCategoryId().GetId()

// 	tests := []struct {
// 		name         string
// 		categoryName string
// 		parentId     *string // Используем указатель
// 		expectError  bool
// 		errorCode    codes.Code
// 	}{
// 		{
// 			name:         "Empty name",
// 			categoryName: "",
// 			parentId:     nil,
// 			expectError:  true,
// 			errorCode:    codes.InvalidArgument,
// 		},
// 		{
// 			name:         "Nonexistent parent",
// 			categoryName: "valid_name",
// 			parentId:     pointer.StringPtr("00000000-0000-0000-0000-000000000000"),
// 			expectError:  true,
// 			errorCode:    codes.NotFound,
// 		},
// 		{
// 			name:         "Valid with parent",
// 			categoryName: "valid_child",
// 			parentId:     &validParentID,
// 			expectError:  false,
// 		},
// 		{
// 			name:         "Valid without parent",
// 			categoryName: "valid_root",
// 			parentId:     nil,
// 			expectError:  false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			resp, err := st.CategoryServiceClient.CreateCategory(ctx, &proto.CreateCategoryRequest{
// 				Name:     tt.categoryName,
// 				ParentId: tt.parentId,
// 			})

// 			if tt.expectError {
// 				require.Error(t, err)
// 				if st, ok := status.FromError(err); ok {
// 					assert.Equal(t, tt.errorCode, st.Code())
// 				}
// 			} else {
// 				require.NoError(t, err)
// 				assert.NotEmpty(t, resp.GetCategoryId())
// 			}
// 		})
// 	}
// }
