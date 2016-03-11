module Update (..) where

import Http
import Task exposing (Task, andThen)
import Models exposing (..)
import Effects exposing (Effects)
import Json.Decode exposing ((:=))
import Json.Encode exposing (..)


type Action
  = NoOp
  | GetManifests (Result Http.Error (List Manifest))
  | SortBy String


update action model =
  case action of
    NoOp ->
      ( model, Effects.none )

    GetManifests result ->
      case result of
        Ok manifests ->
          ( { model | manifests = manifests }, Effects.none )

        Err error ->
          let
            _ =
              reportError error
          in
            ( model, Effects.none )

    SortBy filter ->
      ( { model | manifests = (List.reverse model.manifests) }, Effects.none )


getManifests =
  Http.get (Json.Decode.list manifest) "/api/manifests"
    |> Task.toResult
    |> Task.map GetManifests
    |> Effects.task


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
