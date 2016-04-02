module Actions (..) where


type Action
  = NoOp
  | RoutingAction Routing.Action
  | ManifestAction Manifests.Actions.Action
