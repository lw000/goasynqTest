package tasks

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

// A list of tasks types.
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

//----------------------------------------------
// Write a function NewXXXTask to create a tasks.
// A tasks consists of a type and a payload.
//----------------------------------------------

func NewEmailDeliveryTask(userID int, tmplID string) *asynq.Task {
	payload := map[string]interface{}{"user_id": userID, "template_id": tmplID}
	return asynq.NewTask(TypeEmailDelivery, payload)
}

func NewImageResizeTask(src string) *asynq.Task {
	payload := map[string]interface{}{"src": src}
	return asynq.NewTask(TypeImageResize, payload)
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input tasks.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	userID, err := t.Payload.GetInt("user_id")
	if err != nil {
		return err
	}
	tmplID, err := t.Payload.GetString("template_id")
	if err != nil {
		return err
	}
	fmt.Printf("Send Email to User: user_id = %d, template_id = %s\n", userID, tmplID)
	// Email delivery code ...
	return nil
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for struct
}

func (p *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	src, err := t.Payload.GetString("src")
	if err != nil {
		return err
	}
	fmt.Printf("Resize image: src = %s\n", src)
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	// ... return an instance

	return nil
}
