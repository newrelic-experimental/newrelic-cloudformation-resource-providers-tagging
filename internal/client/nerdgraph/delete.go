package nerdgraph

import (
   "encoding/json"
   "fmt"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/model"
)

type deleteResponse struct {
   Data deleteData `json:"data"`
}
type deleteData struct {
   TaggingDeleteTagFromEntity taggingDeleteTagFromEntity `json:"taggingDeleteTagFromEntity"`
}
type taggingDeleteTagFromEntity struct {
   Errors []errors `json:"errors"`
   Status string   `json:"status"`
}

func (i *nerdgraph) Delete(m *model.Model) error {
   if m == nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, "nil model")
   }
   if m.Guid == nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, "nil guid")
   }

   if m.IsAWSSemantics() {
      // The key must exist to delete it
      err := i.Read(m)
      if err != nil {
         return err
      }
   }

   // 1. Build an array of just the keys
   keys := make([]string, len(m.Tags))
   for _, t := range m.Tags {
      keys = append(keys, *t.Key)
   }
   // 2. JSON Stringify the key array
   ba, err := json.Marshal(keys)
   if err != nil {
      return err
   }
   // 3. Render the key array
   keyString, err := model.Render(string(ba), m.Variables)
   return i.deleteTag(m, keyString)
}

func (i *nerdgraph) deleteTag(m *model.Model, keyString string) error {
   mutation, err := model.Render(deleteMutation, map[string]string{"GUID": *m.Guid, "KEYS": keyString})
   if err != nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   log.Debugf("Delete: guid: %s model: %+v", *m.Guid, *m)
   body, err := i.emit(mutation, *m.APIKey, m.GetEndpoint())
   if err != nil {
      log.Errorf("Delete: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   response := deleteResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   if response.Data.TaggingDeleteTagFromEntity.Errors != nil {
      for _, e := range response.Data.TaggingDeleteTagFromEntity.Errors {
         log.Errorf("Delete: NerdGraph error: %s %s", e.Type, e.Message)
         err = fmt.Errorf("%w %s %s", &model.InvalidRequest{}, e.Type, e.Message)
      }
      return err
   }
   return nil
}

const deleteMutation = `
mutation {
  taggingDeleteTagFromEntity(guid: "{{{GUID}}}", tagKeys: {{{KEYS}}}) {
    errors {
      message
      type
    }
  }
}
`
