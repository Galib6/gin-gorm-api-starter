package service_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/testutil"
	"github.com/zetsux/gin-gorm-clean-starter/testutil/factory"
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

	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	require.NoError(t, req.ParseMultipartForm(int64(len(content)+1024)))
	return req.MultipartForm.File[field][0]
}

func TestUserService_ChangeAndDeletePicture(t *testing.T) {
	// run in temp working directory so util writes to an isolated folder
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	tmp := t.TempDir()
	require.NoError(t, os.Chdir(tmp))

	db := testutil.NewTestDB(t)
	us := service.NewUserService(factory.NewUserRepository(t, db))
	ctx := context.Background()

	created, err := us.CreateNewUser(ctx, dto.UserRegisterRequest{Name: "P", Email: "p@mail.test", Password: "secret"})
	require.NoError(t, err)

	fh := buildFileHeader(t, "picture", "pic.txt", []byte("hello"))

	updated, err := us.ChangePicture(ctx, dto.UserChangePictureRequest{Picture: fh}, created.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updated.Picture)

	err = us.DeletePicture(ctx, created.ID)
	require.NoError(t, err)
}
