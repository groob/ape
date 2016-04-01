module Models (..) where

import Routing
import Manifests.Models exposing (Manifest, manifest)


type alias Model =
  { manifests : List Manifest
  , manifestForm : Maybe Manifest
  , pkgsinfos : List Pkgsinfo
  , routing : Routing.Model
  }


type alias Pkgsinfo =
  { name : String }
