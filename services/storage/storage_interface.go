package storage

import "mime/multipart"



type StorageInterface interface {
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
	UploadFromRequest(fileNamePath string, fileHeader multipart.File) (string, error)

	/*
	 * Delete file
	 * -------------------------------
	 * Hapus file yang berada pada cloud storage service
	 *
	 * @param 	fileNamePath 	nama file beserta path yang berada pada cloud storage service
	 * @return 	string			fileUrl hasil kembalian dari hasil upload
	 * @return 	error			error
	 */
	Delete(fileNamePath string) error
}