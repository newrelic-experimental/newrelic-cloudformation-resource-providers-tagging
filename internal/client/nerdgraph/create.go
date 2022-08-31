package nerdgraph

import (
   "encoding/json"
   "fmt"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/model"
   "newrelic-cloudformation-tagging/internal/utils"
   "strings"
)

/*
{mutation {
  taggingAddTagsToEntity(guid: "", tags: [{key: "key1", values: ["key1-value1"]}]) {
    errors {
      message
      type
    }
  }
}
*/
type createResponse struct {
   Data createData `json:"data"`
}
type createData struct {
   TaggingAddTagsToEntity taggingAddTagsToEntity `json:"taggingAddTagsToEntity"`
}
type taggingAddTagsToEntity struct {
   Errors []errors `json:"errors"`
}

func (i *nerdgraph) Create(m *model.Model) error {
   log.Debugf("nerdgraph/client.Create model: %+v", m)
   // TODO abstract the rendering out of the "framework"

   // Deal with semantics
   if m.IsAWSSemantics() {
      // Map semantics: if a key exists remove it so there are no duplicate keys
      for _, t := range m.Tags {
         err := i.deleteTag(m, `"`+*t.Key+`"`)
         // TODO will we get an error if the key is not found? That needs a test case.
         if err != nil {
            return err
         }
      }
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
   log.Debugln("Create: rendered tags: ", tagString)

   // Render the mutation
   // TODO abstract the substitution map
   mutation, err := model.Render(createMutation, map[string]string{"GUID": *m.Guid, "TAGS": tagString})
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }
   log.Debugln("Create: rendered mutation: ", mutation)
   log.Debugln("")

   // Validate mutation
   err = m.Validate(&mutation)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *m.APIKey, m.GetEndpoint())
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   response := createResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      log.Errorf("Create: %v", err)
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }
   utils.DumpModel(response, "Create: response: ")
   if response.Data.TaggingAddTagsToEntity.Errors != nil {
      for _, e := range response.Data.TaggingAddTagsToEntity.Errors {
         log.Errorf("Create: NerdGraph error: %s %s", e.Type, e.Message)
         err = fmt.Errorf("%w %s %s", &model.InvalidRequest{}, e.Type, e.Message)
      }
      return err
   }
   return i.Read(m)
}

const createMutation = `
mutation {
  taggingAddTagsToEntity(guid: "{{{GUID}}}" tags: {{{TAGS}}} ) {
    errors {
      message
      type
    }
  }
}
`
