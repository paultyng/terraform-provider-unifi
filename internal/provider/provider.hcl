provider "unifi" {
    sdk_package = "github.com/paultyng/go-unifi/unifi"

    client_type = "Client" # default

    # go templates that have resource info object
    read_func   = "Get{{ .TypeName }}" # default
    create_func = "Create{{ .TypeName }}" # default
    update_func = "Update{{ .TypeName }}" # default
    delete_func = "Delete{{ .TypeName }}" # default
}