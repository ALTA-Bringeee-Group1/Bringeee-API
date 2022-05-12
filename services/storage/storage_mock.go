package storage

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	Mock *mock.Mock
}

func NewStorageMock(mock *mock.Mock) *StorageMock {
	return &StorageMock{
		Mock: mock,
	}
}

/*
 * Upload from request
 * -------------------------------
 * Upload ke storage service dengan file yang bersumber dari request
 *
 * @param 	fileNamePath 	nama file beserta path yang berada pada cloud storage service
 * @param 	fileHeader	 	file yang akan di upload
 * @return 	string			fileUrl hasil kembalian dari hasil upload
 * @return 	error			error
 */
func (storage StorageMock) UploadFromRequest(fileNamePath string, fileHeader *multipart.FileHeader) (string, error) {
	args := storage.Mock.Called()
	return args.String(0), args.Error(1)
}

/*
 * Delete file
 * -------------------------------
 * Hapus file yang berada pada cloud storage service
 *
 * @param 	fileNamePath 	nama file beserta path yang berada pada cloud storage service
 * @return 	string			fileUrl hasil kembalian dari hasil upload
 * @return 	error			error
 */
func (storage StorageMock) Delete(fileNamePath string) error {
	args := storage.Mock.Called(fileNamePath)
	return args.Error(0)
}