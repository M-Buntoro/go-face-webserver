package facesim

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/m-buntoro/go-face-webserver/pkg/distances"
)

func (f Facesim) HandleFaceSim(req FacesimRequest) (res FacesimResult, err error) {

	// get images
	images, err := f.resource.GetImagesByteFromUrl(req.Images)
	if err != nil {
		log.Println(err)
	}

	// find faces
	resl, err := f.resource.GetFacesFromImages(images)
	if err != nil {
		log.Println(err)
	}

	for key, img := range resl {
		f, err := os.Create(fmt.Sprintf("%d.jpg", key))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		jpeg.Encode(f, img, nil)
	}

	// generate face embeddings
	embeddings, err := f.resource.GetFacesEmbeddings(resl)
	if err != nil {
		log.Println(err)
	}

	var (
		totalSim      float32
		countSim      float32
		mapComparison map[string]float32 = map[string]float32{}
	)

	// count average similarity
	for key, embed := range embeddings {

		for key2, compared := range embeddings {
			// don't compare with self
			if key == key2 {
				continue
			}

			mapkey := fmt.Sprintf("%d-%d", key, key2)
			mapkeyflip := fmt.Sprintf("%d-%d", key2, key)
			_, flipok := mapComparison[mapkeyflip]
			if flipok {
				continue
			}

			sim := distances.CosineSimilarity(embed, compared)
			totalSim += sim
			countSim += 1.0
			mapComparison[mapkey] = sim
		}

	}

	res.Similarity = totalSim / countSim
	res.DetailSimilarity = mapComparison
	if res.Similarity > 0.8 {
		res.Verdict = true
	}

	return
}
