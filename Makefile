

login:
	curl -XPOST -H "Content-Type: application/json" -d '{"name": "me", "password": "me"}' http://localhost:9090/user/login

clear:
	curl -XPOST http://localhost:9090/user/clear