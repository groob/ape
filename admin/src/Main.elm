module Main (..) where

import Effects exposing (Effects, Never)
import Task exposing (Task)
import Signal exposing (Address)
import StartApp
import Routing
import Manifests.Update
import Manifests.Models exposing (Manifest)
import Models exposing (..)
import View exposing (..)
import Update exposing (..)
import Client.Http exposing (getManifests)


init : ( Model, Effects Action )
init =
  let
    fxs =
      [ Effects.map ManifestAction Client.Http.getManifests ]

    fx =
      Effects.batch fxs
  in
    ( { manifests = []
      , pkgsinfos = []
      , manifestForm = Nothing
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


port routeRunTask : Task () ()
port routeRunTask =
  Routing.run


main =
  app.html
