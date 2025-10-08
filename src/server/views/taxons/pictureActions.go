package taxons

import (
	"context"
	"mime/multipart"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/picture"
)

type pictureActions struct {
	cc     *common.Context
	dsName string
	docRef string
}

func NewPictureActions(cc *common.Context, dsName, docRef string) *pictureActions {
	return &pictureActions{
		cc:     cc,
		dsName: dsName,
		docRef: docRef,
	}
}

func (h *pictureActions) uploadPicture(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	_, err := picture.UploadPictureFromMultipart(ctx, h.cc.User, h.dsName, h.docRef, file, header, -1)
	return err
}

func (h *pictureActions) Register(reg *action.Registry) {
	reg.AppendAction(action.NewActionWithFileUpload("picture-upload", "picture-file", h.uploadPicture))
}
