package buckets

import (
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"

	"github.com/hsuanshao/go-tools/buckets/cloud"
	"github.com/hsuanshao/go-tools/buckets/entity"
	ifc "github.com/hsuanshao/go-tools/buckets/interface"
)

// Init to intial bucket service
// NOTE: Init will be replaced by Prepare
func Init(ctx ctx.CTX, conf *entity.Config) ifc.Service {
	if conf == nil {
		ctx.Error("initial bucket service without configuration of blob bucket service")
	}

	return &impl{
		conf: conf,
	}
}

type impl struct {
	conf *entity.Config
}

func (im *impl) GetBucketReader(ctx ctx.CTX) (r ifc.BucketReader, err error) {
	r, err = cloud.NewReader(ctx, im.conf)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err}).Warn("unable init bucket reader by config doc")
		return nil, entity.ErrInitBucketReader
	}
	return r, nil
}

func (im *impl) GetBucketWriter(ctx ctx.CTX) (w ifc.BucketWriter, err error) {
	w, err = cloud.NewWriter(ctx, im.conf)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err}).Warn("unable init bucket writer by config doc")
		return nil, entity.ErrInitBucketWriter
	}
	return w, nil
}
