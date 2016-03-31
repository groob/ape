module Manifests.Actions (..) where

import Http
import Manifests.Models exposing (..)


type Action
  = GetManifests (Result Http.Error (List Manifest))
  | SortBy String
