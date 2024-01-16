package documents

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"io/ioutil"
	"time"
)

type Remote struct {
	s3  *s3.S3
	cfg *modules.ObjectStorage
}

func NewRemote(s3 *s3.S3, cfg *modules.ObjectStorage) *Remote {
	return &Remote{
		s3:  s3,
		cfg: cfg,
	}
}

func key(doc dto.Document) string {
	return fmt.Sprintf("%s/%s%s", doc.CreatedAt.Format("2006-01-02"), doc.Path.String(), doc.Extension)
}

func (r *Remote) Upload(ctx context.Context, doc dto.Document) (dto.Document, error) {
	object := s3.PutObjectInput{
		Bucket:      aws.String(r.cfg.Bucket),
		Key:         aws.String(key(doc)),
		Body:        doc.RequestContent,
		ContentType: aws.String(doc.Type),
		ACL:         aws.String("private"),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String(doc.Path.String()),
		},
	}

	logrus.Debugf("[object input]: %+v", object)
	out, err := r.s3.PutObjectWithContext(ctx, &object)
	if err != nil {
		return doc, err
	}
	logrus.Debugf("[object output]: %+v", out)

	return doc, nil
}

func (r *Remote) Get(ctx context.Context, doc dto.Document) (dto.Document, error) {
	object := s3.GetObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(key(doc)),
	}

	logrus.Debugf("[object input]: %+v", object)
	out, err := r.s3.GetObjectWithContext(ctx, &object)
	if err != nil {
		return doc, err
	}
	logrus.Debugf("[object output]: %+v", out)

	doc.ResponseContent, err = ioutil.ReadAll(out.Body)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (r *Remote) Delete(ctx context.Context, doc dto.Document) (dto.Document, error) {
	object := s3.DeleteObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(key(doc)),
	}

	logrus.Debugf("[object input]: %+v", object)
	out, err := r.s3.DeleteObjectWithContext(ctx, &object)
	if err != nil {
		return doc, err
	}
	logrus.Debugf("[object output]: %+v", out)

	if err = r.s3.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(key(doc)),
	}); err != nil {
		return doc, err
	}

	return doc, nil
}

func (r *Remote) Share(ctx context.Context, doc dto.Document, duration time.Duration) (dto.Document, error) {
	object := s3.GetObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(key(doc)),
	}

	logrus.Debugf("[object input]: %+v", object)
	req, out := r.s3.GetObjectRequest(&object)
	link, err := req.Presign(duration)
	if err != nil {
		return doc, err
	}
	logrus.Debugf("[object output]: %+v", out)

	doc.ShareLink = link

	return doc, nil
}
