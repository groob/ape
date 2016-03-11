module Main (..) where

import Effects exposing (Effects, Never)
import Http
import Task exposing (Task, andThen)
import Signal exposing (Address)
import StartApp
import Debug


-- app imports

import Models exposing (..)
import View exposing (..)
import Update exposing (..)


init : ( Model, Effects Action )
init =
  ( { manifests = []
    , pkgsinfos = []
    }
  , getManifests
  )


app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = []
    }


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks


main =
  app.html
