module Update (..) where

import Http
import Task exposing (Task, andThen)
import Models exposing (..)
import Effects exposing (Effects)
import Json.Decode exposing ((:=))
import Json.Encode exposing (..)
import Routing
import Manifests.Actions
import Manifests.Update
import Manifests.Models


type Action
  = NoOp
  | RoutingAction Routing.Action
  | ManifestAction Manifests.Actions.Action


update action model =
  case (Debug.log "action" action) of
    RoutingAction subAction ->
      let
        ( updatedRouting, fx ) =
          Routing.update subAction model.routing
      in
        ( { model | routing = updatedRouting }, Effects.map RoutingAction fx )

    NoOp ->
      ( model, Effects.none )

    ManifestAction subAction ->
      let
        --   updatedModel =
        --     { manifests = model.manifests }
        --
        ( updatedModel, fx ) =
          Manifests.Update.update subAction model
      in
        ( updatedModel, Effects.map ManifestAction fx )


reportError : Http.Error -> Http.Error
reportError error =
  case error of
    Http.Timeout ->
      Debug.log "API timeout" error

    Http.NetworkError ->
      Debug.log "Network error contacting API" error

    Http.UnexpectedPayload payload ->
      Debug.log ("Unexpected payload from API: " ++ payload) error

    Http.BadResponse status payload ->
      Debug.log ("Unexpected status/payload from API: " ++ (toString status) ++ "/" ++ payload) error
