/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	govinddevv1alpha1 "govind.dev/iamuser/api/v1alpha1"
)

// IamUserReconciler reconciles a IamUser object
type IamUserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=govind.dev,resources=iamusers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=govind.dev,resources=iamusers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=govind.dev,resources=iamusers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IamUser object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *IamUserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	fmt.Println("We are inside the reconciler")

	var user govinddevv1alpha1.IamUser
	err := r.Get(ctx, req.NamespacedName, &user)
	if err != nil {
		logErr("Couldnt fetch object", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//check for user deletion
	myfinalizer := "iamusers.govind.dev/finalizer"
	err = r.checkDeletion(myfinalizer,&user,ctx,l)
	if err!= nil{
		return ctrl.Result{}, err
	}
	//check for user creation
	if !user.Status.Usercreated{
		err = r.CreateUserReconcile(ctx, &user, l)
		if err != nil {
			fmt.Println(err.Error())
			return ctrl.Result{}, err
		}
	}

	//check for user updation
	if user.Status.Usercreated && user.Spec.Username != user.Status.Username {
		err = r.UpdateUserReconcile(ctx, &user, l)
		if err != nil {
			fmt.Println(err.Error())
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IamUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&govinddevv1alpha1.IamUser{}).
		Complete(r)
}

func logErr(msg string, err error) {
	fmt.Println(msg, "::", err.Error())
}


func (r *IamUserReconciler) CreateUserReconcile(ctx context.Context, user *govinddevv1alpha1.IamUser, l logr.Logger) error {
	l.Info("It's a CREATE request")
	svc, err := AwsIamSession(l)
	if err != nil {
		return err
	}
	existingUser, err := svc.GetUser(&iam.GetUserInput{
		UserName: &user.Spec.Username,
	})
	
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				err = CreateIamUser(svc, user, l)
				if err != nil {
					user.Status.Usercreated = false
					user.Status.Username = user.Spec.Username
					r.Status().Update(ctx, user)
					return err
				}
				user.Status.Usercreated = true
				user.Status.Username = user.Spec.Username
				r.Status().Update(ctx, user)
				l.Info("USER CREATED")
				return nil
			}
		}
		return err
	}
	if existingUser != nil {
		l.Error(err, "User with the same name already exists")
		return nil
	}
	return nil
}



func (r *IamUserReconciler) UpdateUserReconcile(ctx context.Context, user *govinddevv1alpha1.IamUser, l logr.Logger) error {
	l.Info("It's an UPDATE request")
	svc, err := AwsIamSession(l)
	if err != nil {
		return err
	}
	existingUser, err := svc.GetUser(&iam.GetUserInput{
		UserName: &user.Spec.Username,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				err = UpdateIamUser(svc, user, l)
				if err != nil {
					user.Status.Usercreated = false
					user.Status.Username = user.Spec.Username
					r.Status().Update(ctx, user)
					return err
				}
				user.Status.Usercreated = true
				user.Status.Username = user.Spec.Username
				r.Status().Update(ctx, user)
				l.Info("USER UPDATED")
				return nil
			}
		}
		return err
	}
	if existingUser != nil {
		l.Error(err, "User with the same name already exists")
		return nil
	}
	return nil
}

func (r *IamUserReconciler) DeleteUserReconcile(ctx context.Context, user *govinddevv1alpha1.IamUser,l logr.Logger) error {
	l.Info("It's a DELETE request")
	svc, err := AwsIamSession(l)
	if err != nil {
		return err
	}
	fmt.Println("User to be deleted : ", user.Spec.Username)
	err = DeleteIamUser(svc, user, l)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			if aerr.Code() == iam.ErrCodeNoSuchEntityException {
				l.Error(err, "No user with the name exists")
				user.Status.Usercreated = false
				return nil
			}
		}
		return err
	}
	l.Info("USER DELETED")
	return nil
}

func (r *IamUserReconciler) PrintList(ctx context.Context, l logr.Logger) {
	userlist := govinddevv1alpha1.IamUserList{}
	r.List(ctx, &userlist)
	for _, u := range userlist.Items {
		fmt.Println(u.Name, u.GetCreationTimestamp())
	}
}

func CreateIamUser(svc *iam.IAM, user *govinddevv1alpha1.IamUser, l logr.Logger) error {
	userCreate, err := svc.CreateUser(&iam.CreateUserInput{
		UserName: &user.Spec.Username,
	})
	if err != nil {
		return err
	}
	user.Status.UserArn = *userCreate.User.Arn
	return nil
}

func DeleteIamUser(svc *iam.IAM, user *govinddevv1alpha1.IamUser, l logr.Logger) error {

	_, err := svc.DeleteUser(&iam.DeleteUserInput{
		UserName: &user.Spec.Username,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateIamUser(svc *iam.IAM, user *govinddevv1alpha1.IamUser, l logr.Logger) error {
	_, err := svc.UpdateUser(&iam.UpdateUserInput{
		UserName:    &user.Status.Username,
		NewUserName: &user.Spec.Username,
	})
	if err != nil {
		return err
	}
	userUpdate, err := svc.GetUser(&iam.GetUserInput{UserName: &user.Spec.Username})
	user.Status.UserArn = *userUpdate.User.Arn
	return nil

}

func AwsIamSession(l logr.Logger) (*iam.IAM, error) {
	sess, err := session.NewSession()
	if err != nil {
		l.Error(err, "Failed to start aws session")
		return &iam.IAM{}, err
	}
	svc := iam.New(sess)
	return svc, nil
}

func (r *IamUserReconciler) checkDeletion(myfinalizer string,user *govinddevv1alpha1.IamUser,ctx context.Context,l logr.Logger) error{
	if user.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(user, myfinalizer) {
			controllerutil.AddFinalizer(user, myfinalizer)
			if err := r.Update(ctx, user); err != nil {
				return  err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(user, myfinalizer) {
			if err := r.DeleteUserReconcile(ctx,user,l); err != nil {
				return err
			}
			controllerutil.RemoveFinalizer(user, myfinalizer)
			if err := r.Update(ctx, user); err != nil {
				return err
			}
		}
	}
	return nil
}