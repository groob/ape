module Manifests.SingleView (..) where

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Manifests.Actions exposing (..)
import Manifests.Models exposing (Manifest)
import Manifests.Utils exposing (onInput)


-- EDIT/VIEW single Manifest


manifestEdit : Signal.Address Action -> Maybe Manifest -> Html
manifestEdit address manifest =
  case (Debug.log "manifest" manifest) of
    Nothing ->
      div [ onClick address NoOp ] [ text "not found" ]

    Just manifest ->
      editPage address manifest


catalogRow : Signal.Address Action -> String -> Html
catalogRow address catalog =
  div
    [ class "catalogrow" ]
    [ input
        [ type' "checkbox"
        , checked True
        ]
        []
    , text catalog
    ]


catalogView : Signal.Address Action -> Manifest -> Html
catalogView address manifest =
  let
    catalogItems =
      List.map (catalogRow address) (Maybe.withDefault [] manifest.catalogs)
  in
    div [] catalogItems


editPage address manifest =
  div
    [ class "manifest-menu" ]
    [ h1 [] [ text "General" ]
    , h1 [] [ text "Included Manifests" ]
    , h1 [] [ text "Referencing Manifests" ]
    , formScreen address manifest
    ]


formScreen address manifest =
  div
    [ class "form-box" ]
    [ manifestView address manifest
    ]


manifestView address manifest =
  div
    []
    [ text "Name:"
    , input
        [ type' "text"
        , placeholder "Name"
        , value manifest.name
        , name "name"
        , autofocus True
        ]
        []
    , text "Display name"
    , input
        [ type' "text"
        , placeholder "display name"
        , value (Maybe.withDefault "" manifest.displayName)
        , name "DisplayName"
        , autofocus False
        , onInput address UpdateDisplayName
        ]
        []
    , text "User:"
    , input
        [ type' "text"
        , value (Maybe.withDefault "" manifest.user)
        , name "user"
        ]
        []
    , text "Notes:"
    , input
        [ type' "text"
        , value (Maybe.withDefault "" manifest.notes)
        , name "notes"
        ]
        []
    , text "Catalogs:"
    , catalogView address manifest
    , button [ class "save", onClick address (Save manifest) ] [ text "save" ]
    , button [ class "discard", onClick address DiscardSave ] [ text "discard" ]
    ]
