package itswizard_m_s3bucket

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"sort"
	"time"
)

type Bucket struct {
	bucketName *string
	region     string
	s3service  *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func CreateNewBucket(bucketName, region string) (err error, result string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return err, ""
	}
	svc := s3.New(sess)
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	res, err := svc.CreateBucket(input)
	if err != nil {
		/*
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeBucketAlreadyExists:
					fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
				case s3.ErrCodeBucketAlreadyOwnedByYou:
					fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
				default:
					return aerr.Error()
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		*/
		return err, ""
	}
	return err, fmt.Sprint(res)
}

func NewBucket(bucketName string, region string) (ret *Bucket, err error) {
	ret = new(Bucket)
	ret.bucketName = aws.String(bucketName)
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}
	ret.uploader = s3manager.NewUploader(sess)
	ret.downloader = s3manager.NewDownloader(sess)
	ret.s3service = s3.New(sess)
	return ret, nil
}

func (p *Bucket) ContentUpload(filedir string, input []byte) error {
	_, err := p.uploader.Upload(&s3manager.UploadInput{Bucket: p.bucketName, Key: aws.String(filedir), Body: bytes.NewReader(input)})
	return err
}

func (p *Bucket) DownloadContent(filedir string) (out []byte, err error) {

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err = p.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: p.bucketName,
		Key:    aws.String(filedir),
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Bucket) listData(path string) (m map[time.Time]string, err error) {
	resp, err := p.s3service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: p.bucketName,
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	m = make(map[time.Time]string)
	for _, item := range resp.Contents {
		m[*item.LastModified] = *item.Key
	}
	return
}

func (p *Bucket) ListAllFiles(path string) (filedir []string, err error) {
	list, err := p.listData(path)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		filedir = append(filedir, v)
	}
	return
}

func (p *Bucket) LastChangedFile(path string) (filedir string, err error) {
	var dateSlice timeSlice = []time.Time{}

	allFiles, err := p.listData(path)
	if err != nil {
		fmt.Println(err)
	}

	for t, _ := range allFiles {
		dateSlice = append(dateSlice, t)
	}

	sort.Sort(sort.Reverse(dateSlice))

	if len(dateSlice) > 0 {
		filedir = allFiles[dateSlice[0]]
	} else {
		err = errors.New(fmt.Sprint("There is no file in the bucket ", p.bucketName))
	}

	return
}

func (p *Bucket) GetLastUploadedContent(path string) (out []byte, err error) {
	filedir, err := p.LastChangedFile(path)
	if err != nil {
		return
	}
	out, err = p.DownloadContent(filedir)
	return
}

/*
Returns true when file exist and false when it does not exist
*/
func (p *Bucket) CheckIfFileExist(filedir string) bool {
	_, err := p.DownloadContent(filedir)
	if err != nil {
		return false
	} else {
		return true
	}
}
