

login:
	curl -XPOST -H "Content-Type: application/json" -d '{"name": "cyrushiker", "password": "111111"}' http://localhost:9090/user/login

clear:
	curl -XPOST http://localhost:9090/user/clear