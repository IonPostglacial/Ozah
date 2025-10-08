package picture

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"nicolas.galipot.net/hazo/storage/dataset"
	"nicolas.galipot.net/hazo/storage/dsdb"
	"nicolas.galipot.net/hazo/user"
)

type UploadResult struct {
	AttachmentIndex int
	OriginalPath    string
	SmallPath       string
	MediumPath      string
	BigPath         string
}

func UploadPicture(
	ctx context.Context,
	u *user.T,
	dsName string,
	docRef string,
	file io.Reader,
	filename string,
	attachmentIndex int,
) (*UploadResult, error) {
	ds, err := u.GetDataset(dsName)
	if err != nil {
		return nil, fmt.Errorf("could not access dataset: %w", err)
	}

	queries, err := dataset.OpenDb(ds)
	if err != nil {
		return nil, fmt.Errorf("could not open dataset: %w", err)
	}

	index := attachmentIndex
	if index < 0 {
		attachments, err := queries.GetDocumentAttachments(ctx, docRef)
		if err != nil {
			return nil, fmt.Errorf("could not get existing attachments: %w", err)
		}
		index = len(attachments)
	}

	privateDir := filepath.Dir(string(ds))

	ext := filepath.Ext(filename)

	originalPath := dataset.GetAttachmentPath(privateDir, docRef, index, ext)
	smallPath := dataset.GetThumbnailPath(privateDir, docRef, index, dataset.SizeSmall)
	mediumPath := dataset.GetThumbnailPath(privateDir, docRef, index, dataset.SizeMedium)
	bigPath := dataset.GetThumbnailPath(privateDir, docRef, index, dataset.SizeBig)

	if err := saveUploadedFile(file, originalPath); err != nil {
		return nil, fmt.Errorf("could not save uploaded file: %w", err)
	}

	if err := dataset.CopyOrGenerateThumbnail(originalPath, smallPath, dataset.SizeSmall); err != nil {
		return nil, fmt.Errorf("could not generate small thumbnail: %w", err)
	}

	if err := dataset.CopyOrGenerateThumbnail(originalPath, mediumPath, dataset.SizeMedium); err != nil {
		return nil, fmt.Errorf("could not generate medium thumbnail: %w", err)
	}

	if err := dataset.CopyOrGenerateThumbnail(originalPath, bigPath, dataset.SizeBig); err != nil {
		return nil, fmt.Errorf("could not generate big thumbnail: %w", err)
	}

	_, err = queries.InsertDocumentAttachment(ctx, dsdb.InsertDocumentAttachmentParams{
		DocumentRef:     docRef,
		AttachmentIndex: int64(index),
		Source:          filename,
		Path:            originalPath,
		PathSmall:       smallPath,
		PathMedium:      mediumPath,
		PathBig:         bigPath,
	})
	if err != nil {
		return nil, fmt.Errorf("could not insert picture into database: %w", err)
	}

	return &UploadResult{
		AttachmentIndex: index,
		OriginalPath:    originalPath,
		SmallPath:       smallPath,
		MediumPath:      mediumPath,
		BigPath:         bigPath,
	}, nil
}

func saveUploadedFile(file io.Reader, dst string) error {
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return err
	}

	return nil
}
