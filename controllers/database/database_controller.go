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
	"fmt"
	"github.com/yuanyp8/dbhero/controllers/datasource"
	"github.com/yuanyp8/dbhero/internal/random"
	"k8s.io/apimachinery/pkg/api/errors"

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

	// 2. get *sql.DB
	dbVersion, dbType := databaseInstance.Spec.DBVersion, databaseInstance.Spec.DBVersion

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

	// 5. determine database's datasource type
	switch dbType {

	case "postgresql":
		l.Error(fmt.Errorf("this db type is not implement"), "type", databaseInstance.Spec.DBType)
		return ctrl.Result{}, fmt.Errorf("this db type is not implement")

	case "mysql":

	default:
		l.Error(nil, "not implement database type", "type", dbType)
		return ctrl.Result{}, fmt.Errorf("not implement database type")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.Database{}).
		Complete(r)
}

//+kubebuilder:rbac:groups:core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *DatabaseReconciler) getUsername(ctx context.Context, ins *databasev1alpha1.Database) (string, error) {
	l := log.FromContext(ctx).WithName("getUsername")

	// determine if the auth is empty
	if ins.Spec.Auth.Username.IsEmpty() {
		// generate username

		for {
			username := random.UsernameGenerator(8)
			// if username exist in the datasource

		}
	}
}
