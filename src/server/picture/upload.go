package picture

import (
	"context"
	"mime/multipart"

	"nicolas.galipot.net/hazo/storage/picture"
	"nicolas.galipot.net/hazo/user"
)

func UploadPictureFromMultipart(
	ctx context.Context,
	u *user.T,
	dsName string,
	docRef string,
	file multipart.File,
	header *multipart.FileHeader,
	attachmentIndex int,
) (*picture.UploadResult, error) {
	return picture.UploadPicture(ctx, u, dsName, docRef, file, header.Filename, attachmentIndex)
}
