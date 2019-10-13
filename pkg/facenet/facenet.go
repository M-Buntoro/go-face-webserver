package facenet

import (
	"errors"
	"image"
	"log"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	op "github.com/tensorflow/tensorflow/tensorflow/go/op"
	"gocv.io/x/gocv"
)

func (fnet Facenet) Close(self *Facenet) error {
	if self.FacenetSession != nil {
		return fnet.FacenetSession.Close()
	}
	return nil
}

func (fnet Facenet) BeginSession(self *Facenet) (err error) {
	// if session  is not close, do close before creating new session
	if self.FacenetSession != nil {
		err := self.FacenetSession.Close()
		if err != nil {
			log.Printf("[Error][FACENET] %v", err)
			return err
		}
	}

	self.FacenetSession, err = tf.NewSession(self.FacenetModel, nil)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return err
	}

	return nil
}

// run the image through the model to generate embedding
func (fnet Facenet) ExtractFeatures(tensor *tf.Tensor) ([]float32, error) {
	var embeddingarr [][]float32
	if fnet.FacenetSession == nil {
		return nil, errors.New("Session not started")
	}

	t, err := tf.NewTensor(bool(false))
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return nil, err
	}

	output, err := fnet.FacenetSession.Run(
		map[tf.Output]*tf.Tensor{
			fnet.FacenetModel.Operation("input").Output(0):       tensor,
			fnet.FacenetModel.Operation("phase_train").Output(0): t,
		},
		[]tf.Output{
			fnet.FacenetModel.Operation("embeddings").Output(0),
		},
		nil)
	if err != nil {
		log.Printf("[Error][FACENET] Failed to run inference")
		return []float32{}, err
	}

	var ok bool
	embeddingarr, ok = output[0].Value().([][]float32)
	if !ok {
		log.Printf("[Error][FACENET] Failed convert embedding array")
		return []float32{}, errors.New("Failed convert embedding array")
	}

	if len(embeddingarr) >= 1 {
		return embeddingarr[0], nil
	} else {
		return nil, errors.New("Invalid embedding result")
	}
}

// prepares a gocv.Mat with preprocessing before generating embedding, return a tensor containing processed image
func (fnet Facenet) PrepareMat(img gocv.Mat) (output *tf.Tensor, err error) {
	// convert to tensor
	output, err = TensorFromMat(img)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	// normalize
	output, err = fnet.Normalize(output)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	// prewhiten
	arrayfloat, ok := output.Value().([][][][]float32)
	if !ok {
		return output, errors.New("Image tensor format incorrect, check input")
	}
	mean, stddev := CountMean(arrayfloat[0])
	output, err = fnet.prewhitenImage(output, mean, stddev)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	return output, nil
}

// prepares image with preprocessing before generating embedding, return a tensor containing processed image
func (fnet Facenet) Preprocess(img image.Image) (output *tf.Tensor, err error) {
	// convert to tensor
	output, err = TensorFromImg(img)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	// normalize
	output, err = fnet.Normalize(output)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	// prewhiten
	arrayfloat, ok := output.Value().([][][][]float32)
	if !ok {
		return output, errors.New("Image tensor format incorrect, check input")
	}
	mean, stddev := CountMean(arrayfloat[0])
	output, err = fnet.prewhitenImage(output, mean, stddev)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return output, err
	}

	return output, nil
}

// Creates a graph to decode, rezise and normalize an image
func (fnet Facenet) Normalize(tensor *tf.Tensor) (tensorout *tf.Tensor, err error) {
	s := op.NewScope()
	input := op.Placeholder(s, tf.String)
	// 3 return RGB image

	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	// Sub: returns x - y element-wise
	output := op.ResizeBilinear(s,
		op.ExpandDims(s,
			// cast image to float type
			op.Cast(s, decode, tf.Float),
			op.Const(s.SubScope("make_batch"), int32(0))),
		// resize to 160x160 (model specific)
		op.Const(s.SubScope("size"), []int32{160, 160}))
	graph, err := s.Finalize()
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return nil, err
	}
	defer session.Close()

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{
			input: tensor,
		},
		[]tf.Output{
			output,
		},
		nil)
	if err != nil {
		log.Printf("[Error][FACENET] %v", err)
		return nil, err
	}

	if len(normalized) > 0 {
		tensorout = normalized[0]
	}

	return tensorout, nil

}

// func (fnet *Facenet) prewhitenImage(img *tf.Tensor, mean, std float32) (*tf.Tensor, error) {
func (fnet Facenet) prewhitenImage(img *tf.Tensor, mean, std float32) (*tf.Tensor, error) {
	s := op.NewScope()
	pimg := op.Placeholder(s, tf.Float, op.PlaceholderShape(tf.MakeShape(1, -1, -1, 3)))

	out := op.Mul(s, op.Sub(s, pimg, op.Const(s.SubScope("mean"), mean)), op.Const(s.SubScope("scale"), float32(1.0)/std))

	graph, err := s.Finalize()
	if err != nil {
		log.Println("[Error][FACENET]Error on setting up prewhiten net")
		return img, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Println("[Error][FACENET]Error on setting up prewhiten net")
		return img, err
	}
	defer session.Close()

	outs, err := session.Run(
		map[tf.Output]*tf.Tensor{
			pimg: img,
		},
		[]tf.Output{
			out,
		},
		nil)
	if err != nil {
		log.Println("[Error][FACENET]Error on prewhitening image")
		return img, err
	}

	if len(outs) >= 1 {
		return outs[0], nil
	} else {
		return nil, errors.New("Invalid tensor output from prewhiten")
	}

}
