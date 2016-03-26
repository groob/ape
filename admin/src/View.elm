module View (..) where

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Models exposing (..)
import Update exposing (..)
import Routing


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


manifestPageView address model =
  div
    [ id "admin-page" ]
    [ div
        [ id "navbar" ]
        [ h1 [ onClick address Manifests ] [ text "Manifests" ]
        , h1 [] [ text "Pkgsinfos" ]
        ]
    , div
        []
        [ manifestView address model.manifests ]
    ]


pageView : Signal.Address Action -> Model -> Html
pageView address model =
  case model.routing.route of
    Routing.ManifestRoute ->
      div [] [ h2 [] [ text "About" ] ]

    Routing.NotFoundRoute ->
      div [] [ h2 [] [ text "Not found" ] ]


view address model =
  pageView address model


manifestView address manifests =
  div
    []
    [ div
        [ id "reset" ]
        [ button [ class "reset", onClick address Reset ] [ text "reset" ] ]
    , div
        [ id "container" ]
        [ manifestCollection address manifests ]
    ]
