module Manifests.View (..) where

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Manifests.Actions exposing (..)
import Routing
import Manifests.Models exposing (Manifest)
import Manifests.Utils exposing (onInput)


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
    , button [ onClick address (EditManifest manifest.name) ] [ text "edit" ]
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


manifestView address manifests =
  div
    []
    [ div
        [ id "container" ]
        [ manifestCollection address manifests ]
    ]


defaultManifest : Manifest
defaultManifest =
  { name = "foo"
  , catalogs = Nothing
  , displayName = Nothing
  }


manifestEdit : Signal.Address Action -> Maybe Manifest -> Html
manifestEdit address manifest =
  case manifest of
    Nothing ->
      div [] [ text "not found" ]

    Just manifest ->
      div
        []
        [ input
            [ type' "text"
            , placeholder "Name"
            , value manifest.name
            , name "name"
            , autofocus True
            ]
            []
        , input
            [ type' "text"
            , placeholder "display name"
            , value (Maybe.withDefault "" manifest.displayName)
            , name "DisplayName"
            , autofocus False
            , onInput address UpdateDisplayName
            ]
            []
        , button [ class "save", onClick address (Save manifest) ] [ text "save" ]
        , button [ class "discard", onClick address DiscardSave ] [ text "discard" ]
        ]
