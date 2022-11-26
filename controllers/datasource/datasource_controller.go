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

package datasource

import (
	"context"
	"database/sql"
	"fmt"
	datasourcev1alpha1 "github.com/yuanyp8/dbhero/apis/datasource/v1alpha1"
	"github.com/yuanyp8/dbhero/internal/mysql"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sync"
	"time"
)

// DataSourceReconciler reconciles a DataSource object
type DataSourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type DB struct {
	DBType    string  `json:"DBType,omitempty"`
	DBVersion string  `json:"DBVersion,omitempty"`
	DB        *sql.DB `json:"DB,omitempty"`
}

var (
	DbMap = make(map[string]*DB, 10)
	mux   = sync.Mutex{}
)

func GetDB(dbType, dbVersion string) *sql.DB {
	for _, db := range DbMap {
		if db.DBVersion == dbVersion && db.DBType == dbType {
			return db.DB
		}
	}
	return nil
}

//+kubebuilder:rbac:groups=datasource.dbhero.io,resources=datasources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=datasource.dbhero.io,resources=datasources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=datasource.dbhero.io,resources=datasources/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DataSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	flag := true
	l := log.FromContext(ctx).WithName("reconciling datasource")

	datasourceInstance := &datasourcev1alpha1.DataSource{}

	// 1. get datasource instance
	err := r.Get(ctx, req.NamespacedName, datasourceInstance)
	if errors.IsNotFound(err) {
		l.Error(err, "datasource instance not found", "name", req.Name)
		return ctrl.Result{}, err
	}
	if err != nil {
		l.Error(err, "unable to get datasource instance", "namespace", req.NamespacedName, "name", req.Name)
		return ctrl.Result{}, err
	}

	// 2. set annotations
	datasourceInstance.SetAnnotation()

	// 3. get datasource type
	dbType, err := datasourceInstance.GetType()
	if err != nil {
		l.Error(err, "datasource instance type if not allowed", "datasource type", dbType)
		return ctrl.Result{}, err
	}

	datasourceInstance.Status.Type = dbType

	// 4. set objectMeta labels
	datasourceInstance.SetLabels(map[string]string{"dbType": string(dbType), "database": req.Name})
	l.Info("set label succeed")

	// 5. set username and password
	if err := r.setUserNameAndPassword(ctx, datasourceInstance, req); err != nil {
		l.Error(err, "set username and password error")
		return ctrl.Result{}, err
	}
	l.Info("set username and password succeed")

	// 6. init connection
	db, err := r.getDB(datasourceInstance)
	if err != nil {
		flag = false
		l.Error(err, "get datasource backend connection pool error", "dbName", datasourceInstance.Name)
	}

	l.Info("get db succeed", "status", datasourceInstance.Status.IsConnected)

	// 7. ping MySQL
	if err = mysql.Ping(ctx, db); err != nil {
		l.Error(err, "ping mysql instance error", "dbName", datasourceInstance.Name)
	} else {
		l.Info("ping mysql database succeed")
		datasourceInstance.Status.LastPing = string(time.Now().String())
	}

	// dbhero current does not support any database-wide schema properties in Postgresql
	if datasourceInstance.Spec.Connection.Postgre != nil {
		l.Info("ignoring postgres database schema reconcile request")
		return ctrl.Result{}, fmt.Errorf("ignoring postgres database schema reconcile request")
	}

	if flag {
		l.Info("modify everything ok, continue update")
	}

	l.Info("debug", "username", datasourceInstance.Status.Auth.Username.Value, "password", datasourceInstance.Status.Auth.Password.Value, "type", datasourceInstance.Status.Type)

	/*	if err := r.Update(ctx, datasourceInstance); err != nil {
		l.Error(err, "modify the Datasource status failed")
		return ctrl.Result{Requeue: true}, err
	}*/

	// 6. set status.Version
	datasourceInstance.Status.Version = datasourceInstance.Spec.Connection.MySQL.Version

	// 7. update
	if err := r.Status().Update(ctx, datasourceInstance); err != nil {
		l.Error(err, "modify the  status failed")
		return ctrl.Result{Requeue: true}, err
	}

	l.Info("modify the Datasource status successful")

	return ctrl.Result{}, nil
}

//+kubebuilder:rbac:groups:core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *DataSourceReconciler) setUserNameAndPassword(ctx context.Context, ins *datasourcev1alpha1.DataSource, req ctrl.Request) error {
	l := log.FromContext(ctx).WithName("setUserNameAndPassword")

	// determine whether the Datasource.Spec.Connection.Auth is empty
	switch ins.Status.Type {
	case datasourcev1alpha1.MYSQL:
		// set username
		username := ins.GetMysqlAuth().Username
		secretUsername := &v1.Secret{}

		if username.IsEmpty() {
			//ins.Status.Auth.Username.Value = utils.UsernameGenerator(8)
			l.Error(nil, "username must be specified")
			return fmt.Errorf("username must be specified")
		}

		if username.Value == "" {
			// get secret
			if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: username.ValueFrom.SecretKerRef.Name}, secretUsername); err != nil {
				l.Error(err, "get secret from kubernetes failed", "secretName", username.ValueFrom.SecretKerRef.Name)
				return fmt.Errorf("get username from secret error: %s", err)
			}

			ins.Status.Auth.Username.Value = string(secretUsername.Data[username.ValueFrom.SecretKerRef.Key])

		} else {
			ins.Status.Auth.Username.Value = username.Value
		}

		password := ins.GetMysqlAuth().Password
		secretPassword := &v1.Secret{}

		if password.IsEmpty() {
			// ins.Status.Auth.Password.Value = utils.MustPasswordGenerator()
			l.Error(nil, "password must be specified")
			return fmt.Errorf("password must be specified")
		} else {
			if password.Value == "" {
				// get secret
				if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: password.ValueFrom.SecretKerRef.Name}, secretPassword); err != nil {
					l.Error(err, "get secret from kubernetes failed, will generate random password", "secretName", password.ValueFrom.SecretKerRef.Name)
					return fmt.Errorf("get password from secret error: %s", err)
				} else {
					ins.Status.Auth.Password.Value = string(secretPassword.Data[password.ValueFrom.SecretKerRef.Key])
				}
			} else {
				ins.Status.Auth.Password.Value = password.Value
			}
		}
		return nil
	default:
		l.Error(fmt.Errorf("the datasource was not implement"), "")
		return fmt.Errorf("the datasource was not implement")
	}
}

func (r *DataSourceReconciler) getDB(ins *datasourcev1alpha1.DataSource) (*sql.DB, error) {
	if _, exist := DbMap[ins.Name]; exist {
		// 有了也需要重新刷新下
		if err := mysql.Ping(context.TODO(), DbMap[ins.Name].DB); err == nil {
			ins.Status.IsConnected = true
			return DbMap[ins.Name].DB, nil
		}
	}
	mux.Lock()
	defer mux.Unlock()

	db, err := mysql.InitPool(ins)
	if err != nil {
		ins.Status.IsConnected = false
		return nil, err
	}
	ins.Status.IsConnected = true

	DbMap[ins.Name] = &DB{
		DBType:    string(ins.Status.Type),
		DBVersion: ins.Status.Version,
		DB:        db,
	}

	return db, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datasourcev1alpha1.DataSource{}).
		Complete(r)
}
