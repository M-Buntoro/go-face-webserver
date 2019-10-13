package facenet

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func NewFacenet(modelfile string) (facenet Facenet, err error) {
	var newfc Facenet
	var model []byte

	// if modelfile given is url, fetch and load as byte, else load as file
	if _, err = os.Stat(modelfile); os.IsNotExist(err) {
		// Load from url
		var client http.Client
		resp, err := client.Get(modelfile)
		if err != nil {
			log.Println("[Error][FACENET][INITZ]", err)
			return newfc, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			model, err = ioutil.ReadAll(resp.Body)
		} else {
			err = errors.New("Error Fetching from URL")
			log.Println("[Error][FACENET][INITY]", err)
			return newfc, err
		}
	} else {
		// Load model from file
		model, err = ioutil.ReadFile(modelfile)
		if err != nil {
			log.Println("[Error][FACENET][INITD]", err)
			return newfc, err
		}
	}

	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Println("[Error][FACENET][INITF]", err)
		return newfc, err
	}

	newfc = Facenet{
		FacenetModel:   graph,
		FacenetSession: nil,
	}

	return newfc, nil
}
