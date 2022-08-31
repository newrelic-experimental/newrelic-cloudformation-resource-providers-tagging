package nerdgraph

import (
   "encoding/json"
   "fmt"
   "newrelic-cloudformation-tagging/internal/model"
   "newrelic-cloudformation-tagging/internal/utils"
)

/*
{
*/
type readResponse struct {
   Data readData `json:"data"`
}
type readData struct {
   Actor readActor `json:"actor"`
}
type readActor struct {
   Entity *readEntity `json:"entity,omitempty"`
}
type readEntity struct {
   Domain     string `json:"domain"`
   EntityType string `json:"entityType"`
   Guid       string `json:"guid"`
   Name       string `json:"name"`
   Tags       []Tag  `json:"tags,omitempty"`
   Type       string `json:"type"`
}

type Tag struct {
   Key    string   `json:"key"`
   Values []string `json:"values"`
}

func (i *nerdgraph) Read(m *model.Model) error {
   if m.Guid == nil {
      return fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
   }

   mutation, err := model.Render(readQuery, map[string]string{"GUID": *m.Guid})
   if err != nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *m.APIKey, m.GetEndpoint())
   if err != nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   response := readResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      return fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   // An entity query returns no errors so nothing to check
   if response.Data.Actor.Entity == nil {
      return fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
   }

   if m.IsAWSSemantics() {
      if response.Data.Actor.Entity.Tags == nil {
         // No tags in so no tags out is fine
         if m.Tags == nil {
            return nil
         }
         return fmt.Errorf("%w %s", &model.NotFound{}, "tags not found")
      }
      utils.DumpModel(response.Data.Actor.Entity.Tags, "Read: response: ")
      // Build a map of all of the tag keys in the response
      outMap := make(map[string]string, len(response.Data.Actor.Entity.Tags))
      for _, outTag := range response.Data.Actor.Entity.Tags {
         outMap[outTag.Key] = ""
      }

      b, s := m.AllTagsPresent(outMap)
      if !b {
         return fmt.Errorf("%w missing tags: %s", &model.NotFound{}, s)
      }
   }
   return nil
}

const readQuery = `
{
  actor {
    entity(guid: "{{{GUID}}}") {
      tags {
         key
         values
      }
      domain
      entityType
      guid
      name
      type
    }
  }
}
`
