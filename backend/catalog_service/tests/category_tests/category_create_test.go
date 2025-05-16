package category_tests

// import (
// 	"testing"

// 	proto "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
// 	suite "github.com/DENFNC/Zappy/catalog_service/tests/soute"
// 	"github.com/brianvoe/gofakeit/v6"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateCategory_HappyPath(t *testing.T) {
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
// }
