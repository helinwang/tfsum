Tensorflow provides golang api for model training and inference, however currently tensorboard is only supported when using python.

This repository enables using tensorboard with tensorflow golang api. `tfsum.Writer` writes file that tensorboard could understand.

###Example

```
import (
    "github.com/helinwang/tfsum"
    tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	step := 0
	w := &tfsum.Writer{Dir:"./tf-log", Name:"train"}
	var s *tf.Session
	sum, err := s.Run(
		map[tf.Output]*tf.Tensor{
			s.g.Operation("input").Output(0): put,
		},
		[]tf.Output{
			s.g.Operation("MergeSummary/MergeSummary").Output(0),
		},
		[]*tf.Operation{
			s.g.Operation("train_step"),
		})
	if err != nil {
		fmt.Println(err)
	}
	err = w.AddEvent(sum[0].Value().(string), int64(step))
	if err != nil {
		fmt.Println(err)
	}
}
```
