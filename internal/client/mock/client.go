package mock

import (
   "fmt"
   log "github.com/sirupsen/logrus"
   "newrelic-cloudformation-tagging/internal/model"
   "time"
)

type mock struct {
   guids map[string]interface{}
}

// TODO Config
const progressDelaySeconds = 0

func NewClient() *mock {
   return &mock{guids: map[string]interface{}{}}
}

func (i *mock) Create(m *model.Model) (err error) {
   log.Debugf("mock/model.Create guids: %+v", i.guids)
   time.Sleep(progressDelaySeconds * time.Second)
   if m.Guid != nil {
      if _, v := i.guids[*m.Guid]; v {
         return fmt.Errorf("%w %s", &model.AlreadyExists{}, "Already exists")
      }
   }
   i.guids[*m.Guid] = nil
   err = nil
   return
}

func (i *mock) Delete(m *model.Model) (err error) {
   log.Debugf("mock/model.Delete guids: %+v", i.guids)
   time.Sleep(progressDelaySeconds * time.Second)
   if m.Guid != nil {
      if _, v := i.guids[*m.Guid]; v {
         delete(i.guids, *m.Guid)
         return nil
      } else {
         return fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
      }
   } else {
      return fmt.Errorf("%w %s", &model.NotFound{}, "missing guid")
   }
}

func (i *mock) Update(m *model.Model) (err error) {
   log.Debugf("mock/model.Update guids: %+v", i.guids)
   time.Sleep(progressDelaySeconds * time.Second)
   if m.Guid != nil {
      if _, v := i.guids[*m.Guid]; v {
         return nil
      } else {
         return fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
      }
   } else {
      return fmt.Errorf("%w %s", &model.NotFound{}, "missing guid")
   }
}

func (i *mock) Read(m *model.Model) (err error) {
   log.Debugf("mock/model.Read guids: %+v", i.guids)
   time.Sleep(progressDelaySeconds * time.Second)
   if m.Guid != nil {
      if _, v := i.guids[*m.Guid]; v {
         return nil
      } else {
         return fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
      }
   } else {
      return fmt.Errorf("%w %s", &model.NotFound{}, "missing guid")
   }
}

func (i *mock) List(m *model.Model) (r []interface{}, err error) {
   log.Debugf("mock/model.List guids: %+v", i.guids)
   time.Sleep(progressDelaySeconds * time.Second)
   if m.Guid != nil {
      if _, v := i.guids[*m.Guid]; v {
         return []interface{}{m}, nil
      } else {
         return nil, fmt.Errorf("%w %s", &model.NotFound{}, "guid not found")
      }
   } else {
      return nil, fmt.Errorf("%w %s", &model.NotFound{}, "missing guid")
   }
}
