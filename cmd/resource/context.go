package resource

import (
   "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
   "github.com/aws/aws-sdk-go/service/cloudformation"
   log "github.com/sirupsen/logrus"
   "sync"
)

// Singleton
type progressContext struct {
   mu sync.Mutex
   // Track by AWS Logical Rsource ID
   work map[string]*eventProgress
}

// Current state of a Logical Resource ID
type eventProgress struct {
   event *handler.ProgressEvent
   err   error
}

const callbackKey = "requestPID"

// There can be only one
var progressContextInstancePtr *progressContext

func progressContextInstance() *progressContext {
   if progressContextInstancePtr == nil {
      progressContextInstancePtr = &progressContext{work: make(map[string]*eventProgress)}
   }
   return progressContextInstancePtr
}

func (p *progressContext) getProgressEvent(req handler.Request) (*handler.ProgressEvent, error) {
   // No request context, first time for this request
   if req.CallbackContext == nil {
      return nil, nil
   }

   key, found := req.CallbackContext[callbackKey]
   if found {
      defer p.mu.Unlock()
      p.mu.Lock()
      ep := p.work[key.(string)]
      // eventProgress not found, that's a weird error
      if ep == nil {
         return &handler.ProgressEvent{
            OperationStatus:  handler.Failed,
            HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError,
            Message:          "empty eventProgress in progressContext",
         }, nil
      } else {
         return ep.event, ep.err
      }
   } else {
      // Empty request context, first time for this request
      return nil, nil
   }
}

func (p *progressContext) setProgressEvent(lrid string, event handler.ProgressEvent, err error) {
   log.Printf("setProgressContext: model type: %T", event.ResourceModel)
   defer p.mu.Unlock()
   p.mu.Lock()
   ep := &eventProgress{
      event: &event,
      err:   err,
   }
   p.work[lrid] = ep
}

func createCallbackContext(req handler.Request) map[string]interface{} {
   value := map[string]interface{}{
      callbackKey: req.LogicalResourceID,
   }
   log.Debugf("createCallbackContext: %+v", value)
   return value
}
