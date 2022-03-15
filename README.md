# iamuser-operator-demo

## Steps
### 1. Install kubebuilder
 
```

# download kubebuilder and install locally.
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

```
### 2. Create a project
 
```

mkdir -p ~/goperator/iamuser
cd ~/goperator/iamuser
kubebuilder init --domain govind.kv --repo govind.kv/iamuser
```
### 3. Create an API
 ```

kubebuilder create api --version v1alpha1 --kind IamUser

:info:  Press 'y' for both resource and controller prompts.


Create Resource [y/n]
y
Create Controller [y/n]
y
```
### 4. Modify the spec and status structs as required in api/v1alpha1/types.go

### 5. Add the reconcilation logic in the Reconcile method inside controllers/iamuser_controller.go

### 6. Generate CRD
```
make generate
```
### 7. Install CRD
```
make install
```
### 8. Apply the CR
```
kubectl apply -f iamuser.yaml
```
