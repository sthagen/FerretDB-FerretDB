# see https://github.com/FiloSottile/mkcert

export CAROOT=$(PWD)
export TRUST_STORES=system
HOSTS=localhost 127.0.0.1 ::1 host.docker.internal

gen:
	rm -f *.pem
	mkcert -install
	mkcert -client -cert-file client-cert.pem -key-file client-key.pem $(HOSTS)
	cat client-cert.pem client-key.pem > client.pem
	mkcert -cert-file server-cert.pem -key-file server-key.pem $(HOSTS)
	cat server-cert.pem server-key.pem > server.pem
	mkcert -uninstall
	mv rootCA.pem rootCA-cert.pem
	rm -f rootCA-key.pem
