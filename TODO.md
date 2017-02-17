# generate ca cert
# sign sever certs with the ca cert
# have an endpoint which serves the server cert
# client calls the endpoint to retrieve the server cert
# client generates its own key pair
# client sends the public key to the server for signing
# server maintains a copy of clients public key
# have an endpoint on the server which serves the
