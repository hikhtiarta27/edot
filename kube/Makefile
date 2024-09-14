observability:
	kubectl apply -f tools-namespace.yaml
	kubectl create configmap config-default --from-file=./config -n tools
	kubectl apply -f jaeger.yaml
	kubectl apply -f otel.yaml

cron:
	kubectl apply -f cron.yaml

run: observability

deploy-auth:
	kubectl apply -f auth/namespace.yaml
	kubectl create configmap auth-config --from-file=auth/config.yaml -n auth
	kubectl apply -f auth/storage.yaml
	kubectl apply -f auth/deployment.yaml

deploy-product:
	kubectl apply -f product/namespace.yaml
	kubectl create configmap product-config --from-file=product/config.yaml -n product
	kubectl apply -f product/storage.yaml
	kubectl apply -f product/deployment.yaml

deploy-shop:
	kubectl apply -f shop/namespace.yaml
	kubectl create configmap shop-config --from-file=shop/config.yaml -n shop
	kubectl apply -f shop/storage.yaml
	kubectl apply -f shop/deployment.yaml

deploy-warehouse:
	kubectl apply -f warehouse/namespace.yaml
	kubectl create configmap warehouse-config --from-file=warehouse/config.yaml -n warehouse
	kubectl apply -f warehouse/storage.yaml
	kubectl apply -f warehouse/deployment.yaml