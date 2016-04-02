module Manifests.Actions (..) where

import Http
import Manifests.Models exposing (..)


type Action
  = GetManifests (Result Http.Error (List Manifest))
  | GetManifest (Maybe Manifest)
  | EditManifest String
  | DiscardSave
  | UpdateDisplayName String
  | Save Manifest
  | SortBy String
  | HopAction ()
  | NoOp
  | DeleteManifest String
