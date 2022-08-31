package client

import (
   "errors"
   "fmt"
   "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/client/mock"
   "newrelic-cloudformation-tagging/internal/client/nerdgraph"
   "newrelic-cloudformation-tagging/internal/model"
   "os"
)

/*
   Contract adherence here

   Do not modify the model!
*/

type IClient interface {
   Create(model *model.Model) error
   Delete(model *model.Model) error
   Update(m *model.Model) error
   Read(m *model.Model) error
   List(m *model.Model) ([]interface{}, error)
}

type GraphqlClient struct {
   client IClient
}

var graphqlClient *GraphqlClient

func NewGraphqlClient() *GraphqlClient {
   if graphqlClient == nil {
      // Set in template.yml for TestEntrypoint
      if _, b := os.LookupEnv("Mock"); b {
         graphqlClient = &GraphqlClient{
            client: mock.NewClient(),
         }
         log.Debugln("NewGraphqlClient: returning Mock client")
      } else {
         graphqlClient = &GraphqlClient{
            client: nerdgraph.NewClient(),
         }
         log.Debugln("NewGraphqlClient: returning NerdGraph client")
      }
   }
   return graphqlClient
}

func (i *GraphqlClient) CreateMutation(model *model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Create(model)

   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Create complete",
         ResourceModel:   model,
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          fmt.Sprintf("Create failed: %s", err.Error()),
      }
   }
   return event, nil
}

func (i *GraphqlClient) DeleteMutation(model *model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Delete(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Delete complete",
      }
   } else {
      fmt.Printf("DeleteMutation: error: %+v", err)
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          err.Error(),
      }
   }
   return event, nil
}

func (i *GraphqlClient) UpdateMutation(model *model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Update(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Update complete",
         ResourceModel:   model,
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          err.Error(),
         ResourceModel:    model,
      }
   }
   return event, nil
}

func (i *GraphqlClient) ReadQuery(model *model.Model) (event handler.ProgressEvent, err error) {
   err = i.client.Read(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "Read complete",
         ResourceModel:   model,
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus:  handler.Failed,
         HandlerErrorCode: errors.Unwrap(err).Error(),
         Message:          err.Error(),
      }
   }
   return event, nil
}

func (i *GraphqlClient) ListQuery(model *model.Model) (event handler.ProgressEvent, err error) {
   r, err := i.client.List(model)
   if err == nil {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         "List complete",
         ResourceModels:  r,
      }
   } else {
      event = handler.ProgressEvent{
         OperationStatus: handler.Success,
         Message:         err.Error(),
         ResourceModels:  []interface{}{},
      }
   }
   return event, nil
}
