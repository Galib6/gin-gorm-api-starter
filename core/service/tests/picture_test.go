package service_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"myapp/core/entity"
	"myapp/core/helper/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func buildFileHeader(t *testing.T, field, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile(field, filename)
	require.NoError(t, err)
	_, err = fw.Write(content)
	require.NoError(t, err)
	require.NoError(t, mw.Close())

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	require.NoError(t, req.ParseMultipartForm(int64(len(content)+1024)))
	return req.MultipartForm.File[field][0]
}

func setupTemporaryFileDir(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	tmp := t.TempDir()
	require.NoError(t, os.Chdir(tmp))

	tmpD, err := os.Getwd()
	require.NoError(t, err)

	return tmpD
}

func TestUserService_ChangePicture(t *testing.T) {
	tmpDir := setupTemporaryFileDir(t)
	us, repo, _, ctx := setupUserServiceMock()

	expectedUser := entity.User{ID: uuid.New(), Name: "P", Email: "p@mail.test"}
	fh := buildFileHeader(t, "picture", "pic.txt", []byte("hello"))
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", expectedUser.ID.String()).Return(expectedUser, nil).Once()
	repo.On("UpdateUser", ctx, (*gorm.DB)(nil), mock.AnythingOfType("entity.User")).Return(nil)

	updated, err := us.ChangePicture(ctx, dto.UserChangePictureRequest{ID: expectedUser.ID.String(), Picture: fh})
	require.NoError(t, err)
	require.NotEmpty(t, updated.Picture)
	expectedFilePath := filepath.Join(tmpDir, "files", updated.Picture)
	require.FileExists(t, expectedFilePath, "file should have been uploaded")
}

func TestUserService_DeletePicture(t *testing.T) {
	tmpDir := setupTemporaryFileDir(t)
	us, repo, _, ctx := setupUserServiceMock()

	picPath := "user_picture/" + uuid.New().String()
	expectedUser := entity.User{
		ID:      uuid.New(),
		Name:    "P",
		Email:   "p@mail.test",
		Picture: &picPath,
	}

	// Create a dummy file to simulate existing picture
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "files", "user_picture"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "files", picPath), []byte("hello"), 0644))

	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", expectedUser.ID.String()).
		Return(expectedUser, nil).Once()
	repo.On("UpdateUser", ctx, (*gorm.DB)(nil), mock.AnythingOfType("entity.User")).
		Return(nil)

	err := us.DeletePicture(ctx, expectedUser.ID.String())
	require.NoError(t, err)

	expectedFilePath := filepath.Join(tmpDir, "files", picPath)
	require.NoFileExists(t, expectedFilePath, "file should have been deleted")
}
