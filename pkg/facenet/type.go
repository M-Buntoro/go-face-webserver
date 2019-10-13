package facenet

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Facenet struct {
	FacenetModel    *tf.Graph
	NormalizerGraph *tf.Graph
	FacenetSession  *tf.Session
}
