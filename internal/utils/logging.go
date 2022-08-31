package utils

import (
   "encoding/json"
   log "github.com/sirupsen/logrus"
)

func DumpModel(m interface{}, s string) {
   b, err := json.MarshalIndent(m, "", "   ")
   if err == nil {
      log.Debugln(s, " ", string(b))
   } else {
      log.Debugf("Marshal error: %s %v", s, err)
   }
}
