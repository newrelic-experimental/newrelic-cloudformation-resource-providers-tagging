package resource

import (
   "fmt"
   "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/client"
   "newrelic-cloudformation-tagging/internal/model"
   "newrelic-cloudformation-tagging/internal/utils"
   "os"
   "strings"
)

func init() {
   lvl, ok := os.LookupEnv("LOG_LEVEL")
   // LOG_LEVEL not set, let's default to INFO
   if !ok {
      lvl = "info"
   }
   // parse string, this is built-in feature of logrus
   ll, err := log.ParseLevel(lvl)
   if err != nil {
      ll = log.DebugLevel
   }
   // set global log level
   log.SetLevel(ll)
   log.SetFormatter(&log.TextFormatter{ForceQuote: false, DisableQuote: true})
}

// TODO Config
const callbackDelaySeconds = 2

// TODO timeout the goroutines with channels

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
   fmt.Println("")
   utils.DumpModel(currentModel, "Create: In: currentModel")
   utils.DumpModel(prevModel, "Create: In: prevModel")
   setEnv(currentModel)
   // Get the current application global context and test to see if this request is IN_PROGRESS
   pc := progressContextInstance()
   if event, err := pc.getProgressEvent(req); event != nil {
      // Yes, return the current ProgressEvent
      utils.DumpModel(event.ResourceModel, "Create: Out:")
      return *event, err
   }

   // We'll return this now, and save it away as the current default ProgressEvent
   defaultEvent := handler.ProgressEvent{
      OperationStatus:      handler.InProgress,
      Message:              "Create in progress",
      ResourceModel:        shadowModel(currentModel),
      CallbackContext:      createCallbackContext(req),
      CallbackDelaySeconds: callbackDelaySeconds,
   }

   // Spin the work out into a goroutine
   go func() {
      // Save the default Event, used until we FAIL or SUCCESS
      pc.setProgressEvent(req.LogicalResourceID, defaultEvent, nil)
      client := client.NewGraphqlClient()
      sm := shadowModel(currentModel)
      evt, e := client.CreateMutation(sm)
      pc.setProgressEvent(req.LogicalResourceID, evt, e)
      log.Debugf("Create: returning guid: %s", *sm.Guid)
   }()

   fmt.Println("")
   return defaultEvent, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
   fmt.Println("")
   utils.DumpModel(currentModel, "Update: In:")
   setEnv(currentModel)
   // Get the current application global context and test to see if this request is IN_PROGRESS
   pc := progressContextInstance()
   if event, err := pc.getProgressEvent(req); event != nil {
      utils.DumpModel(currentModel, "Update: Out:")
      return *event, err
   }

   // We'll return this now, and save it away as the current default ProgressEvent
   defaultEvent := handler.ProgressEvent{
      OperationStatus:      handler.InProgress,
      Message:              "Update in progress",
      ResourceModel:        shadowModel(currentModel),
      CallbackContext:      createCallbackContext(req),
      CallbackDelaySeconds: callbackDelaySeconds,
   }

   // Spin the work out into a goroutine
   go func() {
      pc.setProgressEvent(req.LogicalResourceID, defaultEvent, nil)
      client := client.NewGraphqlClient()
      sm := shadowModel(currentModel)
      evt, e := client.UpdateMutation(sm)
      pc.setProgressEvent(req.LogicalResourceID, evt, e)
      // log.Debugf("Update: returning guid: %s", *sm.Guid)
   }()

   fmt.Println("")
   return defaultEvent, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
   fmt.Println("")
   utils.DumpModel(currentModel, "Delete: In: currentModel")
   utils.DumpModel(prevModel, "Delete: In: prevModel")
   setEnv(currentModel)
   // Get the current application global context and test to see if this request is IN_PROGRESS
   pc := progressContextInstance()
   if event, err := pc.getProgressEvent(req); event != nil {
      utils.DumpModel(event.ResourceModel, "Delete: Out:")
      return *event, err
   }

   // We'll return this now, and save it away as the current default ProgressEvent
   defaultEvent := handler.ProgressEvent{
      OperationStatus:      handler.InProgress,
      Message:              "delete in progress",
      ResourceModel:        shadowModel(currentModel),
      CallbackContext:      createCallbackContext(req),
      CallbackDelaySeconds: callbackDelaySeconds,
   }

   // Spin the work out into a goroutine
   go func() {
      pc.setProgressEvent(req.LogicalResourceID, defaultEvent, nil)
      client := client.NewGraphqlClient()
      sm := shadowModel(currentModel)
      evt, e := client.DeleteMutation(sm)
      pc.setProgressEvent(req.LogicalResourceID, evt, e)
   }()

   fmt.Println("")
   return defaultEvent, nil
}

//
// Per the contract neither READ nor LIST may return IN_PROGRESS
//

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
   fmt.Println("")
   utils.DumpModel(currentModel, "Read: In: currentModel")
   utils.DumpModel(prevModel, "Read: In: prevModel")
   setEnv(currentModel)
   client := client.NewGraphqlClient()
   sm := shadowModel(currentModel)
   fmt.Println("")
   utils.DumpModel(currentModel, "Read: Out:")
   return client.ReadQuery(sm)
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
   fmt.Println("")
   setEnv(currentModel)
   client := client.NewGraphqlClient()
   sm := shadowModel(currentModel)
   fmt.Println("")
   return client.ListQuery(sm)
}

// UGLY hack to work around circular dependencies
// TODO Abstract
func shadowModel(in *Model) *model.Model {
   out := model.Model{
      Endpoint:        in.Endpoint,
      APIKey:          in.APIKey,
      Guid:            in.EntityGuid,
      EntityGuid:      in.EntityGuid,
      Variables:       in.Variables,
      ListQueryFilter: in.ListQueryFilter,
   }

   tags := make([]model.TagObject, len(in.Tags))
   for i, t := range in.Tags {
      nt := model.TagObject{
         Key:    t.Key,
         Values: t.Values,
      }
      tags[i] = nt
   }
   out.Tags = tags
   return &out
}

// FIXME Ugly hack until we figure-out how to set either the entrypoint or envvars for `test-type` so we don't leak security credentials into AWS
func setEnv(m *Model) {
   if m.APIKey == nil {
      return
   }
   if strings.Contains(strings.ToLower(*m.APIKey), "mockapikey") {
      err := os.Setenv("Mock", "true")
      if err != nil {
         log.Errorf("error setting Mock envvar: %v", err)
      } else {
         log.Traceln(os.Environ())
      }
   }
}
