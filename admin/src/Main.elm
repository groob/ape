module Main (..) where

import Effects exposing (Effects, Never)
import Http
import Task exposing (Task)
import Signal exposing (Address)
import StartApp
import Routing
import Debug
import Manifests.Update
import Manifests.Models


-- app imports

import Models exposing (..)
import View exposing (..)
import Update exposing (..)


init : ( Model, Effects Action )
init =
  let
    fxs =
      [ Effects.map ManifestAction Manifests.Update.getManifests ]

    fx =
      Effects.batch fxs
  in
    ( { manifests = []
      , pkgsinfos = []
      , routing = Routing.initialModel
      }
    , fx
    )


app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = [ routerSignal ]
    }


routerSignal : Signal Action
routerSignal =
  Signal.map RoutingAction Routing.signal


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks


main =
  app.html
