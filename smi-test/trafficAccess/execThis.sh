
echo $(kubectl get service service-a --namespace=$NAMESPACE -o jsonpath="{.spec.ports[0].nodePort}")
kubectl describe pods --namespace=$NAMESPACE
curl --location --request POST 'http://localhost:'$(echo $(kubectl get service service-a --namespace=$NAMESPACE -o jsonpath="{.spec.ports[0].nodePort}"))'/call' -w "%{http_code}" --data-raw '{"host": "http://service-b/post",}'


curl --location --request GET 'localhost:'$(echo $(kubectl get service service-a --namespace=$NAMESPACE -o jsonpath="{.spec.ports[0].nodePort}"))'/metrics' --header 'Content-Type: application/json' --data-raw '{"hello": "bye"}'