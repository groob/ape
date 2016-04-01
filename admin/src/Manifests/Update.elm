module Manifests.Update (..) where

import Manifests.Actions exposing (..)
import Manifests.Models exposing (..)
import Effects exposing (Effects)
import Http
import Hop
import Hop.Navigate exposing (navigateTo)
import Task exposing (Task, andThen)
import Json.Decode exposing ((:=))
import Html exposing (Attribute)
import Html.Events exposing (on, targetValue)
import Signal exposing (Address)


update action model =
  case action of
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

    EditManifest id ->
      let
        path =
          "/manifests/" ++ id ++ "/edit"

        filterByID : String -> List Manifest -> Maybe Manifest
        filterByID id manifests =
          manifests
            |> List.filter (\manifest -> manifest.name == id)
            |> List.head

        manifest : Maybe Manifest
        manifest =
          filterByID id model.manifests
      in
        ( { model | manifestForm = (Debug.log "manifestX" (filterByID id model.manifests)) }, Effects.map HopAction (navigateTo path) )

    DiscardSave ->
      ( model, Effects.map HopAction (navigateTo "/manifests") )

    UpdateDisplayName contents ->
      let
        manifest =
          model.manifestForm

        updateField contents manifest =
          { manifest | displayName = Just contents }

        updatedManifestForm =
          Maybe.map (updateField contents) manifest
      in
        ( { model | manifestForm = updatedManifestForm }, Effects.none )

    Save manifest ->
      let
        updateManifest existing =
          if existing.name == manifest.name then
            manifest
          else
            existing

        updatedCollection =
          List.map updateManifest model.manifests
      in
        ( { model | manifests = updatedCollection }, Effects.map HopAction (navigateTo "/manifests") )

    HopAction _ ->
      ( model, Effects.none )

    NoOp ->
      ( (Debug.log "model" model), Effects.none )


getManifests =
  Http.get (Json.Decode.list manifest) "/api/manifests"
    |> Task.toResult
    |> Task.map GetManifests
    |> Effects.task


onInput : Address a -> (String -> a) -> Attribute
onInput address f =
  on "input" targetValue (\v -> Signal.message address (f v))


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
