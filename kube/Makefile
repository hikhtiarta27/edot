config-app:
	kubectl create configmap auth-config --from-file=../service/auth/config.yaml
	kubectl create configmap auth-config --from-file=../service/auth/config.yaml
	kubectl create configmap product-config --from-file=../service/product/config.yaml
	kubectl create configmap warehouse-config --from-file=../service/warehouse/config.yaml
	kubectl create configmap order-config --from-file=../service/order/config.yaml
	kubectl create configmap shop-config --from-file=../service/shop/config.yaml

observability:
	kubectl apply -f jaeger.yaml
	kubectl apply -f otel.yaml

storage:
	kubectl apply -f mysql.yaml
	kubectl apply -f redis.yaml

config:
	kubectl create configmap config-default --from-file=./config -n edot

run: storage observability
