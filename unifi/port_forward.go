package unifi

// curl 'https://73.212.25.176:8443/api/s/default/rest/portforward' --data-binary '{"name":"test name","enabled":true,"src":"200.200.200.200/24","dst_port":"90","fwd":"10.1.5.7","fwd_port":"90","proto":"tcp","log":true}' --compressed --insecure
// {"meta":{"rc":"ok"},"data":[{"name":"test name","enabled":true,"src":"200.200.200.200/24","dst_port":"90","fwd":"10.1.5.7","fwd_port":"90","proto":"tcp","log":true,"site_id":"5d6d8b07439adf048407dcd9","_id":"5e0f4ffa1e801c052a25d953"}]}
