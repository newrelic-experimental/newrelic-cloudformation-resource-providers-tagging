package resource

import (
   "encoding/json"
   "fmt"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "strings"
)

//
// Generic, should be able to leave these as-is
//

type Payload struct {
   model  *Model
   models []interface{}
}

func NewPayload(m *Model) *Payload {
   m.Guid = m.EntityGuid
   return &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
}

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}

//
// These are API specific, must be configured per API
//

var typeName = "NewRelic::CloudFormation:Tagging"

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   // FIXME Fudge guid as it's the PK but passed in via EntityGuid
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Guid: &s})
}

var emptyFragment = ""

func (p *Payload) GetGraphQLFragment() *string {
   return &emptyFragment
}

func (p *Payload) SetGuid(g *string) {
   p.model.Guid = g
   log.Debugf("SetGuid: %s", *p.model.Guid)
}

func (p *Payload) GetGuid() *string {
   return p.model.Guid
}

func (p *Payload) GetCreateMutation() string {
   return `
mutation {
  taggingAddTagsToEntity(guid: "{{{GUID}}}" tags: {{{TAGS}}} ) {
    errors {
      message
      type
    }
  }
}
`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  taggingDeleteTagFromEntity(guid: "{{{GUID}}}", tagKeys: {{{KEYS}}}) {
    errors {
      message
      type
    }
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
  taggingReplaceTagsOnEntity(guid: "{{{GUID}}}", tags: {{{TAGS}}}) {
    errors {
      message
      type
    }
  }
}
`
}

func (p *Payload) GetReadQuery() string {
   return `
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
}

func (p *Payload) GetListQuery() string {
   return `
{
  actor {
    entitySearch(query: "tags is null") {
      results {
        entities {
          guid
          entityType
          tags {
            key
            values
          }
        }
        nextCursor
      }
    }
  }
}
`
}

func (p *Payload) GetListQueryNextCursor() string {
   return `
{
  actor {
    entitySearch(query: "tags is null") {
      results(cursor: "{{{NEXTCURSOR}}}") {
        entities {
          guid
          entityType
          tags {
            key
            values
          }
        }
        nextCursor
      }
    }
  }
}
`
}

func (p *Payload) GetGuidKey() string {
   // FIXME Only List returns a guid, this causes all other calls to fail
   return "guid"
}

func (p *Payload) GetResultKey(a model.Action) string {
   switch a {
   case model.List:
      return "guid"
   }
   return ""
}

func (p *Payload) GetVariables() map[string]string {
   // ACCOUNTID comes from the configuration
   // NEXTCURSOR is a _convention_

   if p.model.Variables == nil {
      p.model.Variables = make(map[string]string)
   }

   if p.model.EntityGuid != nil {
      p.model.Variables["GUID"] = *p.model.EntityGuid
   }

   if p.model.Tags != nil {
      // JSON Stringify the tags
      ba, err := json.Marshal(p.model.Tags)
      if err != nil {
         panic(err)
      }

      // Fix the case + GraphQL/JSON snafu
      tagString := string(ba)
      tagString = strings.ReplaceAll(tagString, `"Key"`, "key")
      tagString = strings.ReplaceAll(tagString, `"Values"`, "values")
      p.model.Variables["TAGS"] = tagString

      // Build an array of just the keys
      keys := make([]string, 0)
      for _, t := range p.model.Tags {
         keys = append(keys, *t.Key)
      }
      // JSON Stringify the key array
      ba, err = json.Marshal(keys)
      if err != nil {
         panic(err)
      }
      p.model.Variables["KEYS"] = string(ba)

   }

   lqf := ""
   if p.model.ListQueryFilter != nil {
      lqf = *p.model.ListQueryFilter
   }
   p.model.Variables["LISTQUERYFILTER"] = lqf

   return p.model.Variables
}

func (p *Payload) GetErrorKey() string {
   return "type"
}

func (p *Payload) NeedsPropagationDelay(a model.Action) bool {
   // Tags work on existing entities, no delay required
   return false
}

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
   Tags       []tag  `json:"tags,omitempty"`
   Type       string `json:"type"`
}
type tag struct {
   Key    string   `json:"key"`
   Values []string `json:"values"`
}

func (p *Payload) TestReadResponse(data []byte) (err error) {
   r := readResponse{}
   if err = json.Unmarshal(data, &r); err != nil {
      return
   }
   if p.model.Tags == nil && r.Data.Actor.Entity.Tags == nil {
      return
   }
   if p.model.Tags != nil && r.Data.Actor.Entity.Tags == nil {
      err = fmt.Errorf("%w model Tags nil and read tags not nil", &cferror.NotFound{})
      return
   }
   if p.model.Tags == nil && r.Data.Actor.Entity.Tags != nil {
      err = fmt.Errorf("%w model Tags not nil and read tags nil", &cferror.NotFound{})
      return
   }
   tagsEqual := true
   // Compare the model to the read result. Everything in the model must be present in the read- the read may be bigger
   for _, modelTag := range p.model.Tags {
      // If the model tag key is in the read tag array
      if readValue, ok := containsTagKey(*modelTag.Key, r.Data.Actor.Entity.Tags); ok {
         // Test the model key's values and ensure each is in the read tag's value array
         for _, modelValue := range modelTag.Values {
            if containsTagValue(modelValue, readValue) {
               continue
            } else {
               tagsEqual = false
               break
            }
         }
      } else {
         tagsEqual = false
      }
   }
   if tagsEqual {
      return
   }
   err = fmt.Errorf("%w model Tags not nil and read tags nil", &cferror.NotFound{})
   return
}

func containsTagValue(rv string, mvs []string) bool {
   for _, mv := range mvs {
      if rv == mv {
         return true
      }
   }
   return false
}

func containsTagKey(key string, tags []tag) ([]string, bool) {
   for _, t := range tags {
      if t.Key == key {
         return t.Values, true
      }
   }
   return nil, false
}
