
executable delegate.sh {

  plugin template {
    data-file = "data.yml"
  }

  plugin aws-user-data {}

}
