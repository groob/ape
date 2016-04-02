module Client.Http (..) where

import Http
import Effects exposing (Effects)
import Task exposing (Task, andThen)
import Manifests.Models exposing (..)
import Manifests.Actions exposing (..)
import Json.Decode exposing ((:=))


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


getManifests =
  Http.get (Json.Decode.list manifest) "/api/manifests"
    |> Task.toResult
    |> Task.map GetManifests
    |> Effects.task


getManifest name =
  let
    url =
      "/api/manifests/" ++ name
  in
    Http.get manifest url
      |> Task.toMaybe
      |> Task.map GetManifest
      |> Effects.task
