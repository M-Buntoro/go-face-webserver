package facesim

import (
	"log"

	fresource "github.com/m-buntoro/go-face-webserver/resource/facesim"
)

func New(c Config) (res Facesim) {
	fres, err := fresource.New(
		fresource.Config{
			Conf: c.Conf,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	res = Facesim{
		resource: fres,
	}
	return
}

func (f Facesim) Close() (err error) {
	err = f.resource.Close()
	return err
}
