package resizer

import (
	"bytes"
	"log"
	"runtime"

	"github.com/disintegration/imaging"
	c "github.com/faroyam/pic-resizer/config"
	"github.com/faroyam/pic-resizer/msg"
)

var queue chan msg.Pic
var cfg *c.Config

func init() {
	queue = make(chan msg.Pic, 100)
	cfg = c.GetConfig()
}

// StartService starts N workers to wait incoming from queue images, where N = NumCPU
func StartService() {
	for w := 0; w < runtime.NumCPU(); w++ {
		go worker(queue)
	}
}

// Send image to queue
func Send(m msg.Pic) {
	go func() { queue <- m }()
}

func worker(queue <-chan msg.Pic) {
	for pic := range queue {
		if err := resize(pic); err != nil {
			log.Println(err)
		}
	}
}

func resize(m msg.Pic) error {

	initImage, err := imaging.Decode(bytes.NewReader(m.Data))
	if err != nil {
		return err
	}
	err = imaging.Save(initImage, cfg.ImagesPath+m.ID+".jpg")
	if err != nil {
		return err
	}

	resizsedImage := imaging.Thumbnail(initImage, 100, 100, imaging.CatmullRom)

	err = imaging.Save(resizsedImage, cfg.ImagesPath+m.ID+cfg.ImagesPreviewSuff+".jpg")
	if err != nil {
		return err
	}

	return nil
}
