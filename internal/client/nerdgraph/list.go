package nerdgraph

import (
   "encoding/json"
   "fmt"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/model"
)

type listResponse struct {
   Data listData `json:"data"`
}
type listData struct {
   Actor listActor `json:"actor"`
}
type listActor struct {
   EntitySearch listEntitySearch `json:"entitySearch"`
}
type listEntitySearch struct {
   Count   int         `json:"count"`
   Results listResults `json:"results"`
}
type listResults struct {
   Entities   []listEntities `json:"entities"`
   NextCursor string         `json:"nextCursor"`
}
type listEntities struct {
   Guid string `json:"guid"`
}

// List only gets 30 seconds to do its work, IN_PROGRESS is not allowed
// NOTE: entitySearch requires several seconds to index a newly created entity. Read the guid in the model and append it to the list result.
func (i *nerdgraph) List(m *model.Model) ([]interface{}, error) {
   // TODO what is the List search criteria? How do we query for all entities in an account that have tags?
   log.Debugf("List: enter: guid: %s", *m.Guid)
   result := make([]interface{}, 0)

   // Because of the indexing delay on the guid after create do an entity query
   err := i.Read(m)
   if err != nil {
      return result, err
   }
   result = append(result, m)

   filter := ""
   if m.ListQueryFilter != nil {
      filter = *m.ListQueryFilter
   }
   mutation, err := model.Render(listQuery, map[string]string{"LISTQUERYFILTER": filter})
   if err != nil {
      return nil, fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   body, err := i.emit(mutation, *m.APIKey, m.GetEndpoint())
   if err != nil {
      return nil, fmt.Errorf("%w %s", &model.InvalidRequest{}, err.Error())
   }

   response := listResponse{}
   err = json.Unmarshal(body, &response)
   if err != nil {
      return result, err
   }

   // NOTE: entitySearch does not return errors
   // TODO to return the tags. Case on key & value might be an issue
   for _, e := range response.Data.Actor.EntitySearch.Results.Entities {
      result = append(result, &model.Model{
         EntityGuid: &e.Guid,
         APIKey:     m.APIKey,
         Tags:       nil,
      })
   }
   // TODO process cursor
   return result, nil
}

const listQuery = `
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
const listQueryNextCursor = `
{
  actor {
    entitySearch(query: "tags is null") {
      results(cursor: "") {
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
