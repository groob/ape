module View (..) where

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Models exposing (..)
import Update exposing (..)
import Routing
import Manifests.View exposing (manifestView)


pageView : Signal.Address Action -> Model -> Html
pageView address model =
  case model.routing.route of
    Routing.AdminHome ->
      div [] [ h2 [] [ text "Manage Munki Home" ] ]

    Routing.ManifestCollectionRoute ->
      manifestView (Signal.forwardTo address ManifestAction) model.manifests

    Routing.PkgsInfoCollectionRoute ->
      div [] [ h2 [] [ text "PkgsInfo list" ] ]

    Routing.NotFoundRoute ->
      div [] [ h2 [] [ text "Not found" ] ]


view address model =
  pageView address model
