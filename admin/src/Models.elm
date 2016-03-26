module Models (..) where

import Json.Decode exposing ((:=))
import Json.Encode exposing (..)
import Effects exposing (Effects, Never)
import Routing


type alias Model =
  { manifests : List Manifest
  , pkgsinfos : List Pkgsinfo
  , routing : Routing.Model
  }


type alias Pkgsinfo =
  { name : String }


type alias Manifest =
  { name : String
  , catalogs : Maybe (List String)
  , displayName : Maybe String
  }


manifest : Json.Decode.Decoder Manifest
manifest =
  Json.Decode.object3
    Manifest
    ("filename" := Json.Decode.string)
    (Json.Decode.maybe ("catalogs" := Json.Decode.list Json.Decode.string))
    (Json.Decode.maybe ("display_name" := Json.Decode.string))
