module Manifests (..) where

import Effects exposing (Effects, Never)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Json.Decode exposing ((:=))
import Json.Encode exposing (..)
import Task exposing (Task, andThen)
import Signal exposing (Address)
import StartApp
import Debug


app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = []
    }


main =
  app.html



-- MODEL


type alias Model =
  { manifests : List Manifest }


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


init : ( Model, Effects Action )
init =
  ( { manifests = [] }
  , getManifests
  )



-- UPDATE


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


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks


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



-- VIEW


firstCatalog : Maybe (List String) -> String
firstCatalog catalogs =
  case catalogs of
    Just catalogs ->
      catalogs
        |> List.head
        |> Maybe.withDefault ""

    Nothing ->
      ""


manifestRow address manifest =
  div
    [ class "manifestrow" ]
    [ li
        [ class "mitem" ]
        [ text manifest.name
        ]
    , li
        [ class "mitem" ]
        [ text (Maybe.withDefault "" manifest.displayName)
        ]
    , li
        [ class "mitem" ]
        [ text (firstCatalog manifest.catalogs)
        ]
    ]


manifestCollection address manifests =
  let
    manifestItems =
      List.map (manifestRow address) manifests
  in
    div
      [ id "manifests" ]
      [ div
          [ class "manifest_header_row" ]
          [ li
              [ class "manifest_header_item" ]
              [ h1 [] [ text "Manifest" ]
              , button
                  [ class "sort", onClick address (SortBy "name") ]
                  [ text "sort" ]
              ]
          , li
              [ class "manifest_header_item" ]
              [ h1 [] [ text "Display Name" ]
              , button
                  [ class "sort", onClick address (SortBy "name") ]
                  [ text "sort" ]
              ]
          , li
              [ class "manifest_header_item" ]
              [ h1 [] [ text "Catalogs" ]
              , button
                  [ class "sort", onClick address (SortBy "name") ]
                  [ text "sort" ]
              ]
          ]
      , div [] manifestItems
      ]



-- , ul [] manifestItems


view address model =
  div
    [ id "container" ]
    [ manifestCollection address model.manifests ]
