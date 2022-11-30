/*
Copyright 2022 yuanyp8.

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

package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/yuanyp8/dbhero/controllers/datasource"
	"github.com/yuanyp8/dbhero/internal/mysql"
	"github.com/yuanyp8/dbhero/internal/random"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasev1alpha1 "github.com/yuanyp8/dbhero/apis/database/v1alpha1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.dbhero.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.dbhero.io,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.dbhero.io,resources=databases/finalizers,verbs=update
//+kubebuilder:rbac:groups=datasource.dbhero.io,resources=datasources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=datasource.dbhero.io,resources=datasources/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithName("reconciling database")

	databaseInstance := &databasev1alpha1.Database{}

	// 1. get database instance
	err := r.Get(ctx, req.NamespacedName, databaseInstance)
	if errors.IsNotFound(err) {
		l.Error(err, "database instance not found", "name", req.Name)
		return ctrl.Result{}, err
	}
	if err != nil {
		l.Error(err, "unable to get database instance", "namespace", req.NamespacedName, "name", req.Name)
		return ctrl.Result{}, err
	}

	// 2. determine the dbType and dbVersion to detect the *sql.DB

	dbVersion, dbType := databaseInstance.Spec.DBVersion, databaseInstance.Spec.DBVersion

	switch dbType {

	case "postgresql":
		l.Error(fmt.Errorf("this db type is not implement"), "type", databaseInstance.Spec.DBType)
		return ctrl.Result{}, fmt.Errorf("this db type is not implement")

	case "mysql":

	default:
		l.Error(nil, "not implement database type", "type", dbType)
		return ctrl.Result{}, fmt.Errorf("not implement database type")
	}

	db := datasource.GetDB(dbType, dbVersion)
	if db == nil {
		l.Error(fmt.Errorf("there is no vaild datasouce"), "", "DbVersion", dbVersion, "DbType", dbType)
		return ctrl.Result{}, fmt.Errorf("there is no vaild datasouce")
	}

	// 3. set annotations
	databaseInstance.SetAnnotation()

	// 4. add labels
	databaseInstance.SetLabels(map[string]string{"dbType": string(databaseInstance.Spec.DBType), "database": req.Name})
	l.Info("set label succeed")

	// 5. determine username and password
	username, err := r.getUsername(ctx, databaseInstance, db)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get username error")
	}

	password, err := r.getPassword(ctx, databaseInstance, db)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get password error")
	}

	// 6. determine database schema name if already exists.
	if err := mysql.DetectDatabase(ctx, db, databaseInstance.Spec.DBName); err != nil {
		l.Error(err, "database schema is already in the datasource", "dbName", databaseInstance.Spec.DBName)
		return ctrl.Result{}, err
	}

	// 7. create user create database schema and grant privileges.
	err = mysql.CreateUserDatabaseSTMT(ctx, db, username, password, databaseInstance.Spec.DBName, string(databaseInstance.Spec.DefaultCharacterSet), string(databaseInstance.Spec.DefaultCollation))
	if err != nil {
		return ctrl.Result{}, err
	}

	// 8. update status
	if err := r.Status().Update(ctx, databaseInstance); err != nil {
		l.Error(err, "modify the  status failed")
		return ctrl.Result{Requeue: true}, err
	}

	l.Info("modify the Datasource status successful")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.Database{}).
		Complete(r)
}

//+kubebuilder:rbac:groups:core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *DatabaseReconciler) getUsername(ctx context.Context, ins *databasev1alpha1.Database, db *sql.DB) (string, error) {
	l := log.FromContext(ctx).WithName("getUsername")

	// determine if the auth is empty
	if ins.Spec.Auth.Username.IsEmpty() {
		// generate username

		for {
			username := random.UsernameGenerator(8)
			// if username exist in the datasource
			value, err := mysql.DetectUser(ctx, db, username)
			if value == 0 {
				return username, nil
			} else if value == -1 {
				return "", err
			} else {
				continue
			}
		}
	}
	// get user from secret
	if ins.Spec.Auth.Username.ValueFrom != nil {

		secretUsername := &v1.Secret{}

		if err := r.Get(ctx, types.NamespacedName{Namespace: ins.Namespace, Name: ins.Spec.Auth.Username.ValueFrom.SecretKerRef.Name}, secretUsername); err != nil {
			l.Error(err, "get secret from kubernetes failed", "secretName", ins.Spec.Auth.Username.ValueFrom.SecretKerRef.Name)
			return "", fmt.Errorf("get username from secret error: %s", err)
		}

		return string(secretUsername.Data[ins.Spec.Auth.Username.ValueFrom.SecretKerRef.Key]), nil
	}
	return ins.Spec.Auth.Username.Value, nil
}

//+kubebuilder:rbac:groups:core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *DatabaseReconciler) getPassword(ctx context.Context, ins *databasev1alpha1.Database, db *sql.DB) (string, error) {
	l := log.FromContext(ctx).WithName("getPassword")

	// determine if the auth is empty
	if ins.Spec.Auth.Password.IsEmpty() {
		// generate password
		passwd, err := random.PasswordGenerator()
		if err != nil {
			return "", fmt.Errorf("generate password error")
		}
		return passwd, nil
	}

	// get password from secret
	if ins.Spec.Auth.Password.ValueFrom != nil {

		secretPassword := &v1.Secret{}

		if err := r.Get(ctx, types.NamespacedName{Namespace: ins.Namespace, Name: ins.Spec.Auth.Password.ValueFrom.SecretKerRef.Name}, secretPassword); err != nil {
			l.Error(err, "get secret from kubernetes failed", "secretName", ins.Spec.Auth.Password.ValueFrom.SecretKerRef.Name)
			return "", fmt.Errorf("get password from secret error: %s", err)
		}

		return string(secretPassword.Data[ins.Spec.Auth.Password.ValueFrom.SecretKerRef.Key]), nil
	}
	return ins.Spec.Auth.Password.Value, nil
}
