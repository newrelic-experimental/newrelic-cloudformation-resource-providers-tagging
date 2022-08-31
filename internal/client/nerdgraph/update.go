package nerdgraph

import (
   "encoding/json"
   "fmt"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/model"
   "strings"
)

type updateResponse struct {
   Data updateData `json:"data"`
}
type updateData struct {
   TaggingReplaceTagsOnEntity taggingReplaceTagsOnEntity `json:"taggingReplaceTagsOnEntity"`
}
type taggingReplaceTagsOnEntity struct {
   Errors []errors `json:"errors"`
}

func (i *nerdgraph) Update(m *model.Model) error {
   log.Debugf("nerdgraph/model.Update model: %+v", m)
   if m.IsAWSSemantics() {
      // Read will return NotFound if the keys are not present
      // TODO
      err := i.Read(m)
      if err != nil {
         return err
      }
      // As the key is present and we're in AWS mode, Create gets the job done
      return i.Create(m)
   }

   // JSON Stringify the tags
   ba, err := json.Marshal(m.Tags)
   if err != nil {
      return err
   }

   // Render the tags
   tagString, err := model.Render(string(ba), m.Variables)
   tagString = strings.ReplaceAll(tagString, `"Key"`, "key")
   tagString = strings.ReplaceAll(tagString, `"Values"`, "values")

   // Render the mutation
   // TODO abstract the substitution map
   mutation, err := model.Render(updateMutation, map[string]string{"GUID": *m.Guid, "TAGS": tagString})
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }
   log.Debugf("Update- rendered mutation: %s", mutation)

   // Validate mutation
   err = m.Validate(&mutation)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   log.Debugf("Update: mutation: %s\n model: %+v", mutation, m)
   body, err := i.emit(mutation, *m.APIKey, m.GetEndpoint())
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   response := updateResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Update: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   if response.Data.TaggingReplaceTagsOnEntity.Errors != nil {
      for _, e := range response.Data.TaggingReplaceTagsOnEntity.Errors {
         log.Errorf("Update: NerdGraph error: %s %s", e.Type, e.Message)
         err = fmt.Errorf("%w %s %s", &model.InvalidRequest{}, e.Type, e.Message)
      }
      return err
   }
   return nil
}

const updateMutation = `
mutation {
  taggingReplaceTagsOnEntity(guid: "{{{GUID}}}", tags: {{{TAGS}}}) {
    errors {
      message
      type
    }
  }
}
`
